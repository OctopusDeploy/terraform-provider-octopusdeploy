package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	PlatformHubAzureOidcAccountResourceName   = "platform_hub_azure_oidc_account"
	PlatformHubAzureOidcAccountDatasourceName = "platform_hub_azure_oidc_accounts"
)

type PlatformHubAzureOidcAccountSchema struct{}

var _ EntitySchema = PlatformHubAzureOidcAccountSchema{}

type PlatformHubAzureOidcAccountModel struct {
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	SubscriptionID             types.String `tfsdk:"subscription_id"`
	ApplicationID              types.String `tfsdk:"application_id"`
	TenantID                   types.String `tfsdk:"tenant_id"`
	ExecutionSubjectKeys       types.Set    `tfsdk:"execution_subject_keys"`
	HealthSubjectKeys          types.Set    `tfsdk:"health_subject_keys"`
	AccountTestSubjectKeys     types.Set    `tfsdk:"account_test_subject_keys"`
	Audience                   types.String `tfsdk:"audience"`
	AzureEnvironment           types.String `tfsdk:"azure_environment"`
	AuthenticationEndpoint     types.String `tfsdk:"authentication_endpoint"`
	ResourceManagementEndpoint types.String `tfsdk:"resource_management_endpoint"`

	ResourceModel
}

func (a PlatformHubAzureOidcAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Platform Hub Azure OpenID Connect account in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        util.ResourceString().Required().Description("The name of this Azure OpenID Connect account.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this Azure OpenID Connect account.").Build(),
			"subscription_id": util.ResourceString().
				Required().
				Description("The Azure subscription ID.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"application_id": util.ResourceString().
				Required().
				Description("The Azure application ID (client ID).").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"tenant_id": util.ResourceString().
				Required().
				Description("The Azure tenant ID.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"execution_subject_keys": util.ResourceSet(types.StringType).
				Optional().
				Computed().
				Description("Keys to include in a deployment or runbook. Valid options are `space`, `environment`, `project`, `tenant`, `runbook`, `account`, `type`.").
				Build(),
			"health_subject_keys": util.ResourceSet(types.StringType).
				Optional().
				Computed().
				Description("Keys to include in a health check. Valid options are `space`, `account`, `target`, `type`.").
				Build(),
			"account_test_subject_keys": util.ResourceSet(types.StringType).
				Optional().
				Computed().
				Description("Keys to include in an account test. Valid options are `space`, `account`, `type`.").
				Build(),
			"audience": util.ResourceString().
				Optional().
				Computed().
				Description("The audience for the Azure OIDC account.").
				Build(),
			"azure_environment": util.ResourceString().
				Optional().
				Computed().
				Description("The Azure environment. Valid values are `AzureCloud`, `AzureChinaCloud`, `AzureGermanCloud`, `AzureUSGovernment`.").
				Build(),
			"authentication_endpoint": util.ResourceString().
				Optional().
				Computed().
				Description("The Active Directory endpoint base URI.").
				Build(),
			"resource_management_endpoint": util.ResourceString().
				Optional().
				Computed().
				Description("The Azure Resource Management endpoint base URI.").
				Build(),
		},
	}
}

func (a PlatformHubAzureOidcAccountSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Use this data source to retrieve information about Platform Hub Azure OpenID Connect accounts in Octopus Deploy.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":   util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
			"name": util.DataSourceString().Optional().Description("The name of the Azure OpenID Connect account to filter by.").Build(),
			"skip": util.DataSourceInt64().Optional().Description("The number of records to skip.").Build(),
			"take": util.DataSourceInt64().Optional().Description("The number of records to take.").Build(),
			"azure_oidc_accounts": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Provides information about existing Platform Hub Azure OpenID Connect accounts.",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetPlatformHubAzureOidcAccountDatasourceAttributes(),
				},
			},
		},
	}
}

func GetPlatformHubAzureOidcAccountDatasourceAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                           util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
		"name":                         util.DataSourceString().Computed().Description("The name of this Azure OpenID Connect account.").Build(),
		"description":                  util.DataSourceString().Computed().Description("The description of this Azure OpenID Connect account.").Build(),
		"subscription_id":              util.DataSourceString().Computed().Description("The Azure subscription ID.").Build(),
		"application_id":               util.DataSourceString().Computed().Description("The Azure application ID (client ID).").Build(),
		"tenant_id":                    util.DataSourceString().Computed().Description("The Azure tenant ID.").Build(),
		"execution_subject_keys":       util.DataSourceSet(types.StringType).Computed().Description("Keys to include in a deployment or runbook.").Build(),
		"health_subject_keys":          util.DataSourceSet(types.StringType).Computed().Description("Keys to include in a health check.").Build(),
		"account_test_subject_keys":    util.DataSourceSet(types.StringType).Computed().Description("Keys to include in an account test.").Build(),
		"audience":                     util.DataSourceString().Computed().Description("The audience for the Azure OIDC account.").Build(),
		"azure_environment":            util.DataSourceString().Computed().Description("The Azure environment.").Build(),
		"authentication_endpoint":      util.DataSourceString().Computed().Description("The Active Directory endpoint base URI.").Build(),
		"resource_management_endpoint": util.DataSourceString().Computed().Description("The Azure Resource Management endpoint base URI.").Build(),
	}
}
