package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type npmFeedTypeResource struct {
	*Config
}

func NewNpmFeedResource() resource.Resource {
	return &npmFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &npmFeedTypeResource{}

func (r *npmFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("npm_feed")
}

func (r *npmFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.NpmFeedSchema{}.GetResourceSchema()
}

func (r *npmFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *npmFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.NpmFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	npmFeed, err := createNpmResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating NPM feed: %s", npmFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, npmFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create npm feed", err.Error())
		return
	}

	updateDataFromNpmFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.NpmFeed))

	tflog.Info(ctx, fmt.Sprintf("NPM feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *npmFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.NpmFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading NPM feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "npm feed"); err != nil {
			resp.Diagnostics.AddError("unable to load npm feed", err.Error())
		}
		return
	}

	npmFeed := feed.(*feeds.NpmFeed)
	updateDataFromNpmFeed(data, data.SpaceID.ValueString(), npmFeed)

	tflog.Info(ctx, fmt.Sprintf("NPM feed read (%s)", npmFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *npmFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.NpmFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating npm feed '%s'", data.ID.ValueString()))

	feed, err := createNpmResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load npm feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating NPM feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update npm feed", err.Error())
		return
	}

	updateDataFromNpmFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.NpmFeed))

	tflog.Info(ctx, fmt.Sprintf("NPM feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *npmFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.NpmFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete npm feed", err.Error())
		return
	}
}

func createNpmResourceFromData(data *schemas.NpmFeedTypeResourceModel) (*feeds.NpmFeed, error) {
	feed, err := feeds.NewNpmFeed(data.Name.ValueString(), data.FeedUri.ValueString())
	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()
	feed.DownloadAttempts = int(data.DownloadAttempts.ValueInt64())
	feed.DownloadRetryBackoffSeconds = int(data.DownloadRetryBackoffSeconds.ValueInt64())
	feed.FeedURI = data.FeedUri.ValueString()

	if !data.PackageAcquisitionLocationOptions.IsNull() && !data.PackageAcquisitionLocationOptions.IsUnknown() {
		var packageAcquisitionLocationOptions []string
		for _, element := range data.PackageAcquisitionLocationOptions.Elements() {
			packageAcquisitionLocationOptions = append(packageAcquisitionLocationOptions, element.(types.String).ValueString())
		}
		feed.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptions
	}

	feed.Password = core.NewSensitiveValue(data.Password.ValueString())
	feed.SpaceID = data.SpaceID.ValueString()
	feed.Username = data.Username.ValueString()

	return feed, nil
}

func updateDataFromNpmFeed(data *schemas.NpmFeedTypeResourceModel, spaceId string, feed *feeds.NpmFeed) {
	data.DownloadAttempts = types.Int64Value(int64(feed.DownloadAttempts))
	data.DownloadRetryBackoffSeconds = types.Int64Value(int64(feed.DownloadRetryBackoffSeconds))
	data.FeedUri = types.StringValue(feed.FeedURI)
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	if feed.Username != "" {
		data.Username = types.StringValue(feed.Username)
	}

	if feed.PackageAcquisitionLocationOptions != nil {
		packageAcquisitionLocationOptionsList := make([]attr.Value, len(feed.PackageAcquisitionLocationOptions))
		for i, option := range feed.PackageAcquisitionLocationOptions {
			packageAcquisitionLocationOptionsList[i] = types.StringValue(option)
		}

		var packageAcquisitionLocationOptionsListValue, _ = types.ListValue(types.StringType, packageAcquisitionLocationOptionsList)
		data.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptionsListValue
	}

	data.ID = types.StringValue(feed.ID)
}

func (*npmFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
