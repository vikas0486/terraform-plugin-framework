package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerSchema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ThalesProvider struct {
	client *APIClient
}

func New() provider.Provider {
	return &ThalesProvider{}
}

func (p *ThalesProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "thales"
}

func (p *ThalesProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {

	resp.Schema = providerSchema.Schema{
		Attributes: map[string]providerSchema.Attribute{

			"endpoint": providerSchema.StringAttribute{
				Required:    true,
				Description: "Thales REST API Endpoint",
			},
		},
	}
}

func (p *ThalesProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {

	var config ProviderConfig

	resp.Diagnostics.Append(
		req.Config.Get(ctx, &config)...,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	p.client = NewClient(config.Endpoint.ValueString())

	resp.ResourceData = p.client
}

func (p *ThalesProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {

	return []func() resource.Resource{
		NewKeystoreResource,
	}
}

func (p *ThalesProvider) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {

	return nil
}
