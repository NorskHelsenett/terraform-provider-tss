package tss

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vidarno/tss-sdk-go/v2/server"
)

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(int)
	field := d.Get("field").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	secrets, err := server.New(meta.(server.Configuration))

	if err != nil {
		log.Printf("[DEBUG] configuration error: %s", err)
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] getting secret with id %d", id)

	secret, err := secrets.Secret(id)

	if err != nil {
		log.Print("[DEBUG] unable to get secret", err)
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(secret.ID))

	log.Printf("[DEBUG] using '%s' field of secret with id %d", field, id)

	if value, ok := secret.Field(field); ok {
		d.Set("value", value)
		return nil
	}
	diags = append(diags, diag.Errorf("the secret does not contain a '%s' field", field)...)
	return diags
}

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretRead,

		Schema: map[string]*schema.Schema{
			"value": {
				Computed:    true,
				Description: "the value of the field of the secret",
				Sensitive:   true,
				Type:        schema.TypeString,
			},
			"field": {
				Description: "the field to extract from the secret",
				Required:    true,
				Type:        schema.TypeString,
			},
			"id": {
				Description: "the id of the secret",
				Required:    true,
				Type:        schema.TypeInt,
			},
		},
	}
}
