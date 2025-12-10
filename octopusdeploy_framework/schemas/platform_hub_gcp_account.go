package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	PlatformHubGcpAccountResourceName   = "platform_hub_gcp_account"
	PlatformHubGcpAccountDatasourceName = "platform_hub_gcp_accounts"
)

type PlatformHubGcpAccountSchema struct{}

var _ EntitySchema = PlatformHubGcpAccountSchema{}

type PlatformHubGcpAccountModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	JsonKey     types.String `tfsdk:"json_key"`

	ResourceModel
}

func (a PlatformHubGcpAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Platform Hub GCP account in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        util.ResourceString().Required().Description("The name of this GCP account.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this GCP account.").Build(),
			"json_key": util.ResourceString().
				Required().
				Sensitive().
				Description("The JSON key for the GCP account.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
		},
	}
}

func (a PlatformHubGcpAccountSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Use this data source to retrieve information about Platform Hub GCP accounts in Octopus Deploy.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":   util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
			"name": util.DataSourceString().Optional().Description("The name of the GCP account to filter by.").Build(),
			"skip": util.DataSourceInt64().Optional().Description("The number of records to skip.").Build(),
			"take": util.DataSourceInt64().Optional().Description("The number of records to take.").Build(),
			"gcp_accounts": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Provides information about existing Platform Hub GCP accounts.",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetPlatformHubGcpAccountDatasourceAttributes(),
				},
			},
		},
	}
}

func GetPlatformHubGcpAccountDatasourceAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
		"name":        util.DataSourceString().Computed().Description("The name of this GCP account.").Build(),
		"description": util.DataSourceString().Computed().Description("The description of this GCP account.").Build(),
	}
}
