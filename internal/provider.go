package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MyHeritageProvider struct {
	apiKey types.String
}

func New() provider.Provider {
	return &MyHeritageProvider{}
}

func (p *MyHeritageProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "myheritage"
}

func (p *MyHeritageProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "The API key for the MyHeritage API",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

type MyHeritageProviderConfig struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *MyHeritageProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config MyHeritageProviderConfig

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p.apiKey = config.ApiKey
	resp.ResourceData = p
}

func (p *MyHeritageProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProfileResource,
	}
}

func (p *MyHeritageProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
