package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var _ datasource.DataSource = &gitCredentialsDataSource{}

type gitCredentialsDataSource struct {
	*Config
}

type gitCredentialsDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	SpaceID        types.String `tfsdk:"space_id"`
	Name           types.String `tfsdk:"name"`
	Skip           types.Int64  `tfsdk:"skip"`
	Take           types.Int64  `tfsdk:"take"`
	GitCredentials types.List   `tfsdk:"git_credentials"`
}

type GitCredentialDatasourceModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Username    types.String `tfsdk:"username"`

	RepositoryRestrictions *gitCredentialRepositoryRestrictionDataSourceModel `tfsdk:"repository_restrictions"`
}

type gitCredentialRepositoryRestrictionDataSourceModel struct {
	Enabled             types.Bool `tfsdk:"enabled"`
	AllowedRepositories types.Set  `tfsdk:"allowed_repositories"`
}

func NewGitCredentialsDataSource() datasource.DataSource {
	return &gitCredentialsDataSource{}
}

func (g *gitCredentialsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.GitCredentialDatasourceName)
}

func (g *gitCredentialsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.GitCredentialSchema{}.GetDatasourceSchema()
}

func (g *gitCredentialsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	g.Config = DataSourceConfiguration(req, resp)
}

func (g *gitCredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data gitCredentialsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := credentials.Query{
		Name: data.Name.ValueString(),
		Skip: int(data.Skip.ValueInt64()),
		Take: int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "git credentials", query)

	spaceID := data.SpaceID.ValueString()

	existingGitCredentials, err := credentials.Get(g.Client, spaceID, query)
	if err != nil {
		resp.Diagnostics.AddError("Unable to query git credentials", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "git credentials", len(existingGitCredentials.Items))

	flattenedGitCredentials := make([]GitCredentialDatasourceModel, 0, len(existingGitCredentials.Items))
	for _, gitCredential := range existingGitCredentials.Items {
		flattenedGitCredential := FlattenGitCredential(ctx, gitCredential)
		flattenedGitCredentials = append(flattenedGitCredentials, *flattenedGitCredential)
	}

	gitCredentialsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: GetGitCredentialAttrTypes()}, flattenedGitCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.GitCredentials = gitCredentialsList

	data.ID = types.StringValue(fmt.Sprintf("GitCredentials-%s - new sdk", time.Now().UTC().String()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func GetGitCredentialAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"space_id":                types.StringType,
		"name":                    types.StringType,
		"description":             types.StringType,
		"type":                    types.StringType,
		"username":                types.StringType,
		"repository_restrictions": types.ObjectType{AttrTypes: GetGitCredentialRepositoryRestrictionsAttrTypes()},
	}
}

func GetGitCredentialRepositoryRestrictionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":              types.BoolType,
		"allowed_repositories": types.SetType{ElemType: types.StringType},
	}
}

func FlattenGitCredential(ctx context.Context, credential *credentials.Resource) *GitCredentialDatasourceModel {
	if credential == nil {
		return nil
	}

	model := &GitCredentialDatasourceModel{
		ID:          types.StringValue(credential.GetID()),
		SpaceID:     types.StringValue(credential.SpaceID),
		Name:        types.StringValue(credential.Name),
		Description: types.StringValue(credential.Description),
		Type:        types.StringValue(string(credential.Details.Type())),
	}

	if usernamePassword, ok := credential.Details.(*credentials.UsernamePassword); ok {
		model.Username = types.StringValue(usernamePassword.Username)
	}

	if credential.RepositoryRestrictions != nil {
		var allowedRepositories = make([]string, 0, len(credential.RepositoryRestrictions.AllowedRepositories))
		for _, id := range credential.RepositoryRestrictions.AllowedRepositories {
			allowedRepositories = append(allowedRepositories, id)
		}
		repositoriesList, _ := types.SetValueFrom(
			ctx,
			types.StringType,
			allowedRepositories,
		)

		model.RepositoryRestrictions = &gitCredentialRepositoryRestrictionDataSourceModel{
			Enabled:             types.BoolValue(credential.RepositoryRestrictions.Enabled),
			AllowedRepositories: repositoriesList,
		}
	} else { //Default to disabled if resource doesn't have it
		model.RepositoryRestrictions = &gitCredentialRepositoryRestrictionDataSourceModel{
			Enabled:             types.BoolValue(false),
			AllowedRepositories: types.SetValueMust(types.StringType, make([]attr.Value, 0)),
		}
	}

	return model
}
