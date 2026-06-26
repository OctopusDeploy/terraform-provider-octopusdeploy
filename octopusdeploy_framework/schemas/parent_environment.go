package schemas

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments/v2/environments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas/planmodifiers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ParentEnvironmentResourceDescription = "parent environment"

type ParentEnvironmentSchema struct{}

var _ EntitySchema = ParentEnvironmentSchema{}

type ParentEnvironmentModel struct {
	Name                        types.String `tfsdk:"name"`
	SpaceID                     types.String `tfsdk:"space_id"`
	Description                 types.String `tfsdk:"description"`
	Slug                        types.String `tfsdk:"slug"`
	UseGuidedFailure            types.Bool   `tfsdk:"use_guided_failure"`
	AutomaticDeprovisioningRule types.Object `tfsdk:"automatic_deprovisioning_rule"`
	SortOrder                   types.Int64  `tfsdk:"sort_order"`

	ResourceModel
}

type AutomaticDeprovisioningRuleModel struct {
	Days  types.Int64 `tfsdk:"days"`
	Hours types.Int64 `tfsdk:"hours"`
}

func (p ParentEnvironmentSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			// request
			"space_id":     GetSpaceIdDatasourceSchema(ParentEnvironmentResourceDescription, false),
			"ids":          GetQueryIDsDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),

			// response
			"id": GetIdDatasourceSchema(true),
			"parent_environments": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of parent environments that match the filter(s).",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":          GetIdDatasourceSchema(true),
						"space_id":    GetSpaceIdDatasourceSchema(ParentEnvironmentResourceDescription, true),
						"name":        GetReadonlyNameDatasourceSchema(),
						"slug":        GetSlugDatasourceSchema(ParentEnvironmentResourceDescription, true),
						"description": GetDescriptionDatasourceSchema(ParentEnvironmentResourceDescription),
						"use_guided_failure": datasourceSchema.BoolAttribute{
							Description: "Indicates whether guided failure mode is enabled for this parent environment.",
							Computed:    true,
						},
						"automatic_deprovisioning_rule": datasourceSchema.SingleNestedAttribute{
							Description: "Automatic deprovisioning rule for the environment.",
							Computed:    true,
							Attributes: map[string]datasourceSchema.Attribute{
								"days": datasourceSchema.Int64Attribute{
									Description: "Number of days to wait before deprovisioning.",
									Computed:    true,
								},
								"hours": datasourceSchema.Int64Attribute{
									Description: "Number of hours to wait before deprovisioning.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (p ParentEnvironmentSchema) GetResourceSchema() schema.Schema {
	return schema.Schema{
		Description: util.GetResourceSchemaDescription(ParentEnvironmentResourceDescription),
		Attributes: map[string]schema.Attribute{
			"id": GetIdResourceSchema(),
			"space_id": util.ResourceString().
				Description("The space ID associated with this parent environment.").
				Required().
				Build(),
			"name":        GetNameResourceSchema(true),
			"description": GetDescriptionResourceSchema(ParentEnvironmentResourceDescription),
			"slug": util.ResourceString().
				Description("The human-readable unique identifier for the step.").
				Optional().
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"use_guided_failure": util.ResourceBool().
				Description("Indicates whether guided failure mode is enabled for this parent environment.").
				Optional().
				Computed(). // Allow it to be computed by the API
				PlanModifiers(boolplanmodifier.UseStateForUnknown()).
				Build(),
			// AutomaticDeprovisioningRule is a nested attribute with days and hours
			"automatic_deprovisioning_rule": schema.SingleNestedAttribute{
				Description: "Automatic deprovisioning rule for the environment.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"days": schema.Int64Attribute{
						Description: "Number of days to wait before deprovisioning.",
						Optional:    true,
						Computed:    true,
					},
					"hours": schema.Int64Attribute{
						Description: "Number of hours to wait before deprovisioning.",
						Optional:    true,
						Computed:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					planmodifiers.AllowApiDefaultObject(),
				},
			},
			"sort_order": util.ResourceInt64().
				Description("The sort order of this resource.").
				Optional().
				Computed().
				PlanModifiers(int64planmodifier.UseStateForUnknown()).
				Build(),
		},
	}
}

func MapFromParentEnvironment(ctx context.Context, environment *environments.Environment) ParentEnvironmentTypeResourceModel {
	var env ParentEnvironmentTypeResourceModel
	env.ID = types.StringValue(environment.ID)
	env.SpaceID = types.StringValue(environment.SpaceID)
	env.Slug = types.StringValue(environment.Slug)
	env.Name = types.StringValue(environment.Name)
	env.Description = types.StringValue(environment.Description)
	env.SortOrder = types.Int64Value(int64(environment.SortOrder))
	env.UseGuidedFailure = types.BoolValue(environment.UseGuidedFailure)
	env.Type = types.StringValue(environment.Type)

	env.EnvironmentTags, _ = types.SetValueFrom(ctx, types.StringType, environment.EnvironmentTags)

	return env
}

type ParentEnvironmentTypeResourceModel struct {
	Slug             types.String `tfsdk:"slug"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	SortOrder        types.Int64  `tfsdk:"sort_order"`
	UseGuidedFailure types.Bool   `tfsdk:"use_guided_failure"`
	SpaceID          types.String `tfsdk:"space_id"`
	EnvironmentTags  types.Set    `tfsdk:"environment_tags"`
	Type             types.String `tfsdk:"type"`

	ResourceModel
}
