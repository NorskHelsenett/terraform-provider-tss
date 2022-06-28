package tss

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thycotic/tss-sdk-go/server"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"folder_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"all_fields": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_id": {
							Type:     schema.TypeInt,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// If field_id is missing we assume we can find it in a Create or Update
								// so we return true to suppress a diff only on the field_id
								return new == "0"
							},
						},
						"field_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"field_value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}


func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "At beginning of function resourceSecretCreate")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	tss, err := server.New(m.(server.Configuration))
	if err != nil {
		return diag.FromErr(err)
	}

	secretModel := new(server.Secret)
	secretModel.Name, _ = d.Get("name").(string)
	secretModel.SiteID, _ = d.Get("site_id").(int)
	secretModel.FolderID, _ = d.Get("folder_id").(int)
	secretModel.SecretTemplateID = d.Get("secret_template_id").(int)

	fields := d.Get("fields").([]interface{})
	secretModel.Fields = make([]server.SecretField, len(fields))
	for i, field := range fields {
		fieldMap := field.(map[string]interface{})

		secretModel.Fields[i].FieldName, _ = fieldMap["field_name"].(string)

		if fieldMap["field_id"] == 0 {
			// If fieldID is missing we find it from fieldName
			// This requires "View Secret Templates"-permission
			refSecretTemplate, err := tss.SecretTemplate(secretModel.SecretTemplateID)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, field := range refSecretTemplate.Fields {
				if field.Name == secretModel.Fields[i].FieldName {
					secretModel.Fields[i].FieldID = field.SecretTemplateFieldID
					break
				}
			}

		} else {
			secretModel.Fields[i].FieldID, _ = fieldMap["field_id"].(int)
		}

		secretModel.Fields[i].ItemValue, _ = fieldMap["field_value"].(string)
	}

	newSecret, err := tss.CreateSecret(*secretModel)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Secret resource",
			Detail:   "Unable to create a Secret resource with provided attributes",
		})
		return diags
	}

	d.SetId(strconv.Itoa(newSecret.ID))

	tflog.Debug(ctx, "In function resourceSecretCreate, going to call resourceSecretRead")

	resourceSecretRead(ctx, d, m)

	return diags
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "At beginning of function resourceSecretRead")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	secretID, _ := strconv.Atoi(d.Id())

	tss, err := server.New(m.(server.Configuration))
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := tss.Secret(secretID)
	if err != nil {
		// Unable to get secret
		return diag.FromErr(err)
	}

	d.Set("name", secret.Name)
	d.Set("site_id", secret.SiteID)
	d.Set("folder_id", secret.FolderID)
	d.Set("secret_template_id", secret.SecretTemplateID)

	//If all_fields is false (default) we only care about the specified fields, not all
	all_fields, _ := d.Get("all_fields").(bool)
	if !all_fields {
		fields := d.Get("fields").([]interface{})

		tflog.Debug(ctx, "In function resourceSecretRead, going to call flattenSomeSecretFields")
		secretFields := flattenSomeSecretFields(ctx, &fields, &secret.Fields)
		if err := d.Set("fields", secretFields); err != nil {
			tflog.Debug(ctx, "In function resourceSecretRead, after call to flattenSomeSecretFields, failed adding Fields")
			return diag.FromErr(err)
		}

	} else {
		secretFields := flattenSecretFields(ctx, &secret.Fields)
		if err := d.Set("fields", secretFields); err != nil {
			tflog.Debug(ctx, "In function resourceSecretRead, after call to flattenSecretFields, failed adding Fields")
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceSecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChange("name") || d.HasChange("secret_template_id") || d.HasChange("folder_id") || d.HasChange("all_fields") || d.HasChange("fields") {

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		tss, err := server.New(m.(server.Configuration))
		if err != nil {
			return diag.FromErr(err)
		}

		secretModel := new(server.Secret)
		secretModel.ID, _ = strconv.Atoi(d.Id())
		secretModel.Name, _ = d.Get("name").(string)
		secretModel.SiteID, _ = d.Get("site_id").(int)
		secretModel.FolderID, _ = d.Get("folder_id").(int)
		secretModel.SecretTemplateID = d.Get("secret_template_id").(int)

		fields := d.Get("fields").([]interface{})
		secretModel.Fields = make([]server.SecretField, len(fields))
		for i, field := range fields {
			fieldMap := field.(map[string]interface{})

			secretModel.Fields[i].FieldName, _ = fieldMap["field_name"].(string)

			if fieldMap["field_id"] == 0 {
				// If fieldID is missing we find it from fieldName
				// This requires "View Secret Templates"-permission
				refSecretTemplate, err := tss.SecretTemplate(secretModel.SecretTemplateID)
				if err != nil {
					return diag.FromErr(err)
				}
				for _, field := range refSecretTemplate.Fields {
					if field.Name == secretModel.Fields[i].FieldName {
						secretModel.Fields[i].FieldID = field.SecretTemplateFieldID
						break
					}
				}

			} else {
				secretModel.Fields[i].FieldID, _ = fieldMap["field_id"].(int)
			}

			secretModel.Fields[i].ItemValue, _ = fieldMap["field_value"].(string)
		}

		newSecret, err := tss.UpdateSecret(*secretModel)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update Secret resource",
				Detail:   "Unable to update a Secret resource with provided attributes",
			})
			return diags
		}

		d.SetId(strconv.Itoa(newSecret.ID))

	}

	return resourceSecretRead(ctx, d, m)
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	secretID, _ := strconv.Atoi(d.Id())

	tss, err := server.New(m.(server.Configuration))
	if err != nil {
		return diag.FromErr(err)
	}

	err = tss.DeleteSecret(secretID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func flattenSecretFields(ctx context.Context, secretFields *[]server.SecretField) []interface{} {
	if secretFields != nil {
		ois := make([]interface{}, len(*secretFields))

		for i, secretItem := range *secretFields {
			oi := make(map[string]interface{})

			oi["field_id"] = secretItem.FieldID
			oi["field_name"] = secretItem.FieldName
			oi["field_value"] = secretItem.ItemValue
			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}

func flattenSomeSecretFields(ctx context.Context, someFields *[]interface{}, secretFields *[]server.SecretField) []interface{} {
	if secretFields != nil {

		someFieldsMap := make(map[string]string)

		for _, someFieldItem := range *someFields {
			if someFieldItem == nil {
				tflog.Debug(ctx, "In flattenSomeSecretFields. someFieldItem is nil! ")
			} else {
				fieldName, ok := someFieldItem.(map[string]interface{})["field_name"].(string)
				if ok {
					someFieldsMap[fieldName] = fieldName
					tflog.Debug(ctx, "Got fieldName "+fieldName+", added to someFieldsMap")
				} else {
					tflog.Debug(ctx, "Could not get fieldName")

				}
			}
		}

		ois := make([]interface{}, len(someFieldsMap))
		realIndex := 0

		for _, secretItem := range *secretFields {
			oi := make(map[string]interface{})

			if _, ok := someFieldsMap[secretItem.FieldName]; ok {

				oi["field_id"] = secretItem.FieldID
				oi["field_name"] = secretItem.FieldName
				oi["field_value"] = secretItem.ItemValue
				ois[realIndex] = oi
				realIndex++
			} else {
				tflog.Debug(ctx, "In function flattenSomeSecretFields, did not find field with name"+secretItem.FieldName)
			}
		}

		return ois
	}

	return make([]interface{}, 0)
}
