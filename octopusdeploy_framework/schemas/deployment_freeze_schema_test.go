package schemas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Regression test: the recurring schedule validator converts the config into a struct that
// includes utc_offset_in_minutes, so the resource schema must define that attribute.
// Without it, any recurring_schedule fails validation with a "Value Conversion Error".
func TestDeploymentFreezeRecurringScheduleValidatesAgainstResourceSchema(t *testing.T) {
	ctx := context.Background()
	recurring := DeploymentFreezeSchema{}.GetResourceSchema().Attributes["recurring_schedule"].(resourceSchema.SingleNestedAttribute)
	objType := recurring.GetType().(types.ObjectType)

	daysOfWeek, diags := types.ListValueFrom(ctx, types.StringType, []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"})
	if diags.HasError() {
		t.Fatal(diags)
	}
	values := map[string]attr.Value{
		"type":         types.StringValue("Weekly"),
		"unit":         types.Int64Value(1),
		"end_type":     types.StringValue("Never"),
		"days_of_week": daysOfWeek,
	}
	for name, attrType := range objType.AttrTypes {
		if _, set := values[name]; set {
			continue
		}
		nullValue, err := attrType.ValueFromTerraform(ctx, tftypes.NewValue(attrType.TerraformType(ctx), nil))
		if err != nil {
			t.Fatal(err)
		}
		values[name] = nullValue
	}

	config, diags := types.ObjectValue(objType.AttrTypes, values)
	if diags.HasError() {
		t.Fatal(diags)
	}

	resp := &validator.ObjectResponse{}
	NewRecurringScheduleValidator().ValidateObject(ctx, validator.ObjectRequest{ConfigValue: config}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected recurring schedule to validate, got: %v", resp.Diagnostics)
	}
}
