package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_KEY", nil),
				Description: "The API key for the MyHeritage API",
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"profile": resourceProfile(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)

	// Initialize a client or other structures that you need to use the API key with
	return apiKey, nil
}
