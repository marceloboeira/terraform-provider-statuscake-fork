package statuscake

import (
	"fmt"

	"context"

	scapi "github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceStatusCakeContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: CreateContactGroup,
		Update: UpdateContactGroup,
		Delete: DeleteContactGroup,
		Read:   ReadContactGroup,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ping_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile_numbers": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"email_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"integration_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func CreateContactGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scapi.APIClient)

	res, err := client.CreateContactGroup(context.Background()).
		Name(d.Get("name").(string)).
		PingURL(d.Get("ping_url").(string)).
		MobileNumbers(castSetToSliceStrings(d.Get("mobile_numbers").(*schema.Set).List())).
		EmailAddresses(castSetToSliceStrings(d.Get("email_addresses").(*schema.Set).List())).
		Integrations(castSetToSliceStrings(d.Get("integration_ids").(*schema.Set).List())).
		Execute()

	if err != nil {
		return fmt.Errorf("Error creating StatusCake ContactGroup: %s", err.Error())
	}

	d.SetId(res.Data.NewID)

	return ReadContactGroup(d, meta)
}

func UpdateContactGroup(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func DeleteContactGroup(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*scapi.Client)
	// id, _ := strconv.Atoi(d.Id())
	// log.Printf("[DEBUG] Deleting StatusCake ContactGroup: %s", d.Id())
	// err := scapi.NewContactGroups(client).Delete(id)

	// return err
	return nil
}

func ReadContactGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scapi.APIClient)

	res, err := client.GetContactGroup(context.Background(), d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("Error Getting StatusCake ContactGroup Details for %s: Error: %s", d.Id(), err)
	}

	d.Set("name", res.Data.Name)
	d.Set("mobile_numbers", orEmptySlice(res.Data.MobileNumbers))
	d.Set("email_addresses", orEmptySlice(res.Data.EmailAddresses))
	d.Set("integration_ids", orEmptySlice(res.Data.Integrations))
	d.Set("ping_url", res.Data.PingURL)

	return nil
}

func orEmptySlice(a []string) []string {
	if a == nil || len(a) == 0 {
		return []string{}
	}

	return a
}
