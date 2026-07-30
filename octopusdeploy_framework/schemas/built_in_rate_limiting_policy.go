package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BuiltInRateLimitingPolicySchema struct{}

func (s BuiltInRateLimitingPolicySchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages the settings of a built-in rate limiting policy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"slug": util.ResourceString().
				Description("The slug identifying the built-in policy to manage. One of `anon` (unauthenticated), " +
					"`user` (authenticated human), or `agent` (authenticated agent).").
				Required().
				Validators(stringvalidator.OneOf("anon", "user", "agent")).
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"is_enabled": util.ResourceBool().
				Description("Whether rate limiting is enforced for this policy.").
				Required().
				Build(),
			"requests_per_minute": util.ResourceInt64().
				Description("The request rate allowed per minute.").
				Required().
				Validators(int64validator.AtLeast(1)).
				Build(),
			"burst_limit": util.ResourceInt64().
				Description("The maximum burst capacity.").
				Required().
				Validators(int64validator.AtLeast(1)).
				Build(),
			"audit_mode": util.ResourceBool().
				Description("When enabled, requests that exceed the limit are logged but not rejected.").
				Required().
				Build(),

			// Computed properties
			"name":       util.ResourceString().Description("The name of the policy.").Computed().Build(),
			"scope_type": util.ResourceString().Description("The scope this policy applies to.").Computed().Build(),
		},
	}
}

type BuiltInRateLimitingPolicyResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Slug              types.String `tfsdk:"slug"`
	IsEnabled         types.Bool   `tfsdk:"is_enabled"`
	RequestsPerMinute types.Int64  `tfsdk:"requests_per_minute"`
	BurstLimit        types.Int64  `tfsdk:"burst_limit"`
	AuditMode         types.Bool   `tfsdk:"audit_mode"`
	Name              types.String `tfsdk:"name"`
	ScopeType         types.String `tfsdk:"scope_type"`
}
