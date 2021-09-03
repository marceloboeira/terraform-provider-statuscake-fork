package statuscake

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	statuscake "github.com/StatusCakeDev/statuscake-go"
)

func resourceStatusCakeUptimeHTTPTest() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatusCakeUptimeHTTPTestCreate,
		Read:   resourceStatusCakeUptimeHTTPTestRead,
		Update: resourceStatusCakeUptimeHTTPTestUpdate,
		Delete: resourceStatusCakeUptimeHTTPTestDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"check_rate": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: Int32InSlice(statuscake.UptimeTestCheckRateValues()),
			},
			"confirmation": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(0, 3),
			},
			"contact_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringIsNumerical,
				},
			},
			"content_matchers": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"include_header": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"matcher": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "CONTAINS_STRING",
							ValidateFunc: validation.StringInSlice([]string{"CONTAINS_STRING", "NOT_CONTAINS_STRING"}, false),
						},
					},
				},
			},
			"dns_check": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_ips": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema:       schema.TypeString,
								ValidateFunc: validation.IsIPAddress,
							},
						},
						"dns_server": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
					},
				},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"http_check": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: authenticationSchema(),
							},
						},
						"body": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"headers": {
							Type:     schema.TypeMap,
							Computed: true,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"final_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"follow_redirects": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/",
						},
						"request_method": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "HTTP",
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HEAD"}, false),
						},
						"status_codes": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"use_jar": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"user_agent": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"validate_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"icmp_check": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// There are no special fields for an ICMP check. All that is required
				// is the address which is supplied in the `monitoried_resource` block.
				Elem:         &schema.Resource{},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: locationSchema(),
				},
			},
			"monitored_resource": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         scheam.TypeString,
							Required:     true,
							ValidateFunc: validation.Any(validation.IsURLWithHTTPorHTTPS, validation.IsIPAddress),
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tcp_check": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: authenticationSchema(),
							},
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				ExactlyOneOf: []string{"dns_check", "http_check", "icmp_check", "tcp_check"},
			},
			"trigger_rate": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
}

// authenticationSchema returns the schema describing the necessary fields for
// authentication an uptime request. Since authentication features within
// multiple check types its structure has been encapsulated within a function.
func authenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:     schema.TypeString,
			Required: true,
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
	}
}

func resourceStatusCakeUptimeHTTPTestCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)

	req := client.CreateUptimeTest(context.Background()).
		WebsiteURL(d.Get("website_url").(string)).
		CheckRate(statuscake.UptimeTestCheckRate(int32(d.Get("check_rate").(int)))).
		FollowRedirects(d.Get("follow_redirects").(bool)).
		Paused(d.Get("paused").(bool))

	if contactGroups, ok := d.GetOk("contact_groups"); ok {
		req = req.ContactGroups(toStringSlice(contactGroups.([]interface{})))
	}

	if host, ok := d.GetOk("host"); ok {
		req = req.Host(host.(string))
	}

	if userAgent, ok := d.GetOk("user_agent"); ok {
		req = req.UserAgent(userAgent.(string))
	}

	log.Print("[DEBUG] Creating StatusCake uptime test")

	res, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to create uptime test: %w", err)
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeUptimeHTTPTestRead(d, meta)
}

func resourceStatusCakeUptimeHTTPTestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	res, err := client.GetUptimeTest(context.Background(), id).Execute()
	if err != nil {
		return fmt.Errorf("failed to get uptime test with ID: %w", err)
	}

	d.Set("website_url", res.Data.WebsiteURL)
	d.Set("check_rate", int(res.Data.CheckRate))
	d.Set("contact_groups", res.Data.ContactGroups)
	d.Set("follow_redirects", res.Data.FollowRedirects)
	d.Set("paused", res.Data.Paused)

	if res.Data.Host != nil {
		d.Set("host", *res.Data.Host)
	}

	if res.Data.UserAgent != nil {
		d.Set("user_agent", *res.Data.UserAgent)
	}

	return nil
}

func resourceStatusCakeUptimeHTTPTestUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	req := client.UpdateUptimeTest(context.Background(), id)

	if checkRate, ok := d.GetOk("check_rate"); ok {
		req = req.CheckRate(statuscake.UptimeTestCheckRate(int32(checkRate.(int))))
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

	if host, ok := d.GetOk("host"); ok {
		req = req.Host(host.(string))
	}

	if userAgent, ok := d.GetOk("user_agent"); ok {
		req = req.UserAgent(userAgent.(string))
	}

	log.Printf("[DEBUG] Updating StatusCake uptime test with ID: %s", id)

	if err := req.Execute(); err != nil {
		return fmt.Errorf("failed to update uptime test: %w", err)
	}

	return resourceStatusCakeUptimeHTTPTestRead(d, meta)
}

func resourceStatusCakeUptimeHTTPTestDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake uptime test with ID: %s", id)

	if err := client.DeleteUptimeTest(context.Background(), id).Execute(); err != nil {
		return fmt.Errorf("failed to delete uptime test with id %s: %w", id, err)
	}

	return resourceStatusCakeUptimeHTTPTestRead(d, meta)
}
