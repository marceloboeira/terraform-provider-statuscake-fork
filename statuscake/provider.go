package statuscake

import (
	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_USERNAME", nil),
				Description: "Username for StatusCake Account.",
			},
			"apikey": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_APIKEY", nil),
				Description: "API Key for StatusCake",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"statuscake_test":          resourceStatusCakeTest(),
			"statuscake_contact_group": resourceStatusCakeContactGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	auth := statuscake.Auth{
		Username: d.Get("username").(string),
		Apikey:   d.Get("apikey").(string),
	}
	return statuscake.New(auth)
}
