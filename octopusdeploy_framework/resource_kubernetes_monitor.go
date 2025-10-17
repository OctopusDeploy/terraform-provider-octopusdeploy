package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/kubernetesmonitors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type kubernetesMonitorResource struct {
	*Config
}

func NewKubernetesMonitorResource() resource.Resource {
	return &kubernetesMonitorResource{}
}

func (r *kubernetesMonitorResource) Metadata(
	ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = util.GetTypeName(schemas.KubernetesMonitorResourceName)
}

func (r *kubernetesMonitorResource) Schema(
	ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schemas.KubernetesMonitorSchema{}.GetResourceSchema()
}

func (r *kubernetesMonitorResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *kubernetesMonitorResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	// Check Octopus server version compatibility
	resp.Diagnostics.Append(r.Config.EnsureResourceCompatibilityByVersion(schemas.KubernetesMonitorResourceName, "2025.3")...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data *schemas.KubernetesMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set space ID if not provided
	if data.SpaceID.IsNull() || data.SpaceID.IsUnknown() {
		data.SpaceID = types.StringValue(r.Config.SpaceID)
	}

	// Validate that the machine is a Kubernetes Agent deployment target
	machine, err := machines.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.MachineID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get machine", err.Error())
		return
	}

	if machine.Endpoint == nil || machine.Endpoint.GetCommunicationStyle() != "KubernetesTentacle" {
		resp.Diagnostics.AddError(
			"Invalid machine type",
			"The machine_id must reference a Kubernetes Agent deployment target. The provided machine is not a Kubernetes Agent.",
		)
		return
	}

	// Parse installation ID
	installationID, err := uuid.Parse(data.InstallationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid installation ID", "The installation_id must be a valid UUID")
		return
	}

	// Create the registration command
	command := kubernetesmonitors.NewRegisterKubernetesMonitorCommand(
		&installationID,
		data.MachineID.ValueString(),
	)
	command.SpaceID = data.SpaceID.ValueString()

	if !data.PreserveAuthenticationToken.IsNull() {
		preserveToken := data.PreserveAuthenticationToken.ValueBool()
		command.PreserveAuthenticationToken = &preserveToken
	}

	// Register the Kubernetes monitor
	tflog.Info(ctx, fmt.Sprintf("Creating Kubernetes monitor with installation ID %s", installationID.String()))

	response, err := kubernetesmonitors.Register(r.Config.Client, command)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Kubernetes monitor", err.Error())
		return
	}

	// Map the response to state
	var authToken string
	if response.AuthenticationToken != nil {
		authToken = *response.AuthenticationToken
	}

	schemas.MapFromKubernetesMonitorToState(
		data,
		&response.Resource,
		authToken,
		response.CertificateThumbprint,
	)

	tflog.Info(ctx, fmt.Sprintf("Created Kubernetes monitor with ID %s", data.ID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data *schemas.KubernetesMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading Kubernetes monitor %s", data.ID.ValueString()))

	response, err := kubernetesmonitors.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if apiError, ok := err.(*core.APIError); ok && apiError.StatusCode == 404 {
			tflog.Info(ctx, fmt.Sprintf("Kubernetes monitor %s not found, removing from state", data.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read Kubernetes monitor", err.Error())
		return
	}

	// Map to state, preserving sensitive values
	schemas.MapFromKubernetesMonitorToState(
		data,
		&response.Resource,
		data.AuthenticationToken.ValueString(),   // Preserve existing auth token
		data.CertificateThumbprint.ValueString(), // Preserve existing thumbprint
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesMonitorResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	// Kubernetes monitors cannot be updated - all attributes require replacement
	resp.Diagnostics.AddError(
		"Update not supported",
		"Kubernetes monitors cannot be updated. All changes require replacement.",
	)
}

func (r *kubernetesMonitorResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data *schemas.KubernetesMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting Kubernetes monitor %s", data.ID.ValueString()))

	err := kubernetesmonitors.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if apiError, ok := err.(*core.APIError); ok && apiError.StatusCode == 404 {
			tflog.Info(ctx, fmt.Sprintf("Kubernetes monitor %s already deleted", data.ID.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Failed to delete Kubernetes monitor", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted Kubernetes monitor %s", data.ID.ValueString()))
}

func (r *kubernetesMonitorResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
