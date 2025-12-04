package schemas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Int64 = &strategyAttributeValidator{}
var _ validator.String = &strategyAttributeValidator{}

type strategyAttributeValidator struct {
	strategy string
}

func NewStrategyAttributeValidator(forStrategy string) *strategyAttributeValidator {
	return &strategyAttributeValidator{
		strategy: forStrategy,
	}
}

// ValidateString implements validator.String.
func (r *strategyAttributeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	matchedPaths, diags := req.Config.PathMatches(ctx, req.PathExpression.AtParent().AtName("strategy"))
	resp.Diagnostics.Append(diags...)

	var strategy attr.Value = nil
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

func (r *strategyAttributeValidator) Description(context.Context) string {
	return fmt.Sprintf("This field is required when the strategy is set to '%s' and must be omitted when the strategy is not '%s'.", r.strategy, r.strategy)
}

func (r *strategyAttributeValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("This field is required when the strategy is set to '%s' and must be omitted when the strategy is not '%s'.", r.strategy, r.strategy)
}

func (r *strategyAttributeValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	matchedPaths, diags := req.Config.PathMatches(ctx, req.PathExpression.AtParent().AtName("strategy"))
	resp.Diagnostics.Append(diags...)

	var strategy attr.Value = nil
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

	if strategy.Equal(types.StringValue(r.strategy)) && (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Field",
			fmt.Sprintf("%s is required when the strategy is '%s'.", req.Path, r.strategy),
		)
		return
	}
	if !strategy.Equal(types.StringValue(r.strategy)) && !req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Field",
			fmt.Sprintf("%s must not be set when the strategy is '%s'.", req.Path, strategy),
		)
		return
	}
}
