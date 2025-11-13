package schemas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
