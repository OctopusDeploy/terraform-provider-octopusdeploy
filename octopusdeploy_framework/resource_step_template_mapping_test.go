package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccMapStepTemplateParametersFromStateSingleLineText(t *testing.T) {
	ctx := context.Background()

	defaultValue := "defaultValue"
	parameterDefaultValue := core.NewPropertyValue(defaultValue, false)
	parameter := actiontemplates.ActionTemplateParameter{
		Name:     "Param.One",
		Label:    "Parameter One",
		HelpText: "First parameter",
		DisplaySettings: map[string]string{
			"Octopus.ControlType": "SingleLineText",
		},
		DefaultValue: &parameterDefaultValue,
	}
	parameter.SetID("00000000-0000-0000-0000-000000000001")

	parameterConfig := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue(parameter.ID),
			"name":      types.StringValue(parameter.Name),
			"label":     types.StringValue(parameter.Label),
			"help_text": types.StringValue(parameter.HelpText),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("SingleLineText"),
			}),
			"default_value":           types.StringValue(defaultValue),
			"default_sensitive_value": types.StringNull(),
		},
	)

	state := schemas.StepTemplateTypeResourceModel{
		SpaceID:         types.StringValue("Spaces-1"),
		Name:            types.StringValue("Basic Template"),
		ActionType:      types.StringValue("Octopus.Script"),
		StepPackageId:   types.StringValue("Octopus.Script"),
		Packages:        types.ListValueMust(schemas.StepTemplatePackageObjectType(), []attr.Value{}),
		GitDependencies: types.ListValueMust(schemas.StepTemplateGitDependencyObjectType(), []attr.Value{}),
		Parameters: types.ListValueMust(schemas.StepTemplateParameterObjectType(), []attr.Value{
			parameterConfig,
		}),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host 'Test'"),
		}),
	}

	// Act
	template, diags := mapStepTemplateResourceModelToActionTemplate(ctx, state)
	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedTemplate := &actiontemplates.ActionTemplate{
		SpaceID:         "Spaces-1",
		ActionType:      "Octopus.Script",
		Name:            "Basic Template",
		Packages:        []packages.PackageReference{},
		GitDependencies: []gitdependencies.GitDependency{},
		Parameters: []actiontemplates.ActionTemplateParameter{
			parameter,
		},
		Properties: map[string]core.PropertyValue{
			"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host 'Test'", false),
		},
	}
	expectedTemplate.Links = map[string]string{}

	assert.Equal(t, expectedTemplate, template)
}

func TestAccMapStepTemplateParametersToStateSingleLineText(t *testing.T) {
	ctx := context.Background()

	defaultValue := "random-value"
	parameterDefaultValue := core.NewPropertyValue(defaultValue, false)
	parameter := actiontemplates.ActionTemplateParameter{
		Name:     "Param.One",
		Label:    "Parameter One",
		HelpText: "First parameter",
		DisplaySettings: map[string]string{
			"Octopus.ControlType": "SingleLineText",
		},
		DefaultValue: &parameterDefaultValue,
	}
	parameter.SetID("00000000-0000-0000-0000-000000000001")

	template := &actiontemplates.ActionTemplate{
		SpaceID:         "Spaces-1",
		ActionType:      "Octopus.Script",
		Name:            "Basic Template",
		Packages:        []packages.PackageReference{},
		GitDependencies: []gitdependencies.GitDependency{},
		Parameters: []actiontemplates.ActionTemplateParameter{
			parameter,
		},
		Properties: map[string]core.PropertyValue{
			"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host 'Test'", false),
		},
	}
	template.SetID("StepTemplates-22")

	// Act
	state := schemas.StepTemplateTypeResourceModel{}
	diags := mapStepTemplateToResourceModel(ctx, &state, template)
	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	parameterConfig := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue(parameter.ID),
			"name":      types.StringValue(parameter.Name),
			"label":     types.StringValue(parameter.Label),
			"help_text": types.StringValue(parameter.HelpText),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("SingleLineText"),
			}),
			"default_value":           types.StringValue(defaultValue),
			"default_sensitive_value": types.StringNull(),
		},
	)

	expectedState := schemas.StepTemplateTypeResourceModel{
		SpaceID:                   types.StringValue("Spaces-1"),
		Name:                      types.StringValue("Basic Template"),
		Description:               types.StringValue(""),
		ActionType:                types.StringValue("Octopus.Script"),
		StepPackageId:             types.StringNull(), // not set during mapping
		CommunityActionTemplateId: types.StringValue(""),
		Version:                   types.Int32Value(0),
		Packages:                  types.ListValueMust(schemas.StepTemplatePackageObjectType(), []attr.Value{}),
		GitDependencies:           types.ListValueMust(schemas.StepTemplateGitDependencyObjectType(), []attr.Value{}),
		Parameters: types.ListValueMust(schemas.StepTemplateParameterObjectType(), []attr.Value{
			parameterConfig,
		}),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host 'Test'"),
		}),
	}
	expectedState.ID = types.StringValue(template.ID)

	assert.Equal(t, expectedState, state)
}

func TestAccMapStepTemplateParametersFromStateSensitive(t *testing.T) {
	ctx := context.Background()

	sensitiveValue := "secret-value"
	parameterDefaultValue := core.NewPropertyValue(sensitiveValue, true)
	parameter := actiontemplates.ActionTemplateParameter{
		Name:     "Param.Secret",
		Label:    "Parameter Secret",
		HelpText: "Second parameter",
		DisplaySettings: map[string]string{
			"Octopus.ControlType": "Sensitive",
		},
		DefaultValue: &parameterDefaultValue,
	}
	parameter.SetID("00000000-0000-0000-0000-000000000003")

	parameterConfig := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue(parameter.ID),
			"name":      types.StringValue(parameter.Name),
			"label":     types.StringValue(parameter.Label),
			"help_text": types.StringValue(parameter.HelpText),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("Sensitive"),
			}),
			"default_value":           types.StringValue(""),
			"default_sensitive_value": types.StringValue(sensitiveValue),
		},
	)

	state := schemas.StepTemplateTypeResourceModel{
		SpaceID:         types.StringValue("Spaces-1"),
		Name:            types.StringValue("Basic Template"),
		ActionType:      types.StringValue("Octopus.Script"),
		StepPackageId:   types.StringValue("Octopus.Script"),
		Packages:        types.ListValueMust(schemas.StepTemplatePackageObjectType(), []attr.Value{}),
		GitDependencies: types.ListValueMust(schemas.StepTemplateGitDependencyObjectType(), []attr.Value{}),
		Parameters: types.ListValueMust(schemas.StepTemplateParameterObjectType(), []attr.Value{
			parameterConfig,
		}),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host 'Test'"),
		}),
	}

	// Act
	template, diags := mapStepTemplateResourceModelToActionTemplate(ctx, state)
	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedTemplate := &actiontemplates.ActionTemplate{
		SpaceID:         "Spaces-1",
		ActionType:      "Octopus.Script",
		Name:            "Basic Template",
		Packages:        []packages.PackageReference{},
		GitDependencies: []gitdependencies.GitDependency{},
		Parameters: []actiontemplates.ActionTemplateParameter{
			parameter,
		},
		Properties: map[string]core.PropertyValue{
			"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host 'Test'", false),
		},
	}
	expectedTemplate.Links = map[string]string{}

	assert.Equal(t, expectedTemplate, template)
}

func TestAccMapStepTemplateParametersToStateSensitive(t *testing.T) {
	ctx := context.Background()

	parameterDefaultValue := core.NewPropertyValue("", true) // Sensitive value is empty when returned from the Server
	parameter := actiontemplates.ActionTemplateParameter{
		Name:     "Param.Secret",
		Label:    "Parameter Secret",
		HelpText: "Second parameter",
		DisplaySettings: map[string]string{
			"Octopus.ControlType": "Sensitive",
		},
		DefaultValue: &parameterDefaultValue,
	}
	parameter.SetID("00000000-0000-0000-0000-000000000003")

	template := &actiontemplates.ActionTemplate{
		SpaceID:         "Spaces-1",
		ActionType:      "Octopus.Script",
		Name:            "Basic Template",
		Packages:        []packages.PackageReference{},
		GitDependencies: []gitdependencies.GitDependency{},
		Parameters: []actiontemplates.ActionTemplateParameter{
			parameter,
		},
		Properties: map[string]core.PropertyValue{
			"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host 'Test'", false),
		},
	}
	template.SetID("StepTemplates-22")

	// For sensitive default value to avoid conflicts between planned and applied states,
	// we need to copy configured sensitive value to the state after application
	configuredParameter := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue(parameter.ID),
			"name":      types.StringValue(parameter.Name),
			"label":     types.StringValue(parameter.Label),
			"help_text": types.StringValue(parameter.HelpText),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("Sensitive"),
			}),
			"default_value":           types.StringValue(""),
			"default_sensitive_value": types.StringValue("configured-secret"),
		},
	)

	// Act
	state := schemas.StepTemplateTypeResourceModel{
		Parameters: types.ListValueMust(schemas.StepTemplateParameterObjectType(), []attr.Value{
			configuredParameter,
		}),
	}

	diags := mapStepTemplateToResourceModel(ctx, &state, template)
	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	updatedParameter := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue(parameter.ID),
			"name":      types.StringValue(parameter.Name),
			"label":     types.StringValue(parameter.Label),
			"help_text": types.StringValue(parameter.HelpText),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("Sensitive"),
			}),
			"default_value":           types.StringValue(""), // Sensitive values are always empty when writing to state
			"default_sensitive_value": types.StringValue("configured-secret"),
		},
	)

	expectedState := schemas.StepTemplateTypeResourceModel{
		SpaceID:                   types.StringValue("Spaces-1"),
		Name:                      types.StringValue("Basic Template"),
		Description:               types.StringValue(""),
		ActionType:                types.StringValue("Octopus.Script"),
		StepPackageId:             types.StringNull(), // not set during mapping
		CommunityActionTemplateId: types.StringValue(""),
		Version:                   types.Int32Value(0),
		Packages:                  types.ListValueMust(schemas.StepTemplatePackageObjectType(), []attr.Value{}),
		GitDependencies:           types.ListValueMust(schemas.StepTemplateGitDependencyObjectType(), []attr.Value{}),
		Parameters: types.ListValueMust(schemas.StepTemplateParameterObjectType(), []attr.Value{
			updatedParameter,
		}),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host 'Test'"),
		}),
	}
	expectedState.ID = types.StringValue(template.ID)

	assert.Equal(t, expectedState, state)
}

func TestAccMapStepTemplateParametersValidationWhenNonSensitiveDefaultValueSetForSensitiveControlType(t *testing.T) {
	ctx := context.Background()

	sensitiveParameter := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue("00000000-0000-0000-0000-000000000001"),
			"name":      types.StringValue("Parameter Name"),
			"label":     types.StringValue("Parameter Label"),
			"help_text": types.StringNull(),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("Sensitive"),
			}),
			"default_value":           types.StringValue("secret-value"),
			"default_sensitive_value": types.StringNull(),
		},
	)

	textParameter := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue("00000000-0000-0000-0000-000000000001"),
			"name":      types.StringValue("Parameter Name"),
			"label":     types.StringValue("Parameter Label"),
			"help_text": types.StringNull(),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("SingleLineText"),
			}),
			"default_value":           types.StringValue("plain-text"),
			"default_sensitive_value": types.StringNull(),
		},
	)

	state := schemas.StepTemplateTypeResourceModel{
		SpaceID:         types.StringValue("Spaces-1"),
		Name:            types.StringValue("Basic Template"),
		ActionType:      types.StringValue("Octopus.Script"),
		StepPackageId:   types.StringValue("Octopus.Script"),
		Packages:        types.ListValueMust(schemas.StepTemplatePackageObjectType(), []attr.Value{}),
		GitDependencies: types.ListValueMust(schemas.StepTemplateGitDependencyObjectType(), []attr.Value{}),
		Parameters: types.ListValueMust(schemas.StepTemplateParameterObjectType(), []attr.Value{
			sensitiveParameter,
			textParameter,
		}),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host 'Test'"),
		}),
	}

	// Act
	diags := validateStepTemplateParameters(ctx, &state)

	diagnostics := make([]diag.Severity, len(diags))
	for i, d := range diags {
		diagnostics[i] = d.Severity()
	}
	expectedDiagnostics := []diag.Severity{diag.SeverityError}
	assert.Equal(t, expectedDiagnostics, diagnostics, "Expected diagnostics to contain errors")
}

func TestAccMapStepTemplateParametersValidationWhenSensitiveDefaultValueSetForNonSensitiveControlType(t *testing.T) {
	ctx := context.Background()

	sensitiveParameter := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue("00000000-0000-0000-0000-000000000001"),
			"name":      types.StringValue("Parameter Name"),
			"label":     types.StringValue("Parameter Label"),
			"help_text": types.StringNull(),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("Sensitive"),
			}),
			"default_value":           types.StringValue(""),
			"default_sensitive_value": types.StringNull(),
		},
	)

	textParameter := types.ObjectValueMust(
		schemas.GetStepTemplateParameterTypeAttributes(),
		map[string]attr.Value{
			"id":        types.StringValue("00000000-0000-0000-0000-000000000001"),
			"name":      types.StringValue("Parameter Name"),
			"label":     types.StringValue("Parameter Label"),
			"help_text": types.StringNull(),
			"display_settings": types.MapValueMust(types.StringType, map[string]attr.Value{
				"Octopus.ControlType": types.StringValue("SingleLineText"),
			}),
			"default_value":           types.StringValue(""),
			"default_sensitive_value": types.StringValue("plain-text-in-sensitive-attribute"),
		},
	)

	state := schemas.StepTemplateTypeResourceModel{
		SpaceID:         types.StringValue("Spaces-1"),
		Name:            types.StringValue("Basic Template"),
		ActionType:      types.StringValue("Octopus.Script"),
		StepPackageId:   types.StringValue("Octopus.Script"),
		Packages:        types.ListValueMust(schemas.StepTemplatePackageObjectType(), []attr.Value{}),
		GitDependencies: types.ListValueMust(schemas.StepTemplateGitDependencyObjectType(), []attr.Value{}),
		Parameters: types.ListValueMust(schemas.StepTemplateParameterObjectType(), []attr.Value{
			sensitiveParameter,
			textParameter,
		}),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host 'Test'"),
		}),
	}

	// Act
	diags := validateStepTemplateParameters(ctx, &state)

	diagnostics := make([]diag.Severity, len(diags))
	for i, d := range diags {
		diagnostics[i] = d.Severity()
	}
	expectedDiagnostics := []diag.Severity{diag.SeverityError}
	assert.Equal(t, expectedDiagnostics, diagnostics, "Expected diagnostics to contain errors")
}
