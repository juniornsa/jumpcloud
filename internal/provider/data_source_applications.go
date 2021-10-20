package provider

import (
	"context"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJumpCloudApplications() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about the JumpCloud Applications.",
		ReadContext: dataSourceJumpCloudApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "The application name, e.g. `aws or docusign`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"displayname": {
				Description: "The application displayName , e.g. `AWS SSO or  DocuSign`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ssourl": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceJumpCloudApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configv1 := convertV2toV1Config(meta.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	response, _, err := client.ApplicationsApi.ApplicationsList(context.TODO(), "application/json", "accept", nil)
	if err != nil {
		return diag.Errorf("could not find any application. Previous error message: %v", err)
	}
	for _, v := range response.Results {
		if v.DisplayName == d.Get("displayname").(string) || v.Name == d.Get("name").(string) {

			resourceId := v.Id
			d.SetId(resourceId)
			_ = d.Set("id", v.Id)
			_ = d.Set("ssourl", v.SsoUrl)
			_ = d.Set("displayname", v.DisplayName)
			return nil
		}
	}

	// If the object does not exist, unset the ID
	d.SetId("")
	return nil
}
