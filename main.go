package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/NorskHelsenett/terraform-provider-tss/tss"
)

func main() {

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: "terraform.local/nhn/tss",
		ProviderFunc: func() *schema.Provider {
			return tss.Provider()
		},
	}

	plugin.Serve(opts)

}
