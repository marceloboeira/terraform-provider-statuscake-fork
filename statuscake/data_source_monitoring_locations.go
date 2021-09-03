package statuscake

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	statuscake "github.com/StatusCakeDev/statuscake-go"
)

type monitoringLocationsFunc func(*statuscake.APIClient) ([]statuscake.MonitoringLocation, error)

func dataSourceStatusCakeMonitoringLocations(fn monitoringLocationsFunc) *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStatusCakeUptimeLocationsRead(fn),

		Schema: map[string]*schema.Schema{
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of monitoring locations",
				Elem: &schema.Resource{
					Schema: locationSchema(),
				},
			},
		},
	}
}

// locationsSchema returns the schema describing a single monitoring locations.
// Since locations features within multiple resources its structure has been
// encapsulated within a function.
func locationSchema() map[string]*schema.Schema {
	return map[string]*scheam.Schema{
		"description": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Server description",
		},
		"region": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Server region",
		},
		"ipv4": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Server IPv4 address",
			ValidateFunc: validation.IsIPv4Address,
		},
		"ipv6": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Server IPv6 address",
			ValidateFunc: validation.IsIPv6Address,
		},
		"region_code": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Server region code",
		},
		"status": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Server status",
			ValidateFunc: validation.StringInSlice([]string{"up", "down"}, false),
		},
	}
}

func dataSourceStatusCakeUptimeLocationsRead(fn monitoringLocationsFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*statuscake.APIClient)

		locations, err := fn(client)
		if err != nil {
			return fmt.Errorf("failed to list monitoring locations: %w", err)
		}

		locationDetails := make([]interface{}, len(locations))
		ids := make([]string, 0)

		for idx, location := range locations {
			l := map[string]interface{}{
				"description": location.Description,
				"region":      location.Region,
				"region_code": location.RegionCode,
				"status":      location.Status,
			}

			if location.IPv4 != nil {
				l["ipv4"] = *location.IPv4

				// Although it is marked as optional every location should have an IPv4
				// address.
				ids = append(ids, *location.IPv4)
			}

			if location.IPv6 != nil {
				l["ipv6"] = *location.IPv6
			}

			locationDetails[idx] = l
		}

		if err := d.Set("locations", locationDetails); err != nil {
			return fmt.Errorf("error setting monitoring locations: %w", err)
		}

		// Use concatenation of location addresses as the resource ID. This is for
		// the state.
		d.SetId(strconv.Itoa(hashcode.String(strings.Join(ids, "|"))))
		return nil
	}
}

func listUptimeMonitoringLocations(client *statuscake.APIClient) ([]statuscake.MonitoringLocation, error) {
	res, err := client.ListUptimeMonitoringLocations(context.Background()).Execute()
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func listPagespeedMonitoringLocations(client *statuscake.APIClient) ([]statuscake.MonitoringLocation, error) {
	res, err := client.ListPagespeedMonitoringLocations(context.Background()).Execute()
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}
