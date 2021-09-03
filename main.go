package main

import (
	"context"
	"os"

	"github.com/StatusCakeDev/terraform-provider-statuscake/statuscake"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	_, debug := os.LookupEnv("DEBUG")

	opts := &plugin.ServeOpts{ProviderFunc: statuscake.Provider}

	if debug {
		// Run the provider with support for debuggers like delve.
		if err := plugin.Debug(context.Background(), "registry.terraform.io/StatusCakeDev/statuscake", opts); err != nil {
			panic(err)
		}

		return
	}

	plugin.Serve(opts)
}
