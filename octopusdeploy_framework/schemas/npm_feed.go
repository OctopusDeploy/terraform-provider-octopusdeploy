package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const npmFeedDescription = "npm feed"

type NpmFeedSchema struct{}

func (n NpmFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages an NPM feed in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"download_attempts":                    GetDownloadAttemptsResourceSchema(),
			"download_retry_backoff_seconds":       GetDownloadRetryBackoffSecondsResourceSchema(),
			"feed_uri":                             GetFeedUriResourceSchema(),
			"id":                                   GetIdResourceSchema(),
			"name":                                 GetNameResourceSchema(true),
			"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
			"password":                             GetPasswordResourceSchema(false),
			"space_id":                             GetSpaceIdResourceSchema(npmFeedDescription),
			"username":                             GetUsernameResourceSchema(false),
		},
	}
}

func (n NpmFeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = NpmFeedSchema{}

type NpmFeedTypeResourceModel struct {
	DownloadAttempts                  types.Int64  `tfsdk:"download_attempts"`
	DownloadRetryBackoffSeconds       types.Int64  `tfsdk:"download_retry_backoff_seconds"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`

	ResourceModel
}
