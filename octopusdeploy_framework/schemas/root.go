package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type RootSchema struct{}

// GetDatasourceSchema implements EntitySchema.
func (r RootSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			"server_version": datasourceSchema.StringAttribute{
				Description: "The version of the connected Octopus Deploy server.",
				Computed:    true,
			},
		},
	}
}

// GetResourceSchema implements EntitySchema.
func (r RootSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{}
}

var _ EntitySchema = RootSchema{}
