package config

import "github.com/hashicorp/terraform-plugin-framework/types"

type MyHeritageProviderConfig struct {
	ApiKey types.String `tfsdk:"api_key"`
}
