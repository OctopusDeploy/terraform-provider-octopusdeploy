package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	PlatformHubUsernamePasswordAccountResourceName   = "platform_hub_username_password_account"
	PlatformHubUsernamePasswordAccountDatasourceName = "platform_hub_username_password_accounts"
)

type PlatformHubUsernamePasswordAccountSchema struct{}

var _ EntitySchema = PlatformHubUsernamePasswordAccountSchema{}

type PlatformHubUsernamePasswordAccountModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`

	ResourceModel
}

func (a PlatformHubUsernamePasswordAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Platform Hub Username-Password account in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        util.ResourceString().Required().Description("The name of this Username-Password account.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this Username-Password account.").Build(),
			"username": util.ResourceString().
				Required().
				Description("The username for the account.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"password": util.ResourceString().
				Required().
				Sensitive().
				Description("The password for the account.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
		},
	}
}

func (a PlatformHubUsernamePasswordAccountSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Use this data source to retrieve information about Platform Hub Username-Password accounts in Octopus Deploy.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":   util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
			"name": util.DataSourceString().Optional().Description("The name of the Username-Password account to filter by.").Build(),
			"skip": util.DataSourceInt64().Optional().Description("The number of records to skip.").Build(),
			"take": util.DataSourceInt64().Optional().Description("The number of records to take.").Build(),
			"username_password_accounts": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Provides information about existing Platform Hub Username-Password accounts.",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetPlatformHubUsernamePasswordAccountDatasourceAttributes(),
				},
			},
		},
	}
}

func GetPlatformHubUsernamePasswordAccountDatasourceAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
		"name":        util.DataSourceString().Computed().Description("The name of this Username-Password account.").Build(),
		"description": util.DataSourceString().Computed().Description("The description of this Username-Password account.").Build(),
		"username":    util.DataSourceString().Computed().Description("The username for the account.").Build(),
	}
}
