package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SpaceDefaultLifecycleReleaseRetentionPolicySchema struct{}

// GetDatasourceSchema implements EntitySchema.
func (s SpaceDefaultLifecycleReleaseRetentionPolicySchema) GetDatasourceSchema() ds.Schema {
	return ds.Schema{}
}

var _ EntitySchema = SpaceDefaultLifecycleReleaseRetentionPolicySchema{}

func (s SpaceDefaultLifecycleReleaseRetentionPolicySchema) GetResourceSchema() rs.Schema {
	return rs.Schema{
		Description: "Manages a space's default lifecycle release retention policy.",
		Attributes: map[string]rs.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema("space default retention policy"),
			"strategy": util.ResourceString().Description("The strategy for the retention policy.").Required().Validators(stringvalidator.OneOf("Forever", "Count")).Build(),
			// Optional
			"quantity_to_keep": util.ResourceInt64().Description("The quantity of items to keep.").Validators(int64validator.AtLeast(1)).Optional().Build(),
			"unit":             util.ResourceString().Description("The unit of time for the retention policy.").Optional().Validators(stringvalidator.OneOf("Days", "Items")).Build(),
		},
	}
}

type SpaceDefaultLifecycleReleaseRetentionPoliciesResourceModel struct {
	ID             types.String `tfsdk:"id"`
	SpaceID        types.String `tfsdk:"space_id"`
	Strategy       types.String `tfsdk:"strategy"`
	QuantityToKeep types.Int64  `tfsdk:"quantity_to_keep"`
	Unit           types.String `tfsdk:"unit"`
}
