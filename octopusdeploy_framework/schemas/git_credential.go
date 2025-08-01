package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	GitCredentialResourceName   = "git_credential"
	GitCredentialDatasourceName = "git_credentials"
)

type GitCredentialSchema struct{}

var _ EntitySchema = GitCredentialSchema{}

func (g GitCredentialSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Git credential in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"space_id":    util.ResourceString().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Description("The space ID associated with this Git Credential.").Build(),
			"name":        util.ResourceString().Required().Description("The name of this Git Credential.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this Git Credential.").Build(),
			"type": util.ResourceString().
				Optional().
				Description("The Git credential authentication type.").
				Build(),
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

func gitCredentialRepositoryRestrictionAttribute() resourceSchema.SingleNestedAttribute {
	return resourceSchema.SingleNestedAttribute{
		Description: "Sets the repository restrictions associated with the Git credential.",
		Attributes: map[string]resourceSchema.Attribute{
			"enabled": util.ResourceBool().
				Description("Whether repository restrictions are enabled.").
				Required().
				Build(),
			"allowed_repositories": util.ResourceSet(types.StringType).
				Description("Set of allowed repository URL's.").
				Required().
				Build(),
		},
		Optional: true,
		Computed: true,
		Default: objectdefault.StaticValue(
			types.ObjectValueMust(
				map[string]attr.Type{
					"enabled":              types.BoolType,
					"allowed_repositories": types.SetType{ElemType: types.StringType},
				},
				map[string]attr.Value{
					"enabled":              types.BoolValue(false),
					"allowed_repositories": types.SetValueMust(types.StringType, make([]attr.Value, 0)),
				},
			),
		),
	}
}

func (g GitCredentialSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Use this data source to retrieve information about Git credentials in Octopus Deploy.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":       util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
			"space_id": util.DataSourceString().Optional().Description("The space ID associated with this Git Credential.").Build(),
			"name":     util.DataSourceString().Optional().Description("The name of the Git Credential to filter by.").Build(),
			"skip":     util.DataSourceInt64().Optional().Description("The number of records to skip.").Build(),
			"take":     util.DataSourceInt64().Optional().Description("The number of records to take.").Build(),
			"git_credentials": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Provides information about existing GitCredentials.",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetGitCredentialDatasourceAttributes(),
				},
			},
		},
	}
}

func GetGitCredentialDatasourceAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                      util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
		"space_id":                util.DataSourceString().Computed().Description("The space ID associated with this Git Credential.").Build(),
		"name":                    util.DataSourceString().Computed().Description("The name of this Git Credential.").Build(),
		"description":             util.DataSourceString().Computed().Description("The description of this Git Credential.").Build(),
		"type":                    util.DataSourceString().Computed().Description("The Git credential authentication type.").Build(),
		"username":                util.DataSourceString().Computed().Description("The username for the Git credential.").Build(),
		"repository_restrictions": gitCredentialRepositoryRestrictionDataSourceAttribute(),
	}
}

func gitCredentialRepositoryRestrictionDataSourceAttribute() datasourceSchema.SingleNestedAttribute {
	return datasourceSchema.SingleNestedAttribute{
		Description: "Sets the repository restrictions associated with the Git credential.",
		Attributes: map[string]datasourceSchema.Attribute{
			"enabled": util.ResourceBool().
				Description("Whether repository restrictions are enabled.").
				Computed().
				Build(),
			"allowed_repositories": util.ResourceSet(types.StringType).
				Description("Set of allowed repository URL's.").
				Computed().
				Build(),
		},
		Computed: true,
	}
}
