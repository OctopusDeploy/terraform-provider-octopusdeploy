package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const dockerContainerRegistryFeedDescription = "docker container registry feed"

type DockerContainerRegistryFeedSchema struct{}

var _ EntitySchema = DockerContainerRegistryFeedSchema{}

func (d DockerContainerRegistryFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"api_version": resourceSchema.StringAttribute{
				Optional: true,
			},
			"download_attempts":                    GetDownloadAttemptsResourceSchema(),
			"download_retry_backoff_seconds":       GetDownloadRetryBackoffSecondsResourceSchema(),
			"feed_uri":                             GetFeedUriResourceSchema(),
			"id":                                   GetIdResourceSchema(),
			"name":                                 GetNameResourceSchema(true),
			"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
			"password":                             GetPasswordResourceSchema(false),
			"space_id":                             GetSpaceIdResourceSchema(dockerContainerRegistryFeedDescription),
			"username":                             GetUsernameResourceSchema(false),
			"registry_path": resourceSchema.StringAttribute{
				Optional: true,
			},
		},
		Description: "This resource manages a Docker Container Registry in Octopus Deploy.",
	}
}

func (d DockerContainerRegistryFeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type DockerContainerRegistryFeedTypeResourceModel struct {
	APIVersion                        types.String `tfsdk:"api_version"`
	DownloadAttempts                  types.Int64  `tfsdk:"download_attempts"`
	DownloadRetryBackoffSeconds       types.Int64  `tfsdk:"download_retry_backoff_seconds"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`
	RegistryPath                      types.String `tfsdk:"registry_path"`

	ResourceModel
}

// DownloadAttemptsOrDefault returns 5 if downloadAttempts is zero.
// Handles backward compatibility with pre-2025.4 servers that don't return this field.
// Zero indicates missing field (old server), not user-set value (valid range: 1-5).
func DownloadAttemptsOrDefault(downloadAttempts int) int {
	if downloadAttempts == 0 {
		return 5
	}
	return downloadAttempts
}

// DownloadRetryBackoffSecondsOrDefault returns 10 if downloadRetryBackoffSeconds is zero.
// Handles backward compatibility with pre-2025.4 servers that don't return this field.
// Zero indicates missing field (old server), not user-set value.
func DownloadRetryBackoffSecondsOrDefault(downloadRetryBackoffSeconds int) int {
	if downloadRetryBackoffSeconds == 0 {
		return 10
	}
	return downloadRetryBackoffSeconds
}
