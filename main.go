package main

import (
	// Used for debugging

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
		ProviderAddr: "registry.terraform.io/example-namespace/example",
		ProviderFunc: func() *schema.Provider {
			return tss.Provider()
		},
	}

	plugin.Serve(opts)

	// var debugMode bool
	// flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	// flag.Parse()

	// if debugMode {

	// 	// err := plugin.Debug(context.Background(), "terraform.local/nhn/tss", opts)
	// 	// if err != nil {
	// 	// 	log.Fatal(err.Error())
	// 	// }
	// 	// return
	// 	opts := &plugin.ServeOpts{
	// 		ProviderFunc: func() *schema.Provider {
	// 			return tss.Provider()
	// 		},
	// 		Debug: true,
	// 	}
	// 	plugin.Serve(opts)
	// } else {
	// 	opts := &plugin.ServeOpts{
	// 		ProviderFunc: func() *schema.Provider {
	// 			return tss.Provider()
	// 		},
	// 	}
	// 	plugin.Serve(opts)

	// }

}
