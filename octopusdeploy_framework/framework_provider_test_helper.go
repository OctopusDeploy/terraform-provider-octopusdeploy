package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type testOctopusDeployFrameworkProvider struct {
	*octopusDeployFrameworkProvider
	featureToggleOverrides map[string]bool
}

var _ provider.Provider = (*testOctopusDeployFrameworkProvider)(nil)

func newTestOctopusDeployFrameworkProvider(overrides map[string]bool) *testOctopusDeployFrameworkProvider {
	return &testOctopusDeployFrameworkProvider{
		octopusDeployFrameworkProvider: NewOctopusDeployFrameworkProvider(),
		featureToggleOverrides:         overrides,
	}
}

func (p *testOctopusDeployFrameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	p.octopusDeployFrameworkProvider.Metadata(ctx, req, resp)
}

func (p *testOctopusDeployFrameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	p.octopusDeployFrameworkProvider.Schema(ctx, req, resp)
}

func (p *testOctopusDeployFrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.octopusDeployFrameworkProvider.Configure(ctx, req, resp)

	if p.featureToggleOverrides != nil && !resp.Diagnostics.HasError() {
		config, ok := resp.ResourceData.(*Config)
		if !ok {
			resp.Diagnostics.AddError("Test configuration error", "Failed to cast provider data to Config")
			return
		}

		if config.FeatureToggles == nil {
			config.FeatureToggles = make(map[string]bool)
		}

		for key, value := range p.featureToggleOverrides {
			config.FeatureToggles[key] = value
			tflog.Debug(ctx, fmt.Sprintf("Test override: feature toggle %s = %t", key, value))
		}

		resp.ResourceData = config
		resp.DataSourceData = config
	}
}

func (p *testOctopusDeployFrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return p.octopusDeployFrameworkProvider.DataSources(ctx)
}

func (p *testOctopusDeployFrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return p.octopusDeployFrameworkProvider.Resources(ctx)
}
