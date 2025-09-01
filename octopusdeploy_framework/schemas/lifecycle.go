package schemas

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
			"release_retention_policy":  getResourceRetentionPolicyBlockSchema(),
			"tentacle_retention_policy": getResourceRetentionPolicyBlockSchema(),
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
				"release_retention_policy":  getResourceRetentionPolicyBlockSchema(),
				"tentacle_retention_policy": getResourceRetentionPolicyBlockSchema(),
			},
		},
	}
}

var strategyDescription = "How retention will be set. Valid strategies are `Default`, `Forever`, and `Count`. The default value is `Default`.\n  - `strategy = \"Default\"`, is used if the retention is set by the space-wide default lifecycle retention policy. When `Default` is used, no other attributes can be set since the specific retention policy is no longer defined within this lifecycle.\n  - `strategy = \"Forever\"`, is used if items within this lifecycle should never be deleted.\n  - `strategy = \"Count\"`, is used if a specific number of days/releases should be kept."

func getResourceRetentionPolicyBlockSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		Description: "Defines the retention policy for releases or tentacles.",
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"quantity_to_keep": util.ResourceInt64().
					Optional().Computed().
					Default(int64default.StaticInt64(30)).
					Validators(int64validator.AtLeast(0)).
					Description("The number of days/releases to keep. This number should be larger than 0.").
					Build(),
				"should_keep_forever": util.ResourceBool().
					Deprecated("Use strategy instead.").
					Optional().Computed().
					Default(booldefault.StaticBool(false)).
					Description("A depreciated attribute indicating if items should never be deleted. Octopus recommends using `strategy = \"Forever\"` instead.").
					Build(),
				"strategy": util.ResourceString().
					Optional().Computed().
					Validators(stringvalidator.OneOf(core.RetentionStrategyDefault, core.RetentionStrategyCount, core.RetentionStrategyForever)).
					Description(strategyDescription).
					Build(),
				"unit": util.ResourceString().
					Optional().Computed().
					Default(stringdefault.StaticString("Days")).
					Description("The unit of quantity to keep. Valid units are Days or Items.").
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
				"strategy":            util.DataSourceString().Computed().Description("The retention policy strategy.").Build(),
				"quantity_to_keep":    util.DataSourceInt64().Computed().Description("The quantity of releases to keep.").Build(),
				"should_keep_forever": util.DataSourceBool().Computed().Description("Whether releases should be kept forever.").Build(),
				"unit":                util.DataSourceString().Computed().Description("The unit of time for the retention policy.").Build(),
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

	rejectAttribute := func(attributeToReject string, strategy string) {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName(attributeToReject),
			"Invalid retention policy configuration",
			fmt.Sprintf("%s should not be supplied when strategy is set to %s", attributeToReject, strategy),
		)
	}
	// If a strategy is present, it overrides all other retention behaviour. Other unrelated attributes are rejected.
	if retentionPolicy.Strategy.ValueString() == core.RetentionStrategyDefault {
		if !retentionPolicy.QuantityToKeep.IsNull() {
			rejectAttribute("quantity_to_keep", core.RetentionStrategyDefault)
		}
		if !retentionPolicy.Unit.IsNull() {
			rejectAttribute("unit", core.RetentionStrategyDefault)
		}
		if !retentionPolicy.ShouldKeepForever.IsNull() {
			rejectAttribute("should_keep_forever", core.RetentionStrategyDefault)
		}
	}

	if retentionPolicy.Strategy.ValueString() == core.RetentionStrategyForever {
		if !retentionPolicy.QuantityToKeep.IsNull() {
			rejectAttribute("quantity_to_keep", core.RetentionStrategyForever)
		}
		if !retentionPolicy.Unit.IsNull() {
			rejectAttribute("unit", core.RetentionStrategyForever)
		}
		if !retentionPolicy.ShouldKeepForever.IsNull() {
			rejectAttribute("should_keep_forever", core.RetentionStrategyForever)
		}
	}

	if retentionPolicy.Strategy.ValueString() == core.RetentionStrategyCount {
		if !retentionPolicy.ShouldKeepForever.IsNull() {
			rejectAttribute("should_keep_forever", core.RetentionStrategyCount)
		}
	}

	if retentionPolicy.Strategy.IsNull() || retentionPolicy.Strategy.IsUnknown() {
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
