package octopusdeploy_framework

import (
	"context"
	"strings"
	"testing"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateStepTemplateParameters_UnknownList(t *testing.T) {
	data := &schemas.StepTemplateTypeResourceModel{
		Parameters: types.ListUnknown(types.ObjectType{
			AttrTypes: schemas.GetStepTemplateParameterTypeAttributes(),
		}),
	}

	diags := validateStepTemplateParameters(context.Background(), data)

	if diags.HasError() {
		t.Errorf("expected no errors for unknown parameters list, got: %v", diags)
	}
}

func TestValidateStepTemplateParameters_NullList(t *testing.T) {
	data := &schemas.StepTemplateTypeResourceModel{
		Parameters: types.ListNull(types.ObjectType{
			AttrTypes: schemas.GetStepTemplateParameterTypeAttributes(),
		}),
	}

	diags := validateStepTemplateParameters(context.Background(), data)

	if diags.HasError() {
		t.Errorf("expected no errors for null parameters list, got: %v", diags)
	}
}

func TestValidateStepTemplateParameters_Valid(t *testing.T) {
	displaySettings, _ := types.MapValue(types.StringType, map[string]attr.Value{
		"Octopus.ControlType": types.StringValue("SingleLineText"),
	})

	paramObj, _ := types.ObjectValue(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":                      types.StringValue("test-id"),
			"name":                    types.StringValue("TestParam"),
			"label":                   types.StringValue("Test Parameter"),
			"help_text":               types.StringValue("Help text"),
			"default_value":           types.StringValue("default_value"),
			"default_sensitive_value": types.StringNull(),
			"display_settings":        displaySettings,
		},
	)

	paramList, _ := types.ListValue(
		types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()},
		[]attr.Value{paramObj},
	)

	data := &schemas.StepTemplateTypeResourceModel{
		Parameters: paramList,
	}

	diags := validateStepTemplateParameters(context.Background(), data)

	if diags.HasError() {
		t.Errorf("expected no errors for valid parameters, got: %v", diags)
	}
}

// stepTemplateParameterModel builds a single parameter model for validation tests.
// Pass types.StringNull() for either default to leave it unset.
func stepTemplateParameterModel(controlType string, defaultValue attr.Value, defaultSensitiveValue attr.Value) *schemas.StepTemplateTypeResourceModel {
	displaySettings, _ := types.MapValue(types.StringType, map[string]attr.Value{
		"Octopus.ControlType": types.StringValue(controlType),
	})

	paramObj, _ := types.ObjectValue(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":                      types.StringValue("test-id"),
			"name":                    types.StringValue("TestParam"),
			"label":                   types.StringValue("Test Parameter"),
			"help_text":               types.StringValue("Help text"),
			"default_value":           defaultValue,
			"default_sensitive_value": defaultSensitiveValue,
			"display_settings":        displaySettings,
		},
	)

	paramList, _ := types.ListValue(
		types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()},
		[]attr.Value{paramObj},
	)

	return &schemas.StepTemplateTypeResourceModel{Parameters: paramList}
}

func TestValidateStepTemplateParameters_SensitiveWithDefaultValue(t *testing.T) {
	data := stepTemplateParameterModel("Sensitive", types.StringValue("hunter2"), types.StringNull())

	diags := validateStepTemplateParameters(context.Background(), data)

	if diags.HasError() {
		t.Fatalf("expected a warning rather than an error for a literal sensitive default, got: %v", diags)
	}

	if diags.WarningsCount() == 0 {
		t.Fatal("expected a warning for sensitive parameter using default_value")
	}

	found := false
	for _, diag := range diags.Warnings() {
		detail := diag.Detail()
		if strings.Contains(detail, "Sensitive") &&
			strings.Contains(detail, "default_value") &&
			strings.Contains(detail, "default_sensitive_value") {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected warning detail to mention Sensitive, default_value, and default_sensitive_value")
	}
}

func TestValidateStepTemplateParameters_SensitiveWithBindingDefaultValue(t *testing.T) {
	bindings := []string{
		"#{DefaultPasswordVariable}",
		"  #{DefaultPasswordVariable}  ",
		"prefix-#{DefaultPasswordVariable}",
		"#{DefaultPasswordVariable}-suffix",
		"#{First}#{Second}",
		"#{DefaultPasswordVariable | ToUpper}",
	}

	for _, binding := range bindings {
		t.Run(binding, func(t *testing.T) {
			data := stepTemplateParameterModel("Sensitive", types.StringValue(binding), types.StringNull())

			diags := validateStepTemplateParameters(context.Background(), data)

			if diags.HasError() {
				t.Fatalf("expected no error for a variable binding default, got: %v", diags)
			}

			if diags.WarningsCount() > 0 {
				t.Errorf("expected no warning for a variable binding default, got: %v", diags.Warnings())
			}
		})
	}
}

func TestValidateStepTemplateParameters_SensitiveWithBothDefaultValues(t *testing.T) {
	data := stepTemplateParameterModel(
		"Sensitive",
		types.StringValue("#{DefaultPasswordVariable}"),
		types.StringValue("hunter2"),
	)

	diags := validateStepTemplateParameters(context.Background(), data)

	if !diags.HasError() {
		t.Fatal("expected error when both default_value and default_sensitive_value are set")
	}

	found := false
	for _, diag := range diags.Errors() {
		detail := diag.Detail()
		if strings.Contains(detail, "both") &&
			strings.Contains(detail, "default_value") &&
			strings.Contains(detail, "default_sensitive_value") {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected error detail to state that both default_value and default_sensitive_value are set")
	}
}

func TestValidateStepTemplateParameters_SensitiveWithSensitiveValue(t *testing.T) {
	data := stepTemplateParameterModel("Sensitive", types.StringNull(), types.StringValue("hunter2"))

	diags := validateStepTemplateParameters(context.Background(), data)

	if diags.HasError() {
		t.Fatalf("expected no error for sensitive parameter using default_sensitive_value, got: %v", diags)
	}

	if diags.WarningsCount() > 0 {
		t.Errorf("expected no warning for sensitive parameter using default_sensitive_value, got: %v", diags.Warnings())
	}
}

func TestValidateStepTemplateParameters_NonSensitiveWithSensitiveValue(t *testing.T) {
	displaySettings, _ := types.MapValue(types.StringType, map[string]attr.Value{
		"Octopus.ControlType": types.StringValue("SingleLineText"),
	})

	paramObj, _ := types.ObjectValue(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":                      types.StringValue("test-id"),
			"name":                    types.StringValue("TestParam"),
			"label":                   types.StringValue("Test Parameter"),
			"help_text":               types.StringValue("Help text"),
			"default_value":           types.StringNull(),
			"default_sensitive_value": types.StringValue("default_value"),
			"display_settings":        displaySettings,
		},
	)

	paramList, _ := types.ListValue(
		types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()},
		[]attr.Value{paramObj},
	)

	data := &schemas.StepTemplateTypeResourceModel{
		Parameters: paramList,
	}

	diags := validateStepTemplateParameters(context.Background(), data)

	if !diags.HasError() {
		t.Fatal("expected error for non-sensitive parameter using default_sensitive_value")
	}

	found := false
	for _, diag := range diags {
		if diag.Summary() == "Invalid step template parameter configuration" {
			detail := diag.Detail()
			if strings.Contains(detail, "Octopus.ControlType") &&
				strings.Contains(detail, "default_sensitive_value") &&
				strings.Contains(detail, "default_value") {
				found = true
				break
			}
		}
	}

	if !found {
		t.Error("expected diagnostic detail to mention Octopus.ControlType, default_sensitive_value, and default_value")
	}
}
