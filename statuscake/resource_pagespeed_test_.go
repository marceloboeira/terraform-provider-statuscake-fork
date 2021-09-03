package statuscake

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	statuscake "github.com/StatusCakeDev/statuscake-go"
)

func resourceStatusCakePagespeedTest() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatusCakePagespeedTestCreate,
		Read:   resourceStatusCakePagespeedTestRead,
		Update: resourceStatusCakePagespeedTestUpdate,
		Delete: resourceStatusCakePagespeedTestDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alert_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Alert confiuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_bigger": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							Description:  "An alert will be sent if the size of the page is larger than this value (kb). A value of 0 prevents alerts being sent",
							ValidateFunc: validation.IntAtLeast(0),
						},
						"alert_slower": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							Description:  "An alert will be sent if the load time of the page exceeds this value (ms). A value of 0 prevents alerts being sent",
							ValidateFunc: validation.IntAtLeast(0),
						},
						"alert_smaller": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							Description:  "An alert will be sent if the size of the page is smaller than this value (kb). A value of 0 prevents alerts being sent",
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},
			"check_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of minutes between checks",
				ValidateFunc: Int32InSlice(statuscake.PagespeedTestCheckRateValues()),
			},
			"contact_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of contact group IDs",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringIsNumerical,
				},
			},
			"location": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Assigned monitoring location on which checks will be run",
				Elem: &schema.Resource{
					Schema: locationSchema(),
				},
			},
			"location_iso": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Testing location",
				ValidateFunc: validation.StringInSlice(statuscake.PagespeedTestLocationISOValues(), false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the check",
			},
			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the check should be run",
			},
			"website_url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "URL or IP address of the website under test",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
	}
}

func resourceStatusCakePagespeedTestCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)

	config := expandPagespeedCheckAlertConfig(g.Get("alert_config"))
	req := client.CreatePagespeedTest(context.Background()).
		AlertBigger(int32(config["alert_bigger"].(int))).
		AlertSlower(int64(config["alert_slower"].(int))).
		AlertSmaller(int32(config["alert_smaller"].(int))).
		CheckRate(statuscake.PagespeedTestCheckRate(int32(d.Get("check_interval").(int)))).
		LocationISO(statuscake.PagespeedTestLocationISO(d.Get("location_iso").(string))).
		Name(d.Get("name").(string)).
		Paused(d.Get("paused").(bool)).
		WebsiteURL(d.Get("website_url").(string))

	if contactGroups, ok := d.GetOk("contact_groups"); ok {
		req = req.ContactGroups(toStringSlice(contactGroups.([]interface{})))
	}

	log.Print("[DEBUG] Creating StatusCake pagespeed check")

	res, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to create pagespeed check: %w", err)
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakePagespeedTestRead(d, meta)
}

func resourceStatusCakePagespeedTestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	res, err := client.GetPagespeedTest(context.Background(), id).Execute()
	if err != nil {
		return fmt.Errorf("failed to get pagespeed check with ID: %w", err)
	}

	d.Set("check_interval", int(res.Data.CheckRate))
	d.Set("contact_groups", res.Data.ContactGroups)
	d.Set("location", flattenPagespeedCheckLocation(res.Data.Location))
	d.Set("location_iso", string(res.Data.LocationISO))
	d.Set("name", res.Data.Name)
	d.Set("paused", res.Data.Paused)
	d.Set("website_url", res.Data.WebsiteURL)

	// Alert config is an artificial block. The actual data is flattened in the
	// response.
	d.Set("alert_config", flattenPagespeedCheckAlertConfig(res.Data))

	return nil
}

func resourceStatusCakePagespeedTestUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	req := client.UpdatePagespeedTest(context.Background(), id)

	if alertConfig, ok := d.GetOk("alert_config"); ok {
		config := expandPagespeedCheckAlertConfig(alertConfig)
		req = req.AlertBigger(int32(alertBigger.(int))).
			AlertSlower(int64(alertSlower.(int))).
			AlertSmaller(int32(alertSmaller.(int)))
	}

	if checkRate, ok := d.GetOk("check_interval"); ok {
		req = req.CheckRate(statuscake.PagespeedTestCheckRate(int32(checkRate.(int))))
	}

	if contactGroups, ok := d.GetOk("contact_groups"); ok {
		req = req.ContactGroups(toStringSlice(contactGroups.([]interface{})))
	}

	if location, ok := d.GetOk("location_iso"); ok {
		req = req.LocationISO(statuscake.PagespeedTestLocationISO(location.(string)))
	}

	if name, ok := d.GetOk("name"); ok {
		req = req.Name(name.(string))
	}

	if paused, ok := d.GetOk("paused"); ok {
		req = req.Paused(paused.(bool))
	}

	log.Printf("[DEBUG] Updating StatusCake pagespeed check with ID: %s", id)

	if err := req.Execute(); err != nil {
		return fmt.Errorf("failed to update pagespeed check: %w", err)
	}

	return resourceStatusCakePagespeedTestRead(d, meta)
}

func resourceStatusCakePagespeedTestDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake pagespeed check with ID: %s", id)

	if err := client.DeletePagespeedTest(context.Background(), id).Execute(); err != nil {
		return fmt.Errorf("failed to delete pagespeed check with id %s: %w", id, err)
	}

	return resourceStatusCakePagespeedTestRead(d, meta)
}

func expandPagespeedCheckAlertConfig(v []interface{}) map[string]interface{} {
	if v == nil || len(v) == 0 {
		return map[string]interface{}{}
	}

	return v[0].(map[string]interface{})
}

func flattenPagespeedCheckAlertConfig(data *statuscake.PagespeedTest) []interface{} {
	if data == nil {
		return []interface{}{}
	}

	return []map[string]interface{}{
		map[string]interface{}{
			"alert_bigger":  data.AlertBigger,
			"alert_slower":  data.AlertSlower,
			"alert_smaller": data.AlertSmaller,
		},
	}
}

func flattenPagespeedCheckLocation(data *statuscake.MonitoringLocation) interface{} {
	return map[string]interface{}{
		"description": "",
	}
}
