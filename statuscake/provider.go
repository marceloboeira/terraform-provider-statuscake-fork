package statuscake

import (
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	statuscake "github.com/StatusCakeDev/statuscake-go"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("STATUSCAKE_API_TOKEN", nil),
				Description:  "The API token for operations",
				ValidateFunc: validation.StringMatch(regexp.MustCompile("[0-9a-zA-Z)]{20,30}"), "API token must only contain characters 0-9, a-zA-Z and underscores"),
			},
			"rps": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_RPS", 4),
				Description: "RPS limit to apply when making calls to the API",
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_RETRIES", 3),
				Description: "Maximum number of retries to perform when an API request fails",
			},
			"min_backoff": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_MIN_BACKOFF", 1),
				Description: "Minimum backoff period in seconds after failed API calls",
			},
			"max_backoff": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_MAX_BACKOFF", 30),
				Description: "Maximum backoff period in seconds after failed API calls",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"statuscake_contact_group":      resourceStatusCakeContactGroup(),
			"statuscake_maintenance_window": resourceStatusCakeMaintenanceWindow(),
			"statuscake_pagespeed_test":     resourceStatusCakePagespeedTest(),
			"statuscake_ssl_test":           resourceStatusCakeSSLTest(),
			"statuscake_uptime_test":        resourceStatusCakeUptimeTest(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"statuscake_pagespeed_monitoring_locations": dataSourceStatusCakeMonitoringLocations(listPagespeedMonitoringLocations),
			"statuscake_uptime_monitoring_locations":    dataSourceStatusCakeMonitoringLocations(listUptimeMonitoringLocations),
		},
		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure parses the config into the Terraform provider meta object
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// limitOpt := statuscake.WithRateLimit(float64(d.Get("rps").(int)))
	// retryOpt := statuscake.WithRetryPolicy(d.Get("retries").(int), d.Get("min_backoff").(int), d.Get("max_backoff").(int))

	var apiToken string
	if v, ok := d.GetOk("api_token"); ok {
		apiToken = v.(string)
	} else {
		return nil, errors.New("credentials are not set correctly")
	}

	client := statuscake.NewAPIClient(apiToken)
	// config := client.GetConfig()

	// TODO: configure the HTTP client to support expoenential backoff. Try using
	// the HTTP roundtripper interface
	// This will likely be added to the Go SDK directly so it can be used outside
	// of Terraform.

	return client, nil
}
