package statuscake

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	statuscake "github.com/StatusCakeDev/statuscake-go"
)

func resourceStatusCakeMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatusCakeMaintenanceWindowCreate,
		Read:   resourceStatusCakeMaintenanceWindowRead,
		Update: resourceStatusCakeMaintenanceWindowUpdate,
		Delete: resourceStatusCakeMaintenanceWindowDelete,

		// Used by `terraform import`.
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceStatusCakeMaintenanceWindowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)

	req := client.CreateMaintenanceWindow(context.Background()).
		Name(d.Get("name").(string))

	log.Print("[DEBUG] Creating StatusCake maintenance window test")

	res, err := req.Execute()
	if err != nil {
		return fmt.Errorf("failed to create maintenance window test: %w", err)
	}

	d.SetId(res.Data.NewID)
	return resourceStatusCakeMaintenanceWindowRead(d, meta)
}

func resourceStatusCakeMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	res, err := client.GetMaintenanceWindow(context.Background(), id).Execute()
	if err != nil {
		return fmt.Errorf("failed to get maintenance window test with ID: %w", err)
	}

	d.Set("name", res.Data.Name)

	return nil
}

func resourceStatusCakeMaintenanceWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	req := client.UpdateMaintenanceWindow(context.Background(), id)

	if name, ok := d.GetOk("name"); ok {
		req = req.Name(name.(string))
	}

	log.Printf("[DEBUG] Updating StatusCake maintenance window test with ID: %s", id)

	if err := req.Execute(); err != nil {
		return fmt.Errorf("failed to update maintenance window test: %w", err)
	}

	return resourceStatusCakeMaintenanceWindowRead(d, meta)
}

func resourceStatusCakeMaintenanceWindowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.APIClient)
	id := d.Id()

	log.Printf("[DEBUG] Deleting StatusCake maintenance window test with ID: %s", id)

	if err := client.DeleteMaintenanceWindow(context.Background(), id).Execute(); err != nil {
		return fmt.Errorf("failed to delete maintenance window test with id %s: %w", id, err)
	}

	return resourceStatusCakeMaintenanceWindowRead(d, meta)
}
