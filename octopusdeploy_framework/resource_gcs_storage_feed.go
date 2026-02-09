package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type gcsStorageFeedTypeResource struct {
	*Config
}

func NewGcsStorageFeedResource() resource.Resource {
	return &gcsStorageFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &gcsStorageFeedTypeResource{}

func (r *gcsStorageFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("gcs_storage_feed")
}

func (r *gcsStorageFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GcsStorageFeedSchema{}.GetResourceSchema()
}

func (r *gcsStorageFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *gcsStorageFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.GcsStorageFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gcsStorageFeed, err := createGcsStorageResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating GCS feed: %s", gcsStorageFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, gcsStorageFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create GCS feed", err.Error())
		return
	}

	updateDataFromGcsStorageFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.GcsStorageFeed))

	tflog.Info(ctx, fmt.Sprintf("GCS feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *gcsStorageFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.GcsStorageFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading GCS feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "GCS feed"); err != nil {
			resp.Diagnostics.AddError("unable to load GCS feed", err.Error())
		}
		return
	}

	gcsStorageFeed := feed.(*feeds.GcsStorageFeed)
	updateDataFromGcsStorageFeed(data, data.SpaceID.ValueString(), gcsStorageFeed)

	tflog.Info(ctx, fmt.Sprintf("GCS feed read (%s)", gcsStorageFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *gcsStorageFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.GcsStorageFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating GCS feed '%s'", data.ID.ValueString()))

	feed, err := createGcsStorageResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load GCS feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating GCS feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update GCS feed", err.Error())
		return
	}

	updateDataFromGcsStorageFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.GcsStorageFeed))

	tflog.Info(ctx, fmt.Sprintf("GCS feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *gcsStorageFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.GcsStorageFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete GCS feed", err.Error())
		return
	}
}

func createGcsStorageResourceFromData(data *schemas.GcsStorageFeedTypeResourceModel) (*feeds.GcsStorageFeed, error) {
	var oidc *feeds.GoogleContainerRegistryOidcAuthentication

	if data.OidcAuthentication != nil {
		oidc = &feeds.GoogleContainerRegistryOidcAuthentication{
			Audience:    data.OidcAuthentication.Audience.ValueString(),
			SubjectKeys: util.ExpandStringList(data.OidcAuthentication.SubjectKey),
		}
	}

	var serviceAccountKey *core.SensitiveValue
	if !data.ServiceAccountJsonKey.IsNull() && data.ServiceAccountJsonKey.ValueString() != "" {
		serviceAccountKey = core.NewSensitiveValue(data.ServiceAccountJsonKey.ValueString())
	}

	feed, err := feeds.NewGcsStorageFeed(
		data.Name.ValueString(),
		data.UseServiceAccountKey.ValueBool(),
		serviceAccountKey,
		data.Project.ValueString(),
		oidc,
	)

	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()
	feed.SpaceID = data.SpaceID.ValueString()
	feed.OidcAuthentication = oidc

	return feed, nil
}

func updateDataFromGcsStorageFeed(data *schemas.GcsStorageFeedTypeResourceModel, spaceId string, feed *feeds.GcsStorageFeed) {
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	data.UseServiceAccountKey = types.BoolValue(feed.UseServiceAccountKey)

	if feed.Project != "" {
		data.Project = types.StringValue(feed.Project)
	}

	data.ID = types.StringValue(feed.ID)

	if feed.OidcAuthentication != nil {
		data.OidcAuthentication = &schemas.GoogleContainerRegistryOidcAuthenticationResourceModel{
			Audience:   types.StringValue(feed.OidcAuthentication.Audience),
			SubjectKey: util.FlattenStringList(feed.OidcAuthentication.SubjectKeys),
		}
	}
}

func (*gcsStorageFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
