package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	PlatformHubGenericOidcAccountResourceName   = "platform_hub_generic_oidc_account"
	PlatformHubGenericOidcAccountDatasourceName = "platform_hub_generic_oidc_accounts"
)

type PlatformHubGenericOidcAccountSchema struct{}

var _ EntitySchema = PlatformHubGenericOidcAccountSchema{}

type PlatformHubGenericOidcAccountModel struct {
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	ExecutionSubjectKeys types.List   `tfsdk:"execution_subject_keys"`
	Audience             types.String `tfsdk:"audience"`

	ResourceModel
}

func (a PlatformHubGenericOidcAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Platform Hub Generic OIDC account in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        util.ResourceString().Required().Description("The name of this Generic OIDC account.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this Generic OIDC account.").Build(),
			"execution_subject_keys": util.ResourceList(types.StringType).
				Optional().
				Computed().
				Description("Keys to include in a deployment or runbook. Valid options are `space`, `environment`, `project`, `tenant`, `runbook`, `account`, `type`.").
				Build(),
			"audience": util.ResourceString().
				Optional().
				Computed().
				Description("The audience associated with this account.").
				Build(),
		},
	}
}

func (a PlatformHubGenericOidcAccountSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Use this data source to retrieve information about Platform Hub Generic OIDC accounts in Octopus Deploy.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":   util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
			"name": util.DataSourceString().Optional().Description("The name of the Generic OIDC account to filter by.").Build(),
			"skip": util.DataSourceInt64().Optional().Description("The number of records to skip.").Build(),
			"take": util.DataSourceInt64().Optional().Description("The number of records to take.").Build(),
			"generic_oidc_accounts": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Provides information about existing Platform Hub Generic OIDC accounts.",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetPlatformHubGenericOidcAccountDatasourceAttributes(),
				},
			},
		},
	}
}

func GetPlatformHubGenericOidcAccountDatasourceAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                     util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
		"name":                   util.DataSourceString().Computed().Description("The name of this Generic OIDC account.").Build(),
		"description":            util.DataSourceString().Computed().Description("The description of this Generic OIDC account.").Build(),
		"execution_subject_keys": util.DataSourceList(types.StringType).Computed().Description("Keys to include in a deployment or runbook.").Build(),
		"audience":               util.DataSourceString().Computed().Description("The audience associated with this account.").Build(),
	}
}
