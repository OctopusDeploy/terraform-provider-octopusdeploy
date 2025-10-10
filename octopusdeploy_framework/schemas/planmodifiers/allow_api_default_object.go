// Package planmodifiers contains custom plan modifiers for the Terraform provider.
package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AllowApiDefaultObject returns a plan modifier that handles optional+computed objects
// where the API always provides a value. When the field is removed from configuration,
// it marks the plan as unknown (allowing the API to set a default) instead of keeping
// the old state value.
func AllowApiDefaultObject() planmodifier.Object {
	return allowApiDefaultObjectModifier{}
}

type allowApiDefaultObjectModifier struct{}

func (m allowApiDefaultObjectModifier) Description(_ context.Context) string {
	return "When an optional+computed object is removed from configuration, marks it as unknown so the API can provide a default value."
}

func (m allowApiDefaultObjectModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m allowApiDefaultObjectModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// If this is a create operation (state is null), no modification needed
	if req.StateValue.IsNull() {
		return
	}

	// If the config value is null (user removed it from configuration),
	// set the plan to unknown so Terraform knows the API will provide a value.
	// This allows sending nil to the API while accepting the API's default value
	// in the response without causing an "inconsistent result" error.
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.ObjectUnknown(req.StateValue.AttributeTypes(ctx))
		return
	}

	// For all other cases, if plan is unknown (not set), use the config value
	// This handles the case where the field is present in config
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = req.ConfigValue
	}
}
