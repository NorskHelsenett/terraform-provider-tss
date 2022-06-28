package tss

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thycotic/tss-sdk-go/server"
)

func providerConfig(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return server.Configuration{
		ServerURL: d.Get("server_url").(string),
		Credentials: server.UserCredential{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
			Domain:   d.Get("domain").(string),
		},
	}, diags
}

// Provider is a Terraform DataSource
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"tss_secret": resourceSecret(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"tss_secret": dataSourceSecret(),
		},
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Secret Server base URL e.g. https://localhost/SecretServer",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the Secret Server User to connect as",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The password of the Secret Server User",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain of the Secret Server User",
			},
		},
		ConfigureContextFunc: providerConfig,
	}
}
