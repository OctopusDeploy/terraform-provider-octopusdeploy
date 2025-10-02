package schemas

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var _ EntitySchema = LifecycleSchema{}

var AllowDeprecatedAndNewRetentionBlocks = true

type LifecycleSchema struct {
	AllowDeprecatedRetention bool
}

//////////////////
// RESOURCE SCHEMA

func (l LifecycleSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		MarkdownDescription: "This resource manages lifecycles in Octopus Deploy." +
			"\n\nLifecycle retention is set using either the `retention_policy` and `retention_with_strategy` blocks." +
			"\n- When using an octopus version prior to `2025.3`" +
			"\n	- the `release_retention_policy` and `tentacle_retention_policy` blocks are used" +
			"\n- when using an octopus version `2025.3` or later" +
			"\n	- the `release_retention_with_strategy` and `tentacle_retention_with_strategy` blocks may be used",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"space_id":    util.ResourceString().Optional().Computed().Description("The space ID associated with this resource.").PlanModifiers(stringplanmodifier.UseStateForUnknown()).Build(),
			"name":        util.ResourceString().Required().Description("The name of this resource.").Build(),
			"description": util.ResourceString().Optional().Computed().Default("").Description("The description of this lifecycle.").Build(),
		},
		Blocks: getResourceSchemaBlocks(l.AllowDeprecatedRetention, true),
	}
}

func getResourceSchemaPhaseBlock(allowDeprecatedRetention bool) resourceSchema.ListNestedBlock {
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
			Blocks: getResourceSchemaBlocks(allowDeprecatedRetention, false),
		},
	}
}

func getResourceSchemaBlocks(allowDeprecatedRetention bool, includesPhaseBlock bool) map[string]resourceSchema.Block {
	blocks := map[string]resourceSchema.Block{
		"release_retention_with_strategy":  getResourceSchemaRetentionBlock(),
		"tentacle_retention_with_strategy": getResourceSchemaRetentionBlock(),
	}
	if includesPhaseBlock {
		blocks["phase"] = getResourceSchemaPhaseBlock(allowDeprecatedRetention)
	}
	if allowDeprecatedRetention {
		blocks["release_retention_policy"] = getDeprecatedResourceSchemaRetentionBlock(allowDeprecatedRetention)
		blocks["tentacle_retention_policy"] = getDeprecatedResourceSchemaRetentionBlock(allowDeprecatedRetention)
	}
	return blocks
}

func getResourceSchemaRetentionBlock() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		Description: "Defines the retention policy for releases or tentacles.\n	- When this block is not included, the space-wide \"Default\" retention policy is used. \n 	- This block may only be used on Octopus server 2025.3 or later.",
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"strategy": util.ResourceString().
					Required().
					Validators(stringvalidator.OneOf(core.RetentionStrategyDefault, core.RetentionStrategyCount, core.RetentionStrategyForever)).
					Description("How retention will be set. Valid strategies are `Default`, `Forever`, and `Count`. The default value is `Default`." +
						"\n  - `strategy = \"Default\"`, is used if the retention is set by the space-wide default lifecycle retention policy. " +
						"When `Default` is used, no other attributes can be set since the specific retention policy is no longer defined within this lifecycle." +
						"\n  - `strategy = \"Forever\"`, is used if items within this lifecycle should never be deleted." +
						"\n  - `strategy = \"Count\"`, is used if a specific number of days/releases should be kept.").
					Build(),
				"quantity_to_keep": util.ResourceInt64().
					Optional().Computed().
					Validators(int64validator.AtLeast(1)).
					Description("The number of days/releases to keep.").
					Build(),
				"unit": util.ResourceString().
					Optional().Computed().
					Validators(stringvalidator.OneOfCaseInsensitive(core.RetentionUnitDays, core.RetentionUnitItems)).
					Description("The unit of quantity to keep. Valid units are Days or Items.").
					Build(),
			},
			Validators: []validator.Object{
				resourceSchemaRetentionValidator{},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}
func getDeprecatedResourceSchemaRetentionBlock(allowDeprecatedRetention bool) resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		DeprecationMessage: "This block will deprecate when octopus 2025.2 is no longer supported. After upgrading to octopus 2025.3 or higher, please use the `release_retention_with_strategy` and `tentacle_retention_with_strategy` blocks instead.",
		Description:        "Defines the retention policy for releases or tentacles.",
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"quantity_to_keep": util.ResourceInt64().
					Optional().Computed().
					Default(int64default.StaticInt64(30)).
					Validators(int64validator.AtLeast(0)).
					Description("The number of days/releases to keep. The default value is 30. If 0 then all are kept.").
					Build(),
				"should_keep_forever": util.ResourceBool().
					Optional().Computed().
					Default(booldefault.StaticBool(false)).
					Description("Indicates if items should never be deleted. The default value is false.").
					Build(),
				"unit": util.ResourceString().
					Optional().Computed().
					Default(stringdefault.StaticString("Days")).
					Validators(stringvalidator.OneOfCaseInsensitive("Days", "Items")).
					Description("The unit of quantity to keep. Valid units are Days or Items. The default value is Days.").
					Build(),
			},
			Validators: []validator.Object{
				deprecatedRetentionValidator{allowDeprecatedRetention: allowDeprecatedRetention},
			},
		},
	}
}

type resourceSchemaRetentionValidator struct{}

func (v resourceSchemaRetentionValidator) Description(ctx context.Context) string {
	return "ensures only a count strategy has a quantity_to_keep and unit"
}
func (v resourceSchemaRetentionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (v resourceSchemaRetentionValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	var retentionStrategy struct {
		Strategy       types.String `tfsdk:"strategy"`
		QuantityToKeep types.Int64  `tfsdk:"quantity_to_keep"`
		Unit           types.String `tfsdk:"unit"`
	}

	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &retentionStrategy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var unitIsPresent = !retentionStrategy.Unit.IsNull()
	var quantityToKeepIsPresent = !retentionStrategy.QuantityToKeep.IsNull()

	if retentionStrategy.Strategy.ValueString() == core.RetentionStrategyCount {
		if !unitIsPresent {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("unit"),
				"unit",
				"unit must be set when strategy is set to Count.",
			)
		}
		if !quantityToKeepIsPresent {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("quantity_to_keep"),
				"quantity_to_keep",
				"quantity_to_keep must be set when strategy is set to Count.",
			)
		}
	}
	if retentionStrategy.Strategy.ValueString() == core.RetentionStrategyForever || retentionStrategy.Strategy.ValueString() == core.RetentionStrategyDefault {
		if unitIsPresent {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("unit"),
				"unit",
				"unit must not be set when strategy is Forever or Default.",
			)
		}
		if quantityToKeepIsPresent {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("quantity_to_keep"),
				"quantity_to_keep",
				"quantity_to_keep must not be set when strategy is Forever or Default.",
			)
		}
	}
}

type deprecatedRetentionValidator struct {
	allowDeprecatedRetention bool
}

func (v deprecatedRetentionValidator) Description(ctx context.Context) string {
	return "validates that should_keep_forever is true only if quantity_to_keep is 0"
}
func (v deprecatedRetentionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (v deprecatedRetentionValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	var retentionPolicy struct {
		QuantityToKeep    types.Int64  `tfsdk:"quantity_to_keep"`
		ShouldKeepForever types.Bool   `tfsdk:"should_keep_forever"`
		Unit              types.String `tfsdk:"unit"`
	}

	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &retentionPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !v.allowDeprecatedRetention {
		resp.Diagnostics.AddError(
			"release_retention_policy and tentacle_retention_policy are deprecated.",
			"Please use the `release_retention_with_strategy` and `tentacle_retention_with_strategy` blocks instead.")
	}

	if !retentionPolicy.QuantityToKeep.IsNull() && !retentionPolicy.QuantityToKeep.IsUnknown() && !retentionPolicy.ShouldKeepForever.IsNull() && !retentionPolicy.ShouldKeepForever.IsUnknown() {
		quantityToKeep := retentionPolicy.QuantityToKeep.ValueInt64()
		shouldKeepForever := retentionPolicy.ShouldKeepForever.ValueBool()

		if quantityToKeep == 0 && !shouldKeepForever {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("should_keep_forever"),
				"Invalid retention policy configuration",
				"should_keep_forever must be true when quantity_to_keep is 0",
			)
		} else if quantityToKeep != 0 && shouldKeepForever {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("should_keep_forever"),
				"Invalid retention policy configuration",
				"should_keep_forever must be false when quantity_to_keep is not 0",
			)
		}
	}

	if !retentionPolicy.Unit.IsNull() && !retentionPolicy.Unit.IsUnknown() {
		unit := retentionPolicy.Unit.ValueString()
		if !strings.EqualFold(unit, "Days") && !strings.EqualFold(unit, "Items") {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName("unit"),
				"Invalid retention policy unit",
				"Unit must be either 'Days' or 'Items' (case insensitive)",
			)
		}
	}
}

//////////////////
// DATASOURCE SCHEMA

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
			"lifecycles":   util.Ternary(l.AllowDeprecatedRetention, getDeprecatedDatasourceSchemaLifecycles(), getDatasourceSchemaLifecycles()),
		},
	}
}

func getDatasourceSchemaLifecycles() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		Optional: false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"id":                               util.DataSourceString().Computed().Description("The ID of the lifecycle.").Build(),
				"space_id":                         util.DataSourceString().Computed().Description("The space ID associated with this lifecycle.").Build(),
				"name":                             util.DataSourceString().Computed().Description("The name of the lifecycle.").Build(),
				"description":                      util.DataSourceString().Computed().Description("The description of the lifecycle.").Build(),
				"phase":                            getDatasourceSchemaPhases(),
				"release_retention_with_strategy":  getDatasourceSchemaRetention(),
				"tentacle_retention_with_strategy": getDatasourceSchemaRetention(),
			},
		},
	}
}
func getDeprecatedDatasourceSchemaLifecycles() datasourceSchema.ListNestedAttribute {
	var attributes = getDatasourceSchemaLifecycles().NestedObject.Attributes
	attributes["phase"] = getDeprecatedDatasourceSchemaPhases()
	attributes["release_retention_policy"] = getDeprecatedDatasourceSchemaRetention()
	attributes["tentacle_retention_policy"] = getDeprecatedDatasourceSchemaRetention()
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		Optional: false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: attributes,
		},
	}
}

func getDatasourceSchemaPhases() datasourceSchema.ListNestedAttribute {
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
				"release_retention_with_strategy":       getDatasourceSchemaRetention(),
				"tentacle_retention_with_strategy":      getDatasourceSchemaRetention(),
			},
		},
	}
}
func getDeprecatedDatasourceSchemaPhases() datasourceSchema.ListNestedAttribute {
	var attributes = getDatasourceSchemaPhases().NestedObject.Attributes
	attributes["release_retention_policy"] = getDeprecatedDatasourceSchemaRetention()
	attributes["tentacle_retention_policy"] = getDeprecatedDatasourceSchemaRetention()

	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: attributes,
		},
	}
}

func getDatasourceSchemaRetention() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"strategy":         util.DataSourceString().Computed().Description("The retention policy strategy. Can be \"Default\", \"Forever\", and \"Count\". \n  - \"Default\" indicates retention is set by the Space Default retention policy for lifecycles \n  - \"Forever\" indicates releases are never deleted \n  - \"Count\" indicates releases are kept according to `unit` and `quantity_to_keep`").Build(),
				"quantity_to_keep": util.DataSourceInt64().Computed().Description("The unit for `quantity_to_keep`. Dismiss when `strategy` is \"Forever\" or \"Default\".").Build(),
				"unit":             util.DataSourceString().Computed().Description("The number of units to keep. Dismiss when `strategy` is \"Forever\" or \"Default\".").Build(),
			},
		},
	}
}
func getDeprecatedDatasourceSchemaRetention() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"strategy":         util.DataSourceString().Computed().Description("The retention policy strategy. Can be \"Default\", \"Forever\", and \"Count\". \n  - \"Default\" indicates retention is set by the Space Default retention policy for lifecycles \n  - \"Forever\" indicates releases are never deleted \n  - \"Count\" indicates releases are kept according to `unit` and `quantity_to_keep`").Build(),
				"quantity_to_keep": util.DataSourceInt64().Computed().Description("The unit for `quantity_to_keep`. Dismiss when `strategy` is \"Forever\" or \"Default\".").Build(),
				"unit":             util.DataSourceString().Computed().Description("The number of units to keep. Dismiss when `strategy` is \"Forever\" or \"Default\".").Build(),
			},
		},
	}
}
