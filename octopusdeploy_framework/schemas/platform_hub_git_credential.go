package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

const (
	PlatformHubGitCredentialResourceName   = "platform_hub_git_credential"
	PlatformHubGitCredentialDatasourceName = "platform_hub_git_credentials"
)

type PlatformHubGitCredentialSchema struct{}

var _ EntitySchema = PlatformHubGitCredentialSchema{}

func (g PlatformHubGitCredentialSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Platform Hub Git credential in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        util.ResourceString().Required().Description("The name of this Git Credential.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this Git Credential.").Build(),
			"username": util.ResourceString().
				Required().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Description("The username for the Git credential.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"password": util.ResourceString().
				Required().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Sensitive().
				Description("The password for the Git credential.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"repository_restrictions": gitCredentialRepositoryRestrictionAttribute(),
		},
	}
}

func (g PlatformHubGitCredentialSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Use this data source to retrieve information about Platform Hub Git credentials in Octopus Deploy.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":   util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
			"name": util.DataSourceString().Optional().Description("The name of the Git Credential to filter by.").Build(),
			"skip": util.DataSourceInt64().Optional().Description("The number of records to skip.").Build(),
			"take": util.DataSourceInt64().Optional().Description("The number of records to take.").Build(),
			"git_credentials": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Provides information about existing PlatformHubGitCredentials.",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetPlatformHubGitCredentialDatasourceAttributes(),
				},
			},
		},
	}
}

func GetPlatformHubGitCredentialDatasourceAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                      util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
		"name":                    util.DataSourceString().Computed().Description("The name of this Git Credential.").Build(),
		"description":             util.DataSourceString().Computed().Description("The description of this Git Credential.").Build(),
		"username":                util.DataSourceString().Computed().Description("The username for the Git credential.").Build(),
		"repository_restrictions": gitCredentialRepositoryRestrictionDataSourceAttribute(),
	}
}
