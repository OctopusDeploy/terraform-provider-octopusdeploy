package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	PlatformHubAwsAccountResourceName   = "platform_hub_aws_account"
	PlatformHubAwsAccountDatasourceName = "platform_hub_aws_accounts"
)

type PlatformHubAwsAccountSchema struct{}

var _ EntitySchema = PlatformHubAwsAccountSchema{}

type PlatformHubAwsAccountModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	AccessKey   types.String `tfsdk:"access_key"`
	SecretKey   types.String `tfsdk:"secret_key"`

	ResourceModel
}

func (a PlatformHubAwsAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Platform Hub AWS account in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        util.ResourceString().Required().Description("The name of this AWS account.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this AWS account.").Build(),
			"access_key": util.ResourceString().
				Required().
				Description("The access key for the AWS account.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"secret_key": util.ResourceString().
				Required().
				Sensitive().
				Description("The secret key for the AWS account.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
		},
	}
}

func (a PlatformHubAwsAccountSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Use this data source to retrieve information about Platform Hub AWS accounts in Octopus Deploy.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":   util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
			"name": util.DataSourceString().Optional().Description("The name of the AWS account to filter by.").Build(),
			"skip": util.DataSourceInt64().Optional().Description("The number of records to skip.").Build(),
			"take": util.DataSourceInt64().Optional().Description("The number of records to take.").Build(),
			"aws_accounts": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "Provides information about existing Platform Hub AWS accounts.",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetPlatformHubAwsAccountDatasourceAttributes(),
				},
			},
		},
	}
}

func GetPlatformHubAwsAccountDatasourceAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
		"name":        util.DataSourceString().Computed().Description("The name of this AWS account.").Build(),
		"description": util.DataSourceString().Computed().Description("The description of this AWS account.").Build(),
		"access_key":  util.DataSourceString().Computed().Description("The access key for the AWS account.").Build(),
	}
}
