package statuscake

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	statuscake "github.com/StatusCakeDev/statuscake-go"
)

func resourceStatusCakeContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatusCakeContactGroupCreate,
		Read:   resourceStatusCakeContactGroupRead,
		Update: resourceStatusCakeContactGroupUpdate,
		Delete: resourceStatusCakeContactGroupDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email_addressess": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of email addresses",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: IsEmailAddress,
				},
			},
			"integration": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of integration IDs",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringIsNumerical,
				},
			},
			"mobile_numbers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of international format mobile phone numbers",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the contact group",
			},
			"ping_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "URL or IP address of an endpoint to push uptime events. Currently this only supports HTTP GET endpoints",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
	}
}

func resourceStatusCakeContactGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)

	req := client.CreateContactGroup(context.Background()).
		Name(d.Get("name").(string))

	if emailAddresses, ok := d.GetOk("email_addressess"); ok {
		req = req.EmailAddresses(toStringSlice(emailAddresses.(*schema.Set).List()))
	}

	if integrations, ok := d.GetOk("integrations"); ok {
		req = req.Integrations(toStringSlice(integrations.(*schema.Set).List()))
	}

	if mobileNumbers, ok := d.GetOk("mobile_numbers"); ok {
		req = req.MobileNumbers(toStringSlice(mobileNumbers.(*schema.Set).List()))
	}

	if url, ok := d.GetOk("ping_url"); ok {
		req = req.PingURL(url.(string))
	}

	log.Print("[DEBUG] Creating StatusCake contact group")

	res, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to create contact group: %w", err)
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeContactGroupRead(d, meta)
}

func resourceStatusCakeContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	res, err := client.GetContactGroup(context.Background(), id).Execute()
	if err != nil {
		return fmt.Errorf("failed to get contact group with ID: %w", err)
	}

	d.Set("email_addresses", res.Data.EmailAddresses)
	d.Set("integrations", res.Data.Integrations)
	d.Set("mobile_numbers", res.Data.MobileNumbers)
	d.Set("name", res.Data.Name)

	if res.Data.PingURL != nil {
		d.Set("ping_url", *res.Data.PingURL)
	}

	return nil
}

func resourceStatusCakeContactGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	req := client.UpdateContactGroup(context.Background(), id)

	if emailAddresses, ok := d.GetOk("email_addressess"); ok {
		req = req.EmailAddresses(toStringSlice(emailAddresses.(*schema.Set).List()))
	}

	if integrations, ok := d.GetOk("integrations"); ok {
		req = req.Integrations(toStringSlice(integrations.(*schema.Set).List()))
	}

	if mobileNumbers, ok := d.GetOk("mobile_numbers"); ok {
		req = req.MobileNumbers(toStringSlice(mobileNumbers.(*schema.Set).List()))
	}

	if name, ok := d.GetOk("name"); ok {
		req = req.Name(name.(string))
	}

	if url, ok := d.GetOk("ping_url"); ok {
		req = req.PingURL(url.(string))
	}

	log.Printf("[DEBUG] Updating StatusCake contact group with ID: %s", id)

	if err := req.Execute(); err != nil {
		return fmt.Errorf("failed to update contact group: %w", err)
	}

	return resourceStatusCakeContactGroupRead(d, meta)
}

func resourceStatusCakeContactGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake contact group with ID: %s", id)

	if err := client.DeleteContactGroup(context.Background(), id).Execute(); err != nil {
		return fmt.Errorf("failed to delete contact group with id %s: %w", id, err)
	}

	return resourceStatusCakeContactGroupRead(d, meta)
}
