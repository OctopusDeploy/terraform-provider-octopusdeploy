package schemas

import (
	"context"
	"math"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SpaceDefaultLifecycleName is the name Octopus gives the lifecycle it creates alongside a new space.
const SpaceDefaultLifecycleName = "Default Lifecycle"

type SpaceDefaultLifecycleSchema struct{}

var _ EntitySchema = SpaceDefaultLifecycleSchema{}

type SpaceDefaultLifecycleResourceModel struct {
	SpaceID           types.String `tfsdk:"space_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Phases            types.List   `tfsdk:"phases"`
	ReleaseRetention  types.Object `tfsdk:"release_retention_with_strategy"`
	TentacleRetention types.Object `tfsdk:"tentacle_retention_with_strategy"`

	ResourceModel
}

type SpaceDefaultLifecyclePhaseModel struct {
	ID                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	AutomaticDeploymentTargets         types.List   `tfsdk:"automatic_deployment_targets"`
	OptionalDeploymentTargets          types.List   `tfsdk:"optional_deployment_targets"`
	MinimumEnvironmentsBeforePromotion types.Int64  `tfsdk:"minimum_environments_before_promotion"`
	IsOptionalPhase                    types.Bool   `tfsdk:"is_optional_phase"`
	IsPriorityPhase                    types.Bool   `tfsdk:"is_priority_phase"`
	ReleaseRetention                   types.Object `tfsdk:"release_retention_with_strategy"`
	TentacleRetention                  types.Object `tfsdk:"tentacle_retention_with_strategy"`
}

type SpaceDefaultLifecycleRetentionModel struct {
	Strategy       types.String `tfsdk:"strategy"`
	QuantityToKeep types.Int64  `tfsdk:"quantity_to_keep"`
	Unit           types.String `tfsdk:"unit"`
}

func SpaceDefaultLifecycleRetentionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
	}
}

func SpaceDefaultLifecyclePhaseAttrTypes() map[string]attr.Type {
	retentionType := types.ObjectType{AttrTypes: SpaceDefaultLifecycleRetentionAttrTypes()}
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_with_strategy":       retentionType,
		"tentacle_retention_with_strategy":      retentionType,
	}
}

func (s SpaceDefaultLifecycleSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		MarkdownDescription: "Manages the `" + SpaceDefaultLifecycleName + "` that Octopus creates automatically with a space." +
			"\n\nThe lifecycle cannot be created or deleted through the API: creating this resource adopts the existing lifecycle of the space (the oldest lifecycle of the space, which is the one created with it, even if it has since been renamed), " +
			"and destroying it removes its phases, so that referenced environments can be deleted, and removes it from the Terraform state. The lifecycle itself is deleted together with the space." +
			"\n\nCompatible with Octopus 2025.3 and later.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"space_id": util.ResourceString().
				Required().
				Description("The ID of the space whose default lifecycle is managed.").
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"name": util.ResourceString().
				Optional().Computed().
				Description("The name of the default lifecycle. When omitted, the existing name is kept.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"description": util.ResourceString().
				Optional().Computed().
				Description("The description of the default lifecycle. When omitted, the existing description is kept.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"phases": resourceSchema.ListNestedAttribute{
				Description: "The phases of the default lifecycle. When omitted, the lifecycle has no explicit phases and all environments are available.",
				Optional:    true,
				NestedObject: resourceSchema.NestedAttributeObject{
					Attributes: map[string]resourceSchema.Attribute{
						"id": util.ResourceString().
							Computed().
							Description("The unique ID for this phase.").
							PlanModifiers(stringplanmodifier.UseStateForUnknown()).
							Build(),
						"name": util.ResourceString().Required().Description("The name of this phase.").Build(),
						"automatic_deployment_targets": util.ResourceList(types.StringType).
							Optional().Computed().
							DefaultEmpty().
							Description("Environment IDs in this phase that a release is automatically deployed to when it is eligible for this phase").
							Build(),
						"optional_deployment_targets": util.ResourceList(types.StringType).
							Optional().Computed().
							DefaultEmpty().
							Description("Environment IDs in this phase that a release can be deployed to, but is not automatically deployed to").
							Build(),
						"minimum_environments_before_promotion": util.ResourceInt64().
							Optional().Computed().
							Default(int64(0)).
							Validators(int64validator.Between(0, math.MaxInt32)).
							Description("The number of units required before a release can enter the next phase. If 0, all environments are required.").
							Build(),
						"is_optional_phase": util.ResourceBool().
							Optional().Computed().
							Default(false).
							Description("If false a release must be deployed to this phase before it can be deployed to the next phase.").
							Build(),
						"is_priority_phase": util.ResourceBool().
							Optional().Computed().
							Default(false).
							Description("Deployments will be prioritized in this phase").
							Build(),
						"release_retention_with_strategy":  getSpaceDefaultLifecycleRetentionSchema("Defines the release retention policy for this phase. When omitted, the phase inherits the retention from the lifecycle."),
						"tentacle_retention_with_strategy": getSpaceDefaultLifecycleRetentionSchema("Defines the tentacle retention policy for this phase. When omitted, the phase inherits the retention from the lifecycle."),
					},
				},
			},
			"release_retention_with_strategy":  getSpaceDefaultLifecycleRetentionSchema("Defines the release retention policy of the default lifecycle. When omitted, the space-wide \"Default\" retention policy is used."),
			"tentacle_retention_with_strategy": getSpaceDefaultLifecycleRetentionSchema("Defines the tentacle retention policy of the default lifecycle. When omitted, the space-wide \"Default\" retention policy is used."),
		},
	}
}

func (s SpaceDefaultLifecycleSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

func getSpaceDefaultLifecycleRetentionSchema(description string) resourceSchema.SingleNestedAttribute {
	return resourceSchema.SingleNestedAttribute{
		Description: description,
		Optional:    true,
		Attributes: map[string]resourceSchema.Attribute{
			"strategy": util.ResourceString().
				Required().
				Validators(stringvalidator.OneOf(core.RetentionStrategyDefault, core.RetentionStrategyCount, core.RetentionStrategyForever)).
				Description("How retention will be set. Valid strategies are `Default`, `Forever`, and `Count`." +
					"\n  - `strategy = \"Default\"`, is used if the retention is set by the space-wide default lifecycle retention policy." +
					"\n  - `strategy = \"Forever\"`, is used if items should never be deleted." +
					"\n  - `strategy = \"Count\"`, is used if a specific number of days/releases should be kept.").
				Build(),
			"quantity_to_keep": util.ResourceInt64().
				Optional().
				Validators(int64validator.Between(1, math.MaxInt32)).
				Description("The number of days/releases to keep. Required when `strategy` is `Count`.").
				Build(),
			"unit": util.ResourceString().
				Optional().
				Validators(stringvalidator.OneOfCaseInsensitive(core.RetentionUnitDays, core.RetentionUnitItems)).
				Description("The unit of quantity to keep. Valid units are Days or Items. Required when `strategy` is `Count`.").
				Build(),
		},
		Validators: []validator.Object{
			spaceDefaultLifecycleRetentionValidator{},
		},
	}
}

// spaceDefaultLifecycleRetentionValidator applies the lifecycle retention rules to an optional retention object.
type spaceDefaultLifecycleRetentionValidator struct{}

func (v spaceDefaultLifecycleRetentionValidator) Description(ctx context.Context) string {
	return resourceSchemaRetentionValidator{}.Description(ctx)
}

func (v spaceDefaultLifecycleRetentionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v spaceDefaultLifecycleRetentionValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	resourceSchemaRetentionValidator{}.ValidateObject(ctx, req, resp)
}
