package statuscake

import (
	scapi "github.com/StatusCakeDev/statuscake-go"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_APIKEY", nil),
				Description: "API Key for StatusCake",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"statuscake_contact_group": resourceStatusCakeContactGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := scapi.NewAPIClient(d.Get("apikey").(string))

	return client, nil
}
