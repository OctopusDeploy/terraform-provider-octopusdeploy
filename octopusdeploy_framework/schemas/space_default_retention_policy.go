package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetSpaceDefaultRetentionPolicyAttributes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"space_id":         types.StringType,
		"retention_type":   types.StringType,
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
	}
}

type SpaceDefaultRetentionPolicySchema struct{}

var _ EntitySchema = SpaceDefaultRetentionPolicySchema{}

func (s SpaceDefaultRetentionPolicySchema) GetDatasourceSchema() ds.Schema {
	return ds.Schema{
		Description: "Manages a space's default retention policy.",
		Attributes: map[string]ds.Attribute{
			// request
			"space_id":       GetSpaceIdDatasourceSchema("space default retention policy", false),
			"retention_type": util.ResourceString().Description("The type of retention policy.").Required().Validators(stringvalidator.OneOf("LifecycleRelease", "LifecycleTentacle")).Build(),

			// response
			"id":               util.ResourceString().Description("The ID of the retention policy.").Computed().Build(),
			"strategy":         util.ResourceString().Description("The strategy for the retention policy.").Computed().Build(),
			"quantity_to_keep": util.ResourceInt64().Description("The quantity of items to keep.").Computed().Optional().Build(),
			"unit":             util.ResourceString().Description("The unit of time for the retention policy.").Computed().Optional().Build(),
		},
	}
}

func (s SpaceDefaultRetentionPolicySchema) GetResourceSchema() rs.Schema {
	return rs.Schema{}
}

type SpaceDefaultRetentionPoliciesDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	SpaceID        types.String `tfsdk:"space_id"`
	RetentionType  types.String `tfsdk:"retention_type"`
	Strategy       types.String `tfsdk:"strategy"`
	QuantityToKeep types.Int64  `tfsdk:"quantity_to_keep"`
	Unit           types.String `tfsdk:"unit"`
}
