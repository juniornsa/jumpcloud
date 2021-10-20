package provider

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAssociation() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource for associating a JumpCloud user group to objects like SSO applications, G Suite, Office 365, LDAP and more.",
		CreateContext: resourceUserAssociationCreate,
		ReadContext:   resourceUserAssociationRead,
		UpdateContext: resourceUserAssociationUpdate,
		DeleteContext: resourceUserAssociationDelete,
		Schema: map[string]*schema.Schema{
			"system_id": {
				Description: "The ID of the `resource_user_group` resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_id": {
				Description: "The ID of the object to associate to the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The type of the object to associate to the given group. Possible values: `active_directory`, `application`, `command`, `g_suite`, `ldap_server`, `office_365`, `policy`, `radius_server`, `system`, `system_group`.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					allowedValues := []string{
						"system",
						"user",
					}

					v := val.(string)
					if !stringInSlice(v, allowedValues) {
						errors = append(errors, fmt.Errorf("%q must be one of %q", key, allowedValues))
					}
					return
				},
			},
		},
	}
}

func modifyUserAssociation(client *jcapiv2.APIClient,
	d *schema.ResourceData, action string) diag.Diagnostics {

	payload := jcapiv2.SystemGraphManagementReq{
		Op:    action,
		Type_: d.Get("type").(string),
		Id:    d.Get("user_id").(string),
	}

	optional := map[string]interface{}{
		"body": payload,
	}

	_, err := client.GraphApi.GraphSystemAssociationsPost(
		context.TODO(), d.Get("system_id").(string), headerAccept, headerAccept, optional)

	return diag.FromErr(err)
}

func updateUserAssociation(client *jcapiv2.APIClient,
	d *schema.ResourceData, action string) diag.Diagnostics {

	payload := jcapiv2.SystemGraphManagementReq{
		Op:    action,
		Type_: d.Get("type").(string),
		Id:    d.Get("user_id").(string),
	}

	optional := map[string]interface{}{
		"body": payload,
	}

	_, err := client.GraphApi.GraphSystemAssociationsPost(
		context.TODO(), d.Get("system_id").(string), headerAccept, headerAccept, optional)

	return diag.FromErr(err)
}

func resourceUserAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	err := modifyUserAssociation(client, d, "add")
	if err != nil {
		return err
	}
	return resourceUserAssociationRead(ctx, d, meta)
}

func resourceUserAssociationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	err := updateUserAssociation(client, d, "add")
	if err != nil {
		return err
	}
	return resourceUserAssociationRead(ctx, d, meta)
}

func resourceUserAssociationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	optionals := map[string]interface{}{
		"userId": d.Get("user_id").(string),
		"limit":  int32(100),
	}

	graphconnect, _, err := client.UsersApi.GraphUserAssociationsList(
		context.TODO(), d.Get("user_id").(string), "", "", []string{"system"}, optionals)
	if err != nil {
		return diag.FromErr(err)
	}

	// the ID of the specified object is buried in a complex construct
	for _, v := range graphconnect {
		if v.To.Id == d.Get("system_id") {
			resourceId := d.Get("user_id").(string) + "/" + d.Get("system_id").(string)
			d.SetId(resourceId)
			return nil
		}
	}

	// element does not exist; unset ID
	d.SetId("")
	return nil
}

func resourceUserAssociationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)
	return modifyUserAssociation(client, d, "remove")
}
