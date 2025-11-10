package schemas

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
			"quantity_to_keep": rs.Int64Attribute{
				Description: "The quantity of items to keep.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					NewCountStrategyAttributeValidator(),
				},
				Optional: true,
			},
			"unit": rs.StringAttribute{
				Description: "The unit of time for the retention policy.",
				Validators: []validator.String{
					stringvalidator.OneOf("Days", "Weeks", "Months", "Years"),
					NewCountStrategyAttributeValidator(),
				},
				Optional: true,
			},
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

var _ validator.Int64 = &countStrategyAttributeValidator{}
var _ validator.String = &countStrategyAttributeValidator{}

type countStrategyAttributeValidator struct{}

func NewCountStrategyAttributeValidator() *countStrategyAttributeValidator {
	return &countStrategyAttributeValidator{}
}

// ValidateString implements validator.String.
func (r *countStrategyAttributeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	matchedPaths, diags := req.Config.PathMatches(ctx, path.MatchRoot("strategy"))
	resp.Diagnostics.Append(diags...)

	var strategy attr.Value
	for _, p := range matchedPaths {
		d := req.Config.GetAttribute(ctx, p, &strategy)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if strategy.IsNull() || strategy.IsUnknown() {
		return
	}

	if strategy.Equal(types.StringValue("Count")) && (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Field",
			fmt.Sprintf("%s is required when the strategy is 'Count'.", req.Path),
		)
		return
	}
	if strategy.Equal(types.StringValue("Forever")) && !req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Field",
			fmt.Sprintf("%s must not be set when the strategy is 'Forever'.", req.Path),
		)
		return
	}
}

func (r *countStrategyAttributeValidator) Description(context.Context) string {
	return "This field is required when the strategy is set to 'Count' and must be omitted when the strategy is 'Forever'."
}

func (r *countStrategyAttributeValidator) MarkdownDescription(context.Context) string {
	return "This field is required when the strategy is set to 'Count' and must be omitted when the strategy is 'Forever'."
}

func (r *countStrategyAttributeValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	matchedPaths, diags := req.Config.PathMatches(ctx, path.MatchRoot("strategy"))
	resp.Diagnostics.Append(diags...)

	var strategy attr.Value
	for _, p := range matchedPaths {
		d := req.Config.GetAttribute(ctx, p, &strategy)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if strategy.IsNull() || strategy.IsUnknown() {
		return
	}

	if strategy.Equal(types.StringValue("Count")) && (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Field",
			fmt.Sprintf("%s is required when the strategy is 'Count'.", req.Path),
		)
		return
	}
	if strategy.Equal(types.StringValue("Forever")) && !req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Field",
			fmt.Sprintf("%s must not be set when the strategy is 'Forever'.", req.Path),
		)
		return
	}
}
