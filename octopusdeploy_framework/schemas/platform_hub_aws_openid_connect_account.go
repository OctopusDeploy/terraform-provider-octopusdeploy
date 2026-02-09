package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	PlatformHubAwsOpenIDConnectAccountResourceName = "platform_hub_aws_openid_connect_account"
)

type PlatformHubAwsOpenIDConnectAccountSchema struct{}

type PlatformHubAwsOpenIDConnectAccountModel struct {
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	RoleArn                types.String `tfsdk:"role_arn"`
	SessionDuration        types.String `tfsdk:"session_duration"`
	ExecutionSubjectKeys   types.Set    `tfsdk:"execution_subject_keys"`
	HealthSubjectKeys      types.Set    `tfsdk:"health_subject_keys"`
	AccountTestSubjectKeys types.Set    `tfsdk:"account_test_subject_keys"`

	ResourceModel
}

func (a PlatformHubAwsOpenIDConnectAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Platform Hub AWS OpenID Connect account in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        util.ResourceString().Required().Description("The name of this Platform Hub AWS OIDC account.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("A user-friendly description of this Platform Hub AWS OIDC account.").Build(),
			"role_arn": util.ResourceString().
				Required().
				Description("The Amazon Resource Name (ARN) of the role that the caller is assuming.").
				Validators(stringvalidator.LengthAtLeast(1)).
				Build(),
			"session_duration": util.ResourceString().
				Optional().
				Computed().
				Default("3600").
				Description("The duration, in seconds, of the role session.").
				Build(),
			"execution_subject_keys": util.ResourceSet(types.StringType).
				Optional().
				Description("Keys to include in a deployment or runbook. Valid options are `space`, `environment`, `project`, `tenant`, `runbook`, `account`, `type`.").
				Build(),
			"health_subject_keys": util.ResourceSet(types.StringType).
				Optional().
				Description("Keys to include in a health check. Valid options are `space`, `account`, `target`, `type`.").
				Build(),
			"account_test_subject_keys": util.ResourceSet(types.StringType).
				Optional().
				Description("Keys to include in an account test. Valid options are `space`, `account`, `type`.").
				Build(),
		},
	}
}
