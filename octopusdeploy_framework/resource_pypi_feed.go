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

type pyPiFeedTypeResource struct {
	*Config
}

func NewPyPiFeedResource() resource.Resource {
	return &pyPiFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &pyPiFeedTypeResource{}

func (r *pyPiFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("pypi_feed")
}

func (r *pyPiFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PyPiFeedSchema{}.GetResourceSchema()
}

func (r *pyPiFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *pyPiFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.PyPiFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pyPiFeed, err := createPyPiResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating PyPI feed: %s", pyPiFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, pyPiFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create PyPI feed", err.Error())
		return
	}

	updateDataFromPyPiFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.PyPiFeed))

	tflog.Info(ctx, fmt.Sprintf("PyPI feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pyPiFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.PyPiFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading PyPI feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "pypi feed"); err != nil {
			resp.Diagnostics.AddError("unable to load PyPI feed", err.Error())
		}
		return
	}

	pyPiFeed := feed.(*feeds.PyPiFeed)
	updateDataFromPyPiFeed(data, data.SpaceID.ValueString(), pyPiFeed)

	tflog.Info(ctx, fmt.Sprintf("PyPI feed read (%s)", pyPiFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pyPiFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.PyPiFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating PyPI feed '%s'", data.ID.ValueString()))

	feed, err := createPyPiResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load PyPI feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating PyPI feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update PyPI feed", err.Error())
		return
	}

	updateDataFromPyPiFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.PyPiFeed))

	tflog.Info(ctx, fmt.Sprintf("PyPI feed updated (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pyPiFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.PyPiFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete PyPI feed", err.Error())
		return
	}
}

func createPyPiResourceFromData(data *schemas.PyPiFeedTypeResourceModel) (*feeds.PyPiFeed, error) {
	feed, err := feeds.NewPyPiFeed(data.Name.ValueString(), data.FeedUri.ValueString())
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

func updateDataFromPyPiFeed(data *schemas.PyPiFeedTypeResourceModel, spaceId string, feed *feeds.PyPiFeed) {
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

func (*pyPiFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
