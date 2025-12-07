package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &rootDataSource{}

type rootDataSource struct {
	*Config
}

// Metadata implements datasource.DataSource.
func (r *rootDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("root")
}

// Read implements datasource.DataSource.
func (r *rootDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data struct {
		ServerVersion types.String `tfsdk:"server_version"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	root, err := client.GetServerRoot(r.Client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Root Data Source",
			"An unexpected error was encountered trying to read the Root data source. "+
				"Please try again later.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Read Root Data Source")

	data.ServerVersion = types.StringValue(root.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rootDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.Config = DataSourceConfiguration(req, resp)
}

// Schema implements datasource.DataSource.
func (r *rootDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.RootSchema{}.GetDatasourceSchema()
}

func NewRootDataSource() datasource.DataSource {
	return &rootDataSource{}
}
