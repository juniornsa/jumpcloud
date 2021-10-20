package provider

import (
	"context"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJumpCloudSystemDirectory() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about the JumpCloud System.",
		ReadContext: dataSourceJumpCloudSystemDirectoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "The system defined name, e.g. `xxxx-xxx-Mac`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceJumpCloudSystemDirectoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configv1 := convertV2toV1Config(meta.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	filterFunction := make(map[string]interface{}, 0)

	var m interface{}
	m = map[string]interface{}{
		"searchTerm": d.Get("name").(string),
		"fields":     []string{"displayName"},
	}

	s := jcapiv1.Search{
		SearchFilter: &m,
	}
	//s.SearchFilter = &m
	filterFunction = map[string]interface{}{
		"body": s,
	}
	//filterFunction["body"] = filter
	search, _, err := client.SearchApi.SearchSystemsPost(context.TODO(), "application/json", "accept", filterFunction)
	if err != nil {
		return diag.Errorf("could not find system specified. Previous error message: %v", err)
	}
	for _, v := range search.Results {
		if v.DisplayName == d.Get("name").(string) {
			resourceId := v.Id
			d.SetId(resourceId)
			_ = d.Set("id", v.Id)
			_ = d.Set("os", v.Os)
			return nil
		}
	}
	// indicates that everything went well
	return nil
}
