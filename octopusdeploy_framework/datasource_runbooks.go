package octopusdeploy_framework

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/uritemplates"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &runbooksDataSource{}

type runbooksDataSource struct {
	*Config
}

type runbooksDataSourceModel struct {
	ID          types.String                 `tfsdk:"id"`
	SpaceID     types.String                 `tfsdk:"space_id"`
	ProjectID   types.String                 `tfsdk:"project_id"`
	GitRef      types.String                 `tfsdk:"git_ref"`
	IDs         types.List                   `tfsdk:"ids"`
	PartialName types.String                 `tfsdk:"partial_name"`
	Skip        types.Int64                  `tfsdk:"skip"`
	Take        types.Int64                  `tfsdk:"take"`
	Runbooks    []runbookDataSourceItemModel `tfsdk:"runbooks"`
}

type runbookDataSourceItemModel struct {
	ID                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	ProjectID                      types.String `tfsdk:"project_id"`
	SpaceID                        types.String `tfsdk:"space_id"`
	RunbookProcessID               types.String `tfsdk:"runbook_process_id"`
	PublishedRunbookSnapshotID     types.String `tfsdk:"published_runbook_snapshot_id"`
	MultiTenancyMode               types.String `tfsdk:"multi_tenancy_mode"`
	EnvironmentScope               types.String `tfsdk:"environment_scope"`
	Environments                   types.List   `tfsdk:"environments"`
	DefaultGuidedFailureMode       types.String `tfsdk:"default_guided_failure_mode"`
	ForcePackageDownload           types.Bool   `tfsdk:"force_package_download"`
	RunbookTags                    types.Set    `tfsdk:"runbook_tags"`
	ConnectivityPolicy             types.List   `tfsdk:"connectivity_policy"`
	RunRetentionPolicy             types.List   `tfsdk:"retention_policy"`
	RunRetentionPolicyWithStrategy types.List   `tfsdk:"retention_policy_with_strategy"`
}

func NewRunbooksDataSource() datasource.DataSource {
	return &runbooksDataSource{}
}

func (r *runbooksDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.RunbookDataSourceName)
}

func (r *runbooksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.RunbookSchema{}.GetDatasourceSchema()
}

func (r *runbooksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.Config = DataSourceConfiguration(req, resp)
}

func (r *runbooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data runbooksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceID := data.SpaceID.ValueString()
	if spaceID == "" {
		spaceID = r.Client.GetSpaceID()
	}

	ids := util.GetIds(data.IDs)

	var found []*runbooks.Runbook
	var err error

	// The project scoped endpoint is the only one that filters by project. The
	// space wide endpoint accepts a projectIds parameter but ignores it.
	projectID := data.ProjectID.ValueString()
	gitRef := data.GitRef.ValueString()
	if projectID != "" {
		found, err = listRunbooksByProject(r.Client, spaceID, projectID, gitRef, data)
		found = filterRunbooksByID(found, ids)
	} else {
		found, err = listRunbooksInSpace(r.Client, spaceID, ids, data)
	}

	if err != nil {
		detail := err.Error()
		if projectID != "" && gitRef == "" {
			detail += fmt.Sprintf("\n\nIf project %s is version controlled and stores its runbooks in Git, set git_ref to the branch, tag or commit to read them from.", projectID)
		}
		resp.Diagnostics.AddError("Unable to query runbooks", detail)
		return
	}

	data.Runbooks = []runbookDataSourceItemModel{}
	for _, runbook := range found {
		item, diags := flattenRunbookForDataSource(ctx, runbook)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		data.Runbooks = append(data.Runbooks, item)
	}

	data.ID = types.StringValue(fmt.Sprintf("Runbooks-%s", time.Now().UTC().String()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// listRunbooksByProject expands the project scoped endpoints directly because the
// client's runbooks.List and runbooks.ListGitRunbooks helpers do not support skip.
// Runbooks stored in Git are read from the endpoint that includes the git ref.
func listRunbooksByProject(client newclient.Client, spaceID string, projectID string, gitRef string, data runbooksDataSourceModel) ([]*runbooks.Runbook, error) {
	template, templateParams := runbooksByProjectTemplate(spaceID, projectID, gitRef, data)

	expandedURI, err := client.URITemplateCache().Expand(template, templateParams)
	if err != nil {
		return nil, err
	}

	result, err := newclient.Get[resources.Resources[*runbooks.Runbook]](client.HttpSession(), expandedURI)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// runbooksByProjectTemplate picks the project scoped endpoint and its parameters.
func runbooksByProjectTemplate(spaceID string, projectID string, gitRef string, data runbooksDataSourceModel) (string, map[string]any) {
	template := uritemplates.RunbooksByProject
	templateParams := map[string]any{"spaceId": spaceID, "projectId": projectID}
	if gitRef != "" {
		template = uritemplates.GitRunbooksByProject
		templateParams["gitRef"] = gitRef
	}
	if partialName := data.PartialName.ValueString(); partialName != "" {
		templateParams["partialName"] = partialName
	}
	if skip := data.Skip.ValueInt64(); skip > 0 {
		templateParams["skip"] = skip
	}
	if take := data.Take.ValueInt64(); take > 0 {
		templateParams["take"] = take
	}

	return template, templateParams
}

func listRunbooksInSpace(client newclient.Client, spaceID string, ids []string, data runbooksDataSourceModel) ([]*runbooks.Runbook, error) {
	query := runbooks.RunbooksQuery{
		IDs:         ids,
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	result, err := runbooks.Get(client, spaceID, query)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// filterRunbooksByID narrows a result set to the requested IDs. The project scoped
// endpoint has no ids parameter, so the filtering happens here instead.
func filterRunbooksByID(found []*runbooks.Runbook, ids []string) []*runbooks.Runbook {
	if len(ids) == 0 {
		return found
	}

	filtered := make([]*runbooks.Runbook, 0, len(found))
	for _, runbook := range found {
		if slices.Contains(ids, runbook.GetID()) {
			filtered = append(filtered, runbook)
		}
	}

	return filtered
}

func flattenRunbookForDataSource(ctx context.Context, runbook *runbooks.Runbook) (runbookDataSourceItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	item := runbookDataSourceItemModel{
		ID:                         types.StringValue(runbook.GetID()),
		Name:                       types.StringValue(runbook.Name),
		Description:                types.StringValue(runbook.Description),
		ProjectID:                  types.StringValue(runbook.ProjectID),
		SpaceID:                    types.StringValue(runbook.SpaceID),
		RunbookProcessID:           types.StringValue(runbook.RunbookProcessID),
		PublishedRunbookSnapshotID: types.StringValue(runbook.PublishedRunbookSnapshotID),
		MultiTenancyMode:           types.StringValue(string(runbook.MultiTenancyMode)),
		EnvironmentScope:           types.StringValue(runbook.EnvironmentScope),
		Environments:               util.FlattenStringList(runbook.Environments),
		DefaultGuidedFailureMode:   types.StringValue(runbook.DefaultGuidedFailureMode),
		ForcePackageDownload:       types.BoolValue(runbook.ForcePackageDownload),
	}

	tags, d := types.SetValueFrom(ctx, types.StringType, runbook.RunbookTags)
	diags.Append(d...)
	item.RunbookTags = tags

	item.ConnectivityPolicy = types.ListValueMust(
		types.ObjectType{AttrTypes: schemas.GetConnectivityPolicyObjectType()},
		[]attr.Value{schemas.MapFromConnectivityPolicy(runbook.ConnectivityPolicy)},
	)

	item.RunRetentionPolicy = types.ListValueMust(
		types.ObjectType{AttrTypes: schemas.GetLegacyRunbookRetentionPolicyObjectType()},
		[]attr.Value{schemas.MapFromLegacyRunbookRetentionPolicy(runbook.RunRetentionPolicy)},
	)

	item.RunRetentionPolicyWithStrategy = types.ListValueMust(
		types.ObjectType{AttrTypes: schemas.GetRunbookRetentionPolicyObjectType()},
		[]attr.Value{schemas.MapFromRunbookRetentionPolicy(runbook.RunRetentionPolicy)},
	)

	return item, diags
}
