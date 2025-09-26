package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ParentEnvironmentResourceDescription = "This resource manages parent environments in Octopus Deploy."

type ParentEnvironmentSchema struct{}

var _ EntitySchema = ParentEnvironmentSchema{}

type ParentEnvironmentModel struct {
	Name                        types.String                      `tfsdk:"name"`
	SpaceID                     types.String                      `tfsdk:"space_id"`
	Description                 types.String                      `tfsdk:"description"`
	Slug                        types.String                      `tfsdk:"slug"`
	UseGuidedFailure            types.Bool                        `tfsdk:"use_guided_failure"`
	AutomaticDeprovisioningRule *AutomaticDeprovisioningRuleModel `tfsdk:"automatic_deprovisioning_rule"`

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
			"space_id":     GetSpaceIdDatasourceSchema("parent environment", false),
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
						"space_id":    GetSpaceIdDatasourceSchema("parent environment", true),
						"name":        GetReadonlyNameDatasourceSchema(),
						"slug":        GetSlugDatasourceSchema("parent environment", true),
						"description": GetDescriptionDatasourceSchema("parent environment"),
						"use_guided_failure": datasourceSchema.BoolAttribute{
							Description: "Indicates whether guided failure mode is enabled for this parent environment.",
							Computed:    true,
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
			"space_id": schema.StringAttribute{
				Description: "The space ID associated with this parent environment.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of this resource.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of this parent environment.",
				Optional:    true,
				Computed:    true, // API might return computed values
			},
			"slug": schema.StringAttribute{
				Description: "The human-readable unique identifier for this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_guided_failure": schema.BoolAttribute{
				Description: "Indicates whether guided failure mode is enabled for this parent environment.",
				Optional:    true,
				Computed:    true, // Allow it to be computed by the API
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"automatic_deprovisioning_rule": schema.SingleNestedAttribute{
				Description: "Automatic deprovisioning rule for the environment.",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"days":  types.Int64Type,
							"hours": types.Int64Type,
						},
						map[string]attr.Value{
							"days":  types.Int64Value(7),
							"hours": types.Int64Value(0),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"days": schema.Int64Attribute{
						Description: "Number of days to wait before deprovisioning.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(7),
					},
					"hours": schema.Int64Attribute{
						Description: "Number of hours to wait before deprovisioning.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(0),
					},
				},
			},
		},
	}
}
