package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GcsStorageFeedSchema struct{}

var _ EntitySchema = GcsStorageFeedSchema{}

func (g GcsStorageFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a Google Cloud Storage feed in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"name":     GetNameResourceSchema(true),
			"space_id": GetSpaceIdResourceSchema("GCS storage feed"),
			"use_service_account_key": resourceSchema.BoolAttribute{
				Description: "Whether to use a service account key for authentication",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"service_account_json_key": resourceSchema.StringAttribute{
				Description: "The GCP service account JSON key content",
				Optional:    true,
				Sensitive:   true,
			},
			"project": resourceSchema.StringAttribute{
				Description: "The GCP project ID",
				Optional:    true,
			},
			"oidc_authentication": resourceSchema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]resourceSchema.Attribute{
					"audience": resourceSchema.StringAttribute{
						Description: "Audience for OIDC token",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"subject_keys": GetOidcSubjectKeysSchema("Keys to include in a deployment or runbook. Valid options are space, feed.", false),
				},
			},
		},
	}
}

func (g GcsStorageFeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type GcsStorageFeedTypeResourceModel struct {
	Name                  types.String                                            `tfsdk:"name"`
	SpaceID               types.String                                            `tfsdk:"space_id"`
	UseServiceAccountKey  types.Bool                                              `tfsdk:"use_service_account_key"`
	ServiceAccountJsonKey types.String                                            `tfsdk:"service_account_json_key"`
	Project               types.String                                            `tfsdk:"project"`
	OidcAuthentication    *GoogleContainerRegistryOidcAuthenticationResourceModel `tfsdk:"oidc_authentication"`

	ResourceModel
}
