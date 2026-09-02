package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ApprovalRuleResourceDescription = "approval rule"
const ApprovalRuleDataSourceDescription = "approval rules"

type ApprovalRuleSchema struct{}

var _ EntitySchema = ApprovalRuleSchema{}

// ApprovalRuleTagScopeModel represents a tag-based scope entry for an approval rule.
type ApprovalRuleTagScopeModel struct {
	ProjectTags     types.Set `tfsdk:"project_tags"`
	EnvironmentTags types.Set `tfsdk:"environment_tags"`
}

// ApprovalRuleIdScopeModel represents an ID-based scope entry for an approval rule.
type ApprovalRuleIdScopeModel struct {
	ProjectID      types.String `tfsdk:"project_id"`
	EnvironmentIDs types.List   `tfsdk:"environment_ids"`
}

// ApprovalRuleResourceModel represents the octopusdeploy_approval_rule resource defined in Terraform configuration.
type ApprovalRuleResourceModel struct {
	SpaceID                  types.String                  `tfsdk:"space_id"`
	Name                     types.String                  `tfsdk:"name"`
	Description              types.String                  `tfsdk:"description"`
	ScopingStrategy          types.String                  `tfsdk:"scoping_strategy"`
	MinimumApproversRequired types.Int64                   `tfsdk:"minimum_approvers_required"`
	AllowSelfApproval        types.Bool                    `tfsdk:"allow_self_approval"`
	IsDisabled               types.Bool                    `tfsdk:"is_disabled"`
	ApprovingUserIDs         types.List                    `tfsdk:"approving_user_ids"`
	ApprovingTeamIDs         types.List                    `tfsdk:"approving_team_ids"`
	TagScopes                []ApprovalRuleTagScopeModel `tfsdk:"tag_scopes"`
	IdScopes                 []ApprovalRuleIdScopeModel  `tfsdk:"id_scopes"`

	ResourceModel
}

// ApprovalRulesDataSourceModel represents the octopusdeploy_approval_rules data source defined in Terraform configuration.
type ApprovalRulesDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	SpaceID          types.String `tfsdk:"space_id"`
	PartialName      types.String `tfsdk:"partial_name"`
	Skip             types.Int64  `tfsdk:"skip"`
	Take             types.Int64  `tfsdk:"take"`
	ApprovalRules types.List   `tfsdk:"approval_rules"`
}

func (a ApprovalRuleSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages approval rules in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ApprovalRuleResourceDescription),
			"name":     GetNameResourceSchema(true),
			"description": resourceSchema.StringAttribute{
				Description: "The description of this approval rule.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scoping_strategy": resourceSchema.StringAttribute{
				Description: "The scoping strategy used by this approval rule. Valid values are `\"Tag\"` or `\"Id\"`.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Tag", "Id"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"minimum_approvers_required": resourceSchema.Int64Attribute{
				Description: "The minimum number of approvers required for this approval rule.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"allow_self_approval": resourceSchema.BoolAttribute{
				Description: "Whether the deployment creator may approve their own deployment.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_disabled": resourceSchema.BoolAttribute{
				Description: "Whether the policy is disabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"approving_user_ids": resourceSchema.ListAttribute{
				Description: "A list of user IDs that are eligible to approve deployments under this policy.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"approving_team_ids": resourceSchema.ListAttribute{
				Description: "A list of team IDs that are eligible to approve deployments under this policy.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"tag_scopes": resourceSchema.ListNestedAttribute{
				Description: "A list of tag-based scopes this approval rule applies to. Used when `scoping_strategy` is `\"Tag\"`.",
				Optional:    true,
				NestedObject: resourceSchema.NestedAttributeObject{
					Attributes: map[string]resourceSchema.Attribute{
						"project_tags": resourceSchema.SetAttribute{
							Description: "A set of project tags for this scope, specified as tag IDs (e.g. `TagSets-1/Tags-1`) or canonical tag names (e.g. `TagSet/Tag`). Canonical names are resolved to their tag IDs, which are what gets stored in state.",
							Optional:    true,
							Computed:    true,
							ElementType: types.StringType,
						},
						"environment_tags": resourceSchema.SetAttribute{
							Description: "A set of environment tags for this scope, specified as tag IDs (e.g. `TagSets-1/Tags-1`) or canonical tag names (e.g. `TagSet/Tag`). Canonical names are resolved to their tag IDs, which are what gets stored in state.",
							Optional:    true,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"id_scopes": resourceSchema.ListNestedAttribute{
				Description: "A list of ID-based scopes this approval rule applies to. Used when `scoping_strategy` is `\"Id\"`.",
				Optional:    true,
				NestedObject: resourceSchema.NestedAttributeObject{
					Attributes: map[string]resourceSchema.Attribute{
						"project_id": resourceSchema.StringAttribute{
							Description: "The project ID for this scope.",
							Required:    true,
						},
						"environment_ids": resourceSchema.ListAttribute{
							Description: "A list of environment IDs for this scope.",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (a ApprovalRuleSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: util.GetDataSourceDescription(ApprovalRuleDataSourceDescription),
		Attributes: map[string]datasourceSchema.Attribute{
			"id":           GetIdDatasourceSchema(true),
			"space_id":     GetSpaceIdDatasourceSchema(ApprovalRuleResourceDescription, false),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),
			"approval_rules": datasourceSchema.ListNestedAttribute{
				Description: "A list of approval rules that match the filter(s).",
				Computed:    true,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":       GetIdDatasourceSchema(true),
						"space_id": GetSpaceIdDatasourceSchema(ApprovalRuleResourceDescription, true),
						"name":     GetReadonlyNameDatasourceSchema(),
						"description": datasourceSchema.StringAttribute{
							Description: "The description of this approval rule.",
							Computed:    true,
						},
						"scoping_strategy": datasourceSchema.StringAttribute{
							Description: "The scoping strategy used by this approval rule. Valid values are `\"Tag\"` or `\"Id\"`.",
							Computed:    true,
						},
						"minimum_approvers_required": datasourceSchema.Int64Attribute{
							Description: "The minimum number of approvers required for this approval rule.",
							Computed:    true,
						},
						"allow_self_approval": datasourceSchema.BoolAttribute{
							Description: "Whether the deployment creator may approve their own deployment.",
							Computed:    true,
						},
						"is_disabled": datasourceSchema.BoolAttribute{
							Description: "Whether the policy is disabled.",
							Computed:    true,
						},
						"approving_user_ids": datasourceSchema.ListAttribute{
							Description: "A list of user IDs that are eligible to approve deployments under this policy.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"approving_team_ids": datasourceSchema.ListAttribute{
							Description: "A list of team IDs that are eligible to approve deployments under this policy.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"tag_scopes": datasourceSchema.ListNestedAttribute{
							Description: "A list of tag-based scopes this approval rule applies to.",
							Computed:    true,
							NestedObject: datasourceSchema.NestedAttributeObject{
								Attributes: map[string]datasourceSchema.Attribute{
									"project_tags": datasourceSchema.ListAttribute{
										Description: "A list of project tags for this scope.",
										Computed:    true,
										ElementType: types.StringType,
									},
									"environment_tags": datasourceSchema.ListAttribute{
										Description: "A list of environment tags for this scope.",
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
						"id_scopes": datasourceSchema.ListNestedAttribute{
							Description: "A list of ID-based scopes this approval rule applies to.",
							Computed:    true,
							NestedObject: datasourceSchema.NestedAttributeObject{
								Attributes: map[string]datasourceSchema.Attribute{
									"project_id": datasourceSchema.StringAttribute{
										Description: "The project ID for this scope.",
										Computed:    true,
									},
									"environment_ids": datasourceSchema.ListAttribute{
										Description: "A list of environment IDs for this scope.",
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
