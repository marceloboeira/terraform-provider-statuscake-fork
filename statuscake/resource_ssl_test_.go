package statuscake

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	statuscake "github.com/StatusCakeDev/statuscake-go"
)

func resourceStatusCakeSSLCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatusCakeSSLCheckCreate,
		Read:   resourceStatusCakeSSLCheckRead,
		Update: resourceStatusCakeSSLCheckUpdate,
		Delete: resourceStatusCakeSSLCheckDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alert_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Alert configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_at": {
							Type:        schema.TypeSet,
							Optional:    true,
							MinItems:    1,
							MaxItems:    3,
							Description: "List representing when alerts should be sent (days). Must be exactly 3 numerical values",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"on_reminder": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alert reminders",
						},
						"on_expiry": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alerts when the SSL certificate is to expire",
						},
						"on_broken": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alerts when SSL certificate issues are found",
						},
						"on_mixed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable alerts when mixed content is found",
						},
					},
				},
			},
			"check_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of seconds between checks",
				ValidateFunc: Int32InSlice(statuscake.SSLTestCheckRateValues()),
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
			"follow_redirects": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to follow redirects when testing",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname of the server under test",
			},
			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the check should be run",
			},
			"user_agent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom user agent string set when testing",
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "URL of the endpoint under test",
				ValidateFunc: validation.IsURLWithHTTPS,
			},
		},
	}
}

func resourceStatusCakeSSLCheckCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)

	config := expandSSLCheckAlertConfig(d.Get("alert_config"))
	req := client.CreateSslTest(context.Background()).
		AlertAt(toStringSlice(config["alert_at"].([]interface{}))).
		AlertBroken(config["on_broken"].(bool)).
		AlertExpiry(config["on_expiry"].(bool)).
		AlertMixed(config["on_mixed"].(bool)).
		AlertReminder(config["on_reminder"].(bool)).
		CheckRate(statuscake.SSLTestCheckRate(int32(d.Get("check_interval").(int)))).
		FollowRedirects(d.Get("follow_redirects").(bool)).
		Paused(d.Get("paused").(bool)).
		WebsiteURL(d.Get("url").(string))

	if contactGroups, ok := d.GetOk("contact_groups"); ok {
		req = req.ContactGroups(toStringSlice(contactGroups.([]interface{})))
	}

	if hostname, ok := d.GetOk("hostname"); ok {
		req = req.Hostname(hostname.(string))
	}

	if userAgent, ok := d.GetOk("user_agent"); ok {
		req = req.UserAgent(userAgent.(string))
	}

	log.Print("[DEBUG] Creating StatusCake SSL check")

	res, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to create SSL check: %w", err)
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeSSLCheckRead(d, meta)
}

func resourceStatusCakeSSLCheckRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	res, err := client.GetSslTest(context.Background(), id).Execute()
	if err != nil {
		return fmt.Errorf("failed to get SSL check with ID: %w", err)
	}

	d.Set("check_interval", int(res.Data.CheckRate))
	d.Set("contact_groups", res.Data.ContactGroups)
	d.Set("follow_redirects", res.Data.FollowRedirects)
	d.Set("paused", res.Data.Paused)
	d.Set("url", res.Data.WebsiteURL)

	// Alert config is an artificial block. The actual data is flattened in the
	// response.
	d.Set("alert_config", flattenSSLCheckAlertConfig(res.Data))

	if res.Data.Hostname != nil {
		d.Set("hostname", *res.Data.Hostname)
	}

	if res.Data.UserAgent != nil {
		d.Set("user_agent", *res.Data.UserAgent)
	}

	return nil
}

func resourceStatusCakeSSLCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	req := client.UpdateSslTest(context.Background(), id)

	if alertConfig, ok := d.GetOk("alert_config"); ok {
		config := expandSSLCheckAlertConfig(alertConfig)
		req = req.AlertAt(toStringSlice(config["alert_at"].([]interface{}))).
			AlertReminder(config["on_reminder"].(bool)).
			AlertExpiry(config["on_expiry"].(bool)).
			AlertBroken(config["on_broken"].(bool)).
			AlertMixed(config["on_mixed"].(bool))
	}

	if checkRate, ok := d.GetOk("check_interval"); ok {
		req = req.CheckRate(statuscake.SSLTestCheckRate(int32(checkRate.(int))))
	}

	if contactGroups, ok := d.GetOk("contact_groups"); ok {
		req = req.ContactGroups(toStringSlice(contactGroups.([]interface{})))
	}

	if followRedirects, ok := d.GetOk("followRedirects"); ok {
		req = req.FollowRedirects(followRedirects.(bool))
	}

	if paused, ok := d.GetOk("paused"); ok {
		req = req.Paused(paused.(bool))
	}

	if hostname, ok := d.GetOk("hostname"); ok {
		req = req.Hostname(hostname.(string))
	}

	if userAgent, ok := d.GetOk("user_agent"); ok {
		req = req.UserAgent(userAgent.(string))
	}

	log.Printf("[DEBUG] Updating StatusCake SSL check with ID: %s", id)

	if err := req.Execute(); err != nil {
		return fmt.Errorf("failed to update SSL check: %w", err)
	}

	return resourceStatusCakeSSLCheckRead(d, meta)
}

func resourceStatusCakeSSLCheckDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake SSL check with ID: %s", id)

	if err := client.DeleteSslTest(context.Background(), id).Execute(); err != nil {
		return fmt.Errorf("failed to delete SSL check with id %s: %w", id, err)
	}

	return resourceStatusCakeSSLCheckRead(d, meta)
}

func expandSSLCheckAlertConfig(v []interface{}) map[string]interface{} {
	if v == nil || len(v) == 0 {
		return map[string]interface{}{}
	}

	return v[0].(map[string]interface{})
}

func flattenSSLCheckAlertConfig(data *statuscake.SSLTest) []interface{} {
	if data == nil {
		return []interface{}{}
	}

	return []map[string]interface{}{
		map[string]interface{}{
			"alert_at":    data.AlertAt,
			"on_reminder": data.AlertReminder,
			"on_expiry":   data.AlertExpiry,
			"on_broken":   data.AlertBroken,
			"on_mixed":    data.AlertMixed,
		},
	}
}
