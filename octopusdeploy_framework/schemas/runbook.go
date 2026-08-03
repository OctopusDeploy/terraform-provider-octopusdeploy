package schemas

import (
	"context"
	"fmt"

	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"

	//datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"regexp"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const RunbookResourceDescription = "runbook"
const RunbookDataSourceName = "runbooks"

var RunbookSchemaAttributeNames = struct {
	ID                          string
	Name                        string
	Description                 string
	ProjectID                   string
	RunbookProcessID            string
	PublishedRunbookSnapshotID  string
	SpaceID                     string
	MultiTenancyMode            string
	ConnectivityPolicy          string
	EnvironmentScope            string
	Environments                string
	DefaultGuidedFailureMode    string
	RetentionPolicy             string
	RetentionPolicyWithStrategy string
	ForcePackageDownload        string
	RunbookTags                 string
}{
	ID:                          "id",
	Name:                        "name",
	Description:                 "description",
	ProjectID:                   "project_id",
	RunbookProcessID:            "runbook_process_id",
	PublishedRunbookSnapshotID:  "published_runbook_snapshot_id",
	SpaceID:                     "space_id",
	MultiTenancyMode:            "multi_tenancy_mode",
	ConnectivityPolicy:          "connectivity_policy",
	EnvironmentScope:            "environment_scope",
	Environments:                "environments",
	DefaultGuidedFailureMode:    "default_guided_failure_mode",
	RetentionPolicy:             "retention_policy",
	RetentionPolicyWithStrategy: "retention_policy_with_strategy",
	ForcePackageDownload:        "force_package_download",
	RunbookTags:                 "runbook_tags",
}

var tenantedDeploymentModeNames = struct {
	Untenanted           string
	TenantedOrUntenanted string
	Tenanted             string
}{
	Untenanted:           "Untenanted",
	TenantedOrUntenanted: "TenantedOrUntenanted",
	Tenanted:             "Tenanted",
}

var tenantedDeploymentModes = []string{
	tenantedDeploymentModeNames.Untenanted,
	tenantedDeploymentModeNames.TenantedOrUntenanted,
	tenantedDeploymentModeNames.Tenanted,
}

var environmentScopeNames = struct {
	All                   string
	Specified             string
	FromProjectLifecycles string
}{
	All:                   "All",
	Specified:             "Specified",
	FromProjectLifecycles: "FromProjectLifecycles",
}

var environmentScopeTypes = []string{
	environmentScopeNames.All,
	environmentScopeNames.Specified,
	environmentScopeNames.FromProjectLifecycles,
}

var defaultGuidedFailureModeNames = struct {
	EnvironmentDefault string
	Off                string
	On                 string
}{
	EnvironmentDefault: "EnvironmentDefault",
	Off:                "Off",
	On:                 "On",
}

var defaultGuidedFailureModes = []string{
	defaultGuidedFailureModeNames.EnvironmentDefault,
	defaultGuidedFailureModeNames.Off,
	defaultGuidedFailureModeNames.On,
}

type RunbookTypeResourceModel struct {
	Name                           types.String `tfsdk:"name"`
	ProjectID                      types.String `tfsdk:"project_id"`
	Description                    types.String `tfsdk:"description"`
	RunbookProcessID               types.String `tfsdk:"runbook_process_id"`
	PublishedRunbookSnapshotID     types.String `tfsdk:"published_runbook_snapshot_id"`
	SpaceID                        types.String `tfsdk:"space_id"`
	MultiTenancyMode               types.String `tfsdk:"multi_tenancy_mode"`
	ConnectivityPolicy             types.List   `tfsdk:"connectivity_policy"`
	EnvironmentScope               types.String `tfsdk:"environment_scope"`
	Environments                   types.List   `tfsdk:"environments"`
	DefaultGuidedFailureMode       types.String `tfsdk:"default_guided_failure_mode"`
	RunRetentionPolicy             types.List   `tfsdk:"retention_policy"`
	RunRetentionPolicyWithStrategy types.List   `tfsdk:"retention_policy_with_strategy"`
	ForcePackageDownload           types.Bool   `tfsdk:"force_package_download"`
	RunbookTags                    types.Set    `tfsdk:"runbook_tags"`

	ResourceModel
}

type RunbookConnectivityPolicyModel struct {
	AllowDeploymentsToNoTargets types.Bool   `tfsdk:"allow_deployments_to_no_targets"`
	ExcludeUnhealthyTargets     types.Bool   `tfsdk:"exclude_unhealthy_targets"`
	SkipMachineBehavior         types.String `tfsdk:"skip_machine_behaviour"`
	TargetRoles                 types.List   `tfsdk:"target_roles"`
}

type RunbookSchema struct{}

func (r RunbookSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing Octopus Deploy runbooks.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id": GetIdDatasourceSchema(true),
			"space_id": util.DataSourceString().
				Optional().
				Description("A Space ID to filter by. Will revert what is specified on the provider if not set.").
				Build(),
			RunbookSchemaAttributeNames.ProjectID: util.DataSourceString().
				Optional().
				Description("A project ID to filter by. When set, only runbooks belonging to that project are returned.").
				Build(),
			"ids":          GetQueryIDsDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),
			"runbooks":     getRunbooksDatasourceAttribute(),
		},
	}
}

func getRunbooksDatasourceAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "A list of runbooks that match the filter(s).",
		Computed:    true,
		Optional:    false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				RunbookSchemaAttributeNames.ID:   util.DataSourceString().Computed().Build(),
				RunbookSchemaAttributeNames.Name: util.DataSourceString().Computed().Description("The name of the runbook.").Build(),
				RunbookSchemaAttributeNames.Description: util.DataSourceString().Computed().
					Description("The description of the runbook.").Build(),
				RunbookSchemaAttributeNames.ProjectID: util.DataSourceString().Computed().
					Description("The project that this runbook belongs to.").Build(),
				RunbookSchemaAttributeNames.SpaceID: util.DataSourceString().Computed().
					Description("The space that this runbook belongs to.").Build(),
				RunbookSchemaAttributeNames.RunbookProcessID: util.DataSourceString().Computed().
					Description("The runbook process ID.").Build(),
				RunbookSchemaAttributeNames.PublishedRunbookSnapshotID: util.DataSourceString().Computed().
					Description("The published snapshot ID.").Build(),
				RunbookSchemaAttributeNames.MultiTenancyMode: util.DataSourceString().Computed().
					Description("The tenanted deployment mode of the runbook.").Build(),
				RunbookSchemaAttributeNames.EnvironmentScope: util.DataSourceString().Computed().
					Description("Determines how the runbook is scoped to environments.").Build(),
				RunbookSchemaAttributeNames.Environments: util.DataSourceList(types.StringType).Computed().
					Description(fmt.Sprintf("When %s is \"%s\", the environments the runbook can be run against.", RunbookSchemaAttributeNames.EnvironmentScope, environmentScopeNames.Specified)).Build(),
				RunbookSchemaAttributeNames.DefaultGuidedFailureMode: util.DataSourceString().Computed().
					Description("The runbook guided failure mode.").Build(),
				RunbookSchemaAttributeNames.ForcePackageDownload: util.DataSourceBool().Computed().
					Description("Whether packages are re-downloaded.").Build(),
				RunbookSchemaAttributeNames.RunbookTags: util.DataSourceSet(types.StringType).Computed().
					Description("The tags associated with this runbook.").Build(),
				RunbookSchemaAttributeNames.ConnectivityPolicy:          getRunbookConnectivityPolicyDatasourceAttribute(),
				RunbookSchemaAttributeNames.RetentionPolicy:             getLegacyRunbookRetentionPolicyDatasourceAttribute(),
				RunbookSchemaAttributeNames.RetentionPolicyWithStrategy: getRunbookRetentionPolicyDatasourceAttribute(),
			},
		},
	}
}

func getRunbookConnectivityPolicyDatasourceAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "The connectivity policy of the runbook.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				runbookConnectivityPolicySchemeAttributeNames.AllowDeploymentsToNoTargets: util.DataSourceBool().Computed().Build(),
				runbookConnectivityPolicySchemeAttributeNames.ExcludeUnhealthyTargets:     util.DataSourceBool().Computed().Build(),
				runbookConnectivityPolicySchemeAttributeNames.SkipMachineBehavior:         util.DataSourceString().Computed().Build(),
				runbookConnectivityPolicySchemeAttributeNames.TargetRoles:                 util.DataSourceList(types.StringType).Computed().Build(),
			},
		},
	}
}

func getLegacyRunbookRetentionPolicyDatasourceAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "The runbook retention policy, in the form kept for backwards compatibility.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				legacyRunbookRetentionPolicySchemeAttributeNames.QuantityToKeep: util.DataSourceInt64().Computed().
					Description("How many runs to keep per environment.").Build(),
				legacyRunbookRetentionPolicySchemeAttributeNames.ShouldKeepForever: util.DataSourceBool().Computed().
					Description("Whether runs are never deleted.").Build(),
			},
		},
	}
}

func getRunbookRetentionPolicyDatasourceAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "The runbook retention policy, including its strategy.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				runbookRetentionPolicySchemeAttributeNames.QuantityToKeep: util.DataSourceInt64().Computed().
					Description("The number of runs or days of runs kept, depending on the unit.").Build(),
				runbookRetentionPolicySchemeAttributeNames.Strategy: util.DataSourceString().Computed().
					Description("How retention is set. One of `Default`, `Forever` or `Count`.").Build(),
				runbookRetentionPolicySchemeAttributeNames.Unit: util.DataSourceString().Computed().
					Description("The unit of the quantity kept, either `Items` or `Days`.").Build(),
			},
		},
	}
}

var _ EntitySchema = RunbookSchema{}

func (r RunbookSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(RunbookResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			RunbookSchemaAttributeNames.ID: GetIdResourceSchema(),
			RunbookSchemaAttributeNames.Name: resourceSchema.StringAttribute{
				Description: "The name of the runbook in Octopus Deploy. This name must be unique.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`\S+`),
						"expected value to not be an empty string or whitespace",
					),
				},
			},
			RunbookSchemaAttributeNames.Description: GetDescriptionResourceSchema(RunbookResourceDescription),
			RunbookSchemaAttributeNames.ProjectID: resourceSchema.StringAttribute{
				Description: "The project that this runbook belongs to.",
				Required:    true,
			},
			RunbookSchemaAttributeNames.RunbookProcessID: resourceSchema.StringAttribute{
				Description: "The runbook process ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.PublishedRunbookSnapshotID: resourceSchema.StringAttribute{
				Description: "The published snapshot ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.SpaceID: GetSpaceIdResourceSchema(RunbookResourceDescription),
			RunbookSchemaAttributeNames.MultiTenancyMode: resourceSchema.StringAttribute{
				Description: fmt.Sprintf("The tenanted deployment mode of the runbook. Valid modes are %s", strings.Join(util.Map(tenantedDeploymentModes, func(item string) string { return fmt.Sprintf("`%s`", item) }), ", ")),
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(tenantedDeploymentModes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.EnvironmentScope: resourceSchema.StringAttribute{
				Description: "Determines how the runbook is scoped to environments.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(environmentScopeTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.Environments: resourceSchema.ListAttribute{
				Description: fmt.Sprintf("When %s is set to \"%s\", this is the list of environments the runbook can be run against.", RunbookSchemaAttributeNames.EnvironmentScope, environmentScopeNames.Specified),
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.DefaultGuidedFailureMode: resourceSchema.StringAttribute{
				Description: "Sets the runbook guided failure mode.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(defaultGuidedFailureModes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.ForcePackageDownload: resourceSchema.BoolAttribute{
				Description: "Whether to force packages to be re-downloaded or not.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.RunbookTags: resourceSchema.SetAttribute{
				Description: "A list of tags associated with this runbook. Valid tags must be defined in a tag set with the 'Runbook' scope enabled.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]resourceSchema.Block{
			RunbookSchemaAttributeNames.ConnectivityPolicy: resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: getConnectivityPolicySchema(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			RunbookSchemaAttributeNames.RetentionPolicy: resourceSchema.ListNestedBlock{
				Description:        "Sets the runbook retention policy.",
				DeprecationMessage: "Runbook Retention Policies will soon require strategy",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: getLegacyRunbookRetentionPolicySchema(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			RunbookSchemaAttributeNames.RetentionPolicyWithStrategy: resourceSchema.ListNestedBlock{
				Description: "Sets the runbook retention policy with strategy.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: getRunbookRetentionPolicySchema(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.ConflictsWith(path.MatchRoot(RunbookSchemaAttributeNames.RetentionPolicy)),
				},
			},
		},
	}
}

func (data *RunbookTypeResourceModel) RefreshFromApiResponse(ctx context.Context, runbook *runbooks.Runbook) diag.Diagnostics {
	var diags diag.Diagnostics

	if runbook == nil {
		return diags
	}

	data.ID = types.StringValue(runbook.ID)
	data.Name = types.StringValue(runbook.Name)
	data.ProjectID = types.StringValue(runbook.ProjectID)
	data.Description = types.StringValue(runbook.Description)
	data.RunbookProcessID = types.StringValue(runbook.RunbookProcessID)
	data.PublishedRunbookSnapshotID = types.StringValue(runbook.PublishedRunbookSnapshotID)
	data.SpaceID = types.StringValue(runbook.SpaceID)
	data.MultiTenancyMode = types.StringValue(string(runbook.MultiTenancyMode))
	data.EnvironmentScope = types.StringValue(runbook.EnvironmentScope)
	data.Environments = util.FlattenStringList(runbook.Environments)
	data.DefaultGuidedFailureMode = types.StringValue(runbook.DefaultGuidedFailureMode)
	data.ForcePackageDownload = types.BoolValue(runbook.ForcePackageDownload)
	if !data.ConnectivityPolicy.IsNull() {
		result, d := types.ListValueFrom(
			ctx,
			types.ObjectType{AttrTypes: GetConnectivityPolicyObjectType()},
			[]attr.Value{MapFromConnectivityPolicy(runbook.ConnectivityPolicy)},
		)
		diags.Append(d...)
		data.ConnectivityPolicy = result
	} /*else {
		data.ConnectivityPolicy = types.ListValueMust(
			types.ObjectType{AttrTypes: GetConnectivityPolicyObjectType()},
			[]attr.Value{MapFromConnectivityPolicy(GetDefaultConnectivityPolicy())},
		)
	}*/
	if !data.RunRetentionPolicy.IsNull() {
		result, d := types.ListValueFrom(
			ctx,
			types.ObjectType{AttrTypes: GetLegacyRunbookRetentionPolicyObjectType()},
			[]attr.Value{MapFromLegacyRunbookRetentionPolicy(runbook.RunRetentionPolicy)},
		)
		diags.Append(d...)
		data.RunRetentionPolicy = result
	}

	data.RunbookTags, _ = types.SetValueFrom(ctx, types.StringType, runbook.RunbookTags)

	return diags
}
