package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PlatformHubVersionControlAnonymousSettingsSchema struct{}

var _ EntitySchema = PlatformHubVersionControlAnonymousSettingsSchema{}

func (p PlatformHubVersionControlAnonymousSettingsSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

func (p PlatformHubVersionControlAnonymousSettingsSchema) GetResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "This resource manages Platform Hub version control settings with anonymous authentication in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id":             GetIdResourceSchema(),
			"url":            util.ResourceString().Required().Description("The URL of the git repository.").Build(),
			"default_branch": util.ResourceString().Required().Description("The default branch of the git repository.").Build(),
			"base_path":      util.ResourceString().Required().Description("The base path within the repository.").Build(),
		},
	}
}

type PlatformHubVersionControlAnonymousSettingsResourceModel struct {
	URL           types.String `tfsdk:"url"`
	DefaultBranch types.String `tfsdk:"default_branch"`
	BasePath      types.String `tfsdk:"base_path"`

	ResourceModel
}
