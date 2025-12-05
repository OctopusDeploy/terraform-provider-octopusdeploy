package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SpaceDefaultLifecycleTentacleRetentionPolicySchema struct{}

// GetDatasourceSchema implements EntitySchema.
func (s SpaceDefaultLifecycleTentacleRetentionPolicySchema) GetDatasourceSchema() ds.Schema {
	return ds.Schema{
		Description: "Manages a space's default retention policy for how files on tentacles are retained.",
		Attributes: map[string]ds.Attribute{
			// request
			"space_id": ds.StringAttribute{
				Description: "The ID of the space.",
				Required:    true,
			},

			// response
			"id": util.ResourceString().Description("The ID of the retention policy.").Computed().Build(),
			"strategy": util.ResourceString().Description("How retention will be set. Valid strategies are `Forever`, and `Count`." +
				"\n  - `strategy = \"Forever\"`, is used if files on tentacles should never be deleted." +
				"\n  - `strategy = \"Count\"`, is used if a specific number of days/releases files on tentacles should be kept for.").Computed().Build(),
			"quantity_to_keep": util.ResourceInt64().Description("The number of days/releases to keep files on tentacles.").Computed().Build(),
			"unit":             util.ResourceString().Description("The unit of quantity to keep. Valid Units are `Days` or `Items`").Computed().Build(),
		},
	}
}

var _ EntitySchema = SpaceDefaultLifecycleTentacleRetentionPolicySchema{}

func (s SpaceDefaultLifecycleTentacleRetentionPolicySchema) GetResourceSchema() rs.Schema {
	return rs.Schema{
		Description: "Manages a space's default retention policy for how files on tentacles are retained.",
		Attributes: map[string]rs.Attribute{
			"id": GetIdResourceSchema(),
			"space_id": rs.StringAttribute{
				Description: "The ID of the space.",
				Required:    true,
			},
			"strategy": util.ResourceString().Description("How retention will be set. Valid strategies are `Forever`, and `Count`." +
				"\n  - `strategy = \"Forever\"`, is used if files on tentacles should never be deleted." +
				"\n  - `strategy = \"Count\"`, is used if a specific number of days/releases files on tentacles should be kept for.").Required().Validators(stringvalidator.OneOf("Forever", "Count")).Build(),
			// Optional
			"quantity_to_keep": rs.Int64Attribute{
				Description: "The number of days/releases to keep files on tentacles.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					NewStrategyAttributeValidator("Count"),
				},
				Optional: true,
			},
			"unit": rs.StringAttribute{
				Description: "The unit of quantity to keep. Valid Units are `Days` or `Items`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Days", "Items"),
					NewStrategyAttributeValidator("Count"),
				},
				Optional: true,
			},
		},
	}
}

type SpaceDefaultLifecycleTentacleRetentionPoliciesResourceModel struct {
	ID             types.String `tfsdk:"id"`
	SpaceID        types.String `tfsdk:"space_id"`
	Strategy       types.String `tfsdk:"strategy"`
	QuantityToKeep types.Int64  `tfsdk:"quantity_to_keep"`
	Unit           types.String `tfsdk:"unit"`
}

type SpaceDefaultLifecycleTentacleRetentionPoliciesDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	SpaceID        types.String `tfsdk:"space_id"`
	Strategy       types.String `tfsdk:"strategy"`
	QuantityToKeep types.Int64  `tfsdk:"quantity_to_keep"`
	Unit           types.String `tfsdk:"unit"`
}
