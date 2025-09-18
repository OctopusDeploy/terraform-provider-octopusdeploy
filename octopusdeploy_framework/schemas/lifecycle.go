package schemas

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"strings"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ EntitySchema = LifecycleSchema{}

type LifecycleSchema struct{}

func (l LifecycleSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages lifecycles in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"space_id":    util.ResourceString().Optional().Computed().Description("The space ID associated with this resource.").PlanModifiers(stringplanmodifier.UseStateForUnknown()).Build(),
			"name":        util.ResourceString().Required().Description("The name of this resource.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this lifecycle.").Build(),
		},
		Blocks: map[string]resourceSchema.Block{
			"phase":                     getResourcePhaseBlockSchema(),
			"release_retention_policy":  getResourceRetentionPolicyBlockSchema(false),
			"tentacle_retention_policy": getResourceRetentionPolicyBlockSchema(false),
		},
	}
}

func (l LifecycleSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing lifecycles.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":           util.DataSourceString().Computed().Description("The ID of the lifecycle.").Build(),
			"space_id":     util.DataSourceString().Optional().Description("The space ID associated with this lifecycle.").Build(),
			"ids":          util.DataSourceList(types.StringType).Optional().Description("A list of lifecycle IDs to filter by.").Build(),
			"partial_name": util.DataSourceString().Optional().Description("A partial name to filter lifecycles by.").Build(),
			"skip":         util.DataSourceInt64().Optional().Description("A filter to specify the number of items to skip in the response.").Build(),
			"take":         util.DataSourceInt64().Optional().Description("A filter to specify the number of items to take (or return) in the response.").Build(),
			"lifecycles":   getLifecyclesAttribute(),
		},
	}
}

func getResourcePhaseBlockSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		Description: "Defines a phase in the lifecycle.",
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"id": util.ResourceString().
					Optional().Computed().
					Description("The unique ID for this resource.").
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"name": util.ResourceString().Required().Description("The name of this resource.").Build(),
				"automatic_deployment_targets": util.ResourceList(types.StringType).
					Optional().Computed().
					Description("Environment IDs in this phase that a release is automatically deployed to when it is eligible for this phase").
					PlanModifiers(listplanmodifier.UseStateForUnknown()).
					Build(),
				"optional_deployment_targets": util.ResourceList(types.StringType).
					Optional().Computed().
					Description("Environment IDs in this phase that a release can be deployed to, but is not automatically deployed to").
					PlanModifiers(listplanmodifier.UseStateForUnknown()).
					Build(),
				"minimum_environments_before_promotion": util.ResourceInt64().
					Optional().Computed().
					Default(int64default.StaticInt64(0)).
					Description("The number of units required before a release can enter the next phase. If 0, all environments are required.").
					PlanModifiers(int64planmodifier.UseStateForUnknown()).
					Build(),
				"is_optional_phase": util.ResourceBool().
					Optional().Computed().
					Description("If false a release must be deployed to this phase before it can be deployed to the next phase.").
					Default(booldefault.StaticBool(false)).
					PlanModifiers(boolplanmodifier.UseStateForUnknown()).
					Build(),
				"is_priority_phase": util.ResourceBool().
					Optional().Computed().
					Default(booldefault.StaticBool(false)).
					Description("Deployments will be prioritized in this phase").
					PlanModifiers(boolplanmodifier.UseStateForUnknown()).
					Build(),
			},
			Blocks: map[string]resourceSchema.Block{
				"release_retention_policy":  getResourceRetentionPolicyBlockSchema(true),
				"tentacle_retention_policy": getResourceRetentionPolicyBlockSchema(true),
			},
		},
	}
}

func getResourceRetentionPolicyBlockDescription(isPhase bool) string {
	if isPhase {
		return "Defines the retention policy for releases or tentacles within the phase. \n If this block is not set, the retention policy will be inherited from the lifecycle."
	}
	return "Defines the retention policy for releases or tentacles. \n If this block is not set, the strategy will be` strategy = \"Count\"` with `30` `Days` of saved releases."
}

func getResourceRetentionPolicyBlockSchema(isPhase bool) resourceSchema.ListNestedBlock {
	var strategyDescription = "How retention will be set. Valid strategies are `Default`, `Forever`, and `Count`. The default value is `Default`.\n  - `strategy = \"Default\"`, is used if the retention is set by the space-wide default lifecycle retention policy. When `Default` is used, no other attributes can be set since the specific retention policy is no longer defined within this lifecycle.\n  - `strategy = \"Forever\"`, is used if items within this lifecycle should never be deleted.\n  - `strategy = \"Count\"`, is used if a specific number of days/releases should be kept."

	return resourceSchema.ListNestedBlock{
		Description: getResourceRetentionPolicyBlockDescription(isPhase),
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"strategy": util.ResourceString().
					Optional().Computed().
					Validators(stringvalidator.OneOf(core.RetentionStrategyDefault, core.RetentionStrategyCount, core.RetentionStrategyForever)).
					Description(strategyDescription).
					Build(),
				"quantity_to_keep": util.ResourceInt64().
					Optional().Computed().
					Validators(int64validator.AtLeast(0)).
					Description("The number of days/releases to keep. If 0 then all are kept.").
					Build(),
				"should_keep_forever": util.ResourceBool().
					Optional().Computed().
					Description("Indicates if items should never be deleted. For best practice, use strategy=\"forever\" instead").
					Build(),
				"unit": util.ResourceString().
					Optional().Computed().
					Description("The unit of quantity to keep. Valid units are Days or Items.").
					Validators(stringvalidator.OneOfCaseInsensitive("Days", "Items")).
					Build(),
			},
			Validators: []validator.Object{
				retentionPolicyValidator{},
			},
		},
	}
}

func getLifecyclesAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		Optional: false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"id":                        util.DataSourceString().Computed().Description("The ID of the lifecycle.").Build(),
				"space_id":                  util.DataSourceString().Computed().Description("The space ID associated with this lifecycle.").Build(),
				"name":                      util.DataSourceString().Computed().Description("The name of the lifecycle.").Build(),
				"description":               util.DataSourceString().Computed().Description("The description of the lifecycle.").Build(),
				"phase":                     getPhasesAttribute(),
				"release_retention_policy":  getRetentionPolicyAttribute(),
				"tentacle_retention_policy": getRetentionPolicyAttribute(),
			},
		},
	}
}

func getPhasesAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"id":                                    util.DataSourceString().Computed().Description("The ID of the phase.").Build(),
				"name":                                  util.DataSourceString().Computed().Description("The name of the phase.").Build(),
				"automatic_deployment_targets":          util.DataSourceList(types.StringType).Computed().Description("The automatic deployment targets for this phase.").Build(),
				"optional_deployment_targets":           util.DataSourceList(types.StringType).Computed().Description("The optional deployment targets for this phase.").Build(),
				"minimum_environments_before_promotion": util.DataSourceInt64().Computed().Description("The minimum number of environments before promotion.").Build(),
				"is_optional_phase":                     util.DataSourceBool().Computed().Description("Whether this phase is optional.").Build(),
				"is_priority_phase":                     util.DataSourceBool().Computed().Description("Deployments will be prioritized in this phase").Build(),
				"release_retention_policy":              getRetentionPolicyAttribute(),
				"tentacle_retention_policy":             getRetentionPolicyAttribute(),
			},
		},
	}
}

func getRetentionPolicyAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"strategy":            util.DataSourceString().Computed().Description("The retention policy strategy. Can be \"Default\", \"Forever\", and \"Count\". \n  - \"Default\" indicates retention is set by the Space Default retention policy for lifecycles \n  - \"Forever\" indicates releases are never deleted \n  - \"Count\" indicates releases are kept according to `unit` and `quantity_to_keep`").Build(),
				"quantity_to_keep":    util.DataSourceInt64().Computed().Description("The number of units to keep. Dismiss when `strategy` is \"Forever\" or \"Default\".").Build(),
				"should_keep_forever": util.DataSourceBool().Computed().Description("Whether releases should be kept forever. Dismiss when `strategy` is \"Forever\".").Build(),
				"unit":                util.DataSourceString().Computed().Description("The unit for `quantity_to_keep`. Dismiss when `strategy` is \"Forever\" or \"Default\".").Build(),
			},
		},
	}
}

type retentionPolicyValidator struct{}

func (v retentionPolicyValidator) Description(ctx context.Context) string {
	return "validates that should_keep_forever is true only if quantity_to_keep is 0"
}

func (v retentionPolicyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v retentionPolicyValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	var retentionPolicy struct {
		Strategy          types.String `tfsdk:"strategy"`
		QuantityToKeep    types.Int64  `tfsdk:"quantity_to_keep"`
		ShouldKeepForever types.Bool   `tfsdk:"should_keep_forever"`
		Unit              types.String `tfsdk:"unit"`
	}

	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &retentionPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var strategy = retentionPolicy.Strategy
	var quantityToKeep = retentionPolicy.QuantityToKeep
	var shouldKeepForever = retentionPolicy.ShouldKeepForever
	var unit = retentionPolicy.Unit

	if strategy.IsNull() || strategy.IsUnknown() {
		v.ValidateRetentionObjectWithoutStrategy(req, resp, quantityToKeep, shouldKeepForever, unit)
	} else {
		v.ValidateRetentionObjectWithStrategy(req, resp, strategy, quantityToKeep, shouldKeepForever, unit)
	}

}

func (v retentionPolicyValidator) ValidateRetentionObjectWithoutStrategy(req validator.ObjectRequest, resp *validator.ObjectResponse, quantityToKeep types.Int64, shouldKeepForever types.Bool, unit types.String) {
	unitPresent := !unit.IsNull() && !unit.IsUnknown()
	quantityToKeepPresent := !quantityToKeep.IsNull() && !quantityToKeep.IsUnknown()
	shouldKeepForeverPresent := !shouldKeepForever.IsNull() && !shouldKeepForever.IsUnknown()
	shouldKeepForeverIsTrue := shouldKeepForeverPresent && shouldKeepForever.ValueBool() == true
	shouldKeepForeverIsFalse := shouldKeepForeverPresent && shouldKeepForever.ValueBool() == false
	quantityToKeepIsMoreThanZero := quantityToKeepPresent && quantityToKeep.ValueInt64() > 0
	errorGuidingToStrategy := "Incorrect use of retention attributes.\nFor best practice use strategy attribute."

	if !unitPresent && !quantityToKeepPresent && !shouldKeepForeverPresent {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("strategy"),
			"Invalid retention policy configuration",
			"please either add retention policy attributes or remove the entire block",
		)
	}

	// count strategy validations
	if quantityToKeepIsMoreThanZero {
		if shouldKeepForeverIsTrue {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("should_keep_forever"),
				"Invalid retention policy configuration",
				errorGuidingToStrategy,
				//should_keep_forever must be false when quantity_to_keep is greater than 0
			)
		}
		if !unitPresent {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("unit"),
				"Invalid retention policy configuration",
				errorGuidingToStrategy,
				//unit is required when quantity_to_keep is greater than 0
			)
		}
	}

	// keep forever strategy validation
	if !quantityToKeepIsMoreThanZero && shouldKeepForeverIsFalse {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("should_keep_forever"),
			"Invalid retention policy configuration",
			errorGuidingToStrategy,
			//should_keep_forever must be true when quantity_to_keep is zero or missing
		)
	}

	if unitPresent && !quantityToKeepIsMoreThanZero {
		if strings.EqualFold(unit.ValueString(), "Items") {
			// do not throw an error for backwards compatability.
		} else {
			resp.Diagnostics.AddAttributeError(
				// replaces a confusing state change to "unit = Items" error at the api
				req.Path.AtName("unit"),
				"Invalid retention policy configuration",
				errorGuidingToStrategy,
				//"unit is only used when quantity_to_keep is greater than 0"
			)
		}
	}
}

func (v retentionPolicyValidator) ValidateRetentionObjectWithStrategy(req validator.ObjectRequest, resp *validator.ObjectResponse, strategy types.String, quantityToKeep types.Int64, shouldKeepForever types.Bool, unit types.String) {
	unitPresent := !unit.IsNull() && !unit.IsUnknown()
	quantityToKeepPresent := !quantityToKeep.IsNull() && !quantityToKeep.IsUnknown()
	shouldKeepForeverPresent := !shouldKeepForever.IsNull() && !shouldKeepForever.IsUnknown()

	if strategy.ValueString() == core.RetentionStrategyDefault {
		if quantityToKeepPresent {
			rejectAttribute("quantity_to_keep", core.RetentionStrategyDefault, resp, req)
		}
		if unitPresent {
			rejectAttribute("unit", core.RetentionStrategyDefault, resp, req)
		}
		if shouldKeepForeverPresent {
			rejectAttribute("should_keep_forever", core.RetentionStrategyDefault, resp, req)
		}
	}
	if strategy.ValueString() == core.RetentionStrategyForever {
		if quantityToKeepPresent {
			rejectAttribute("quantity_to_keep", core.RetentionStrategyForever, resp, req)
		}
		if unitPresent {
			rejectAttribute("unit", core.RetentionStrategyForever, resp, req)
		}
		if shouldKeepForeverPresent {
			rejectAttribute("should_keep_forever", core.RetentionStrategyForever, resp, req)
		}
	}
	if strategy.ValueString() == core.RetentionStrategyCount {
		if shouldKeepForeverPresent {
			rejectAttribute("should_keep_forever", core.RetentionStrategyCount, resp, req)
		}
		if !strings.EqualFold(unit.ValueString(), "Days") && !strings.EqualFold(unit.ValueString(), "Items") {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("unit"),
				"Invalid retention policy unit",
				"Unit must be either 'Days' or 'Items' (case insensitive) when strategy is set to Count",
			)
		}
		if !quantityToKeepPresent || quantityToKeep.ValueInt64() == 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("quantity_to_keep"),
				"Invalid retention policy configuration",
				"quantity_to_keep must be greater than 0 when strategy is set to Count",
			)
		}
	}

}

func rejectAttribute(attributeToReject string, strategy string, resp *validator.ObjectResponse, req validator.ObjectRequest) {
	resp.Diagnostics.AddAttributeError(
		req.Path.AtName(attributeToReject),
		"Invalid retention policy configuration",
		fmt.Sprintf("%s should not be supplied when strategy is set to %s", attributeToReject, strategy),
	)
}
