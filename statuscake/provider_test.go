package statuscake_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/StatusCakeDev/terraform-provider-statuscake/statuscake"
)

var testProviders = map[string]terraform.ResourceProvider{
	"statuscake": statuscake.Provider(),
}

func TestProvider(t *testing.T) {
	if err := testProviders["statuscake"].(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
