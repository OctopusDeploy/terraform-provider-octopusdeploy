package octopusdeploy_framework

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/ratelimitingpolicies"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type builtInRateLimitingPolicyResource struct {
	*Config
}

func NewBuiltInRateLimitingPolicyResource() resource.Resource {
	return &builtInRateLimitingPolicyResource{}
}

var _ resource.ResourceWithImportState = &builtInRateLimitingPolicyResource{}

// The three built-in policies have no persisted slugs yet, so we address them by these hardcoded pseudo-slugs. If/when
// real, persisted slugs are introduced this will no longer be necessary.
var builtInPolicySlugs = []string{"anon", "user", "agent"}

func scopeFromBuiltInSlug(slug string) (ratelimitingpolicies.RateLimitingPolicyScopeType, bool) {
	switch slug {
	case "anon":
		return ratelimitingpolicies.Unauthenticated, true
	case "user":
		return ratelimitingpolicies.AuthenticatedHuman, true
	case "agent":
		return ratelimitingpolicies.AuthenticatedAgent, true
	default:
		return 0, false
	}
}

func (r *builtInRateLimitingPolicyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("built_in_rate_limiting_policy")
}

func (r *builtInRateLimitingPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.BuiltInRateLimitingPolicySchema{}.GetResourceSchema()
}

func (r *builtInRateLimitingPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)

	if r.Config != nil {
		resp.Diagnostics.Append(r.Config.EnsureResourceCompatibilityByVersion("built_in_rate_limiting_policy", "2026.3")...)
	}
}

// findPolicyBySlug resolves a pseudo-slug to its built-in policy by listing all policies and matching via the scope.
func (r *builtInRateLimitingPolicyResource) findPolicyBySlug(slug string) (*ratelimitingpolicies.RateLimitingPolicy, error) {
	scope, ok := scopeFromBuiltInSlug(slug)
	if !ok {
		return nil, fmt.Errorf("'%s' is not a known policy slug; valid values are: %v", slug, builtInPolicySlugs)
	}

	list, err := ratelimitingpolicies.List(r.Client, ratelimitingpolicies.ListRateLimitingPoliciesRequest{Take: math.MaxInt32})
	if err != nil {
		// An instance that does not serve this endpoint answers 404 with no JSON body, which
		// the client turns into an APIError carrying no message at all. Left alone that reaches
		// the practitioner as "Octopus API error:  []", so name the likely cause instead.
		var apiErr *core.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorMessage == "" && len(apiErr.Errors) == 0 {
			return nil, fmt.Errorf("the rate limiting policies endpoint is not available on this Octopus instance; it requires Octopus Server 2026.3 or later and is not offered on every hosting type")
		}
		return nil, fmt.Errorf("unable to list rate limiting policies: %w", err)
	}

	for i := range list.Items {
		if policy := &list.Items[i]; policy.IsBuiltIn && policy.ScopeType == scope {
			return policy, nil
		}
	}

	return nil, nil
}

// updateStateFromPolicy copies the server's view of a policy into the Terraform state model.
func updateStateFromPolicy(data *schemas.BuiltInRateLimitingPolicyResourceModel, policy *ratelimitingpolicies.RateLimitingPolicy) {
	data.ID = types.StringValue(policy.ID)
	data.Name = types.StringValue(policy.Name)
	data.ScopeType = types.StringValue(policy.ScopeType.String())
	data.IsEnabled = types.BoolValue(policy.IsEnabled)
	data.RequestsPerMinute = types.Int64Value(int64(policy.RequestsPerMinute))
	data.BurstLimit = types.Int64Value(int64(policy.BurstLimit))
	data.AuditMode = types.BoolValue(policy.AuditMode)
}

// Create resolves the pseudo-slug to a built-in policy and applies the desired settings via Modify, as creating new
// built-in rate limiting policies is not supported.
func (r *builtInRateLimitingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.BuiltInRateLimitingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.findPolicyBySlug(data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to look up rate limiting policy", err.Error())
		return
	}
	if policy == nil {
		resp.Diagnostics.AddError("Rate limiting policy not found", fmt.Sprintf("No built-in rate limiting policy found for slug '%s'.", data.Slug.ValueString()))
		return
	}

	updated, err := r.modify(data, policy)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create rate limiting policy", err.Error())
		return
	}

	updateStateFromPolicy(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *builtInRateLimitingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.BuiltInRateLimitingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.findPolicyBySlug(data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read rate limiting policy", err.Error())
		return
	}
	if policy == nil {
		// If the policy no longer exists somehow, drop it from state
		resp.State.RemoveResource(ctx)
		return
	}

	updateStateFromPolicy(&data, policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *builtInRateLimitingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schemas.BuiltInRateLimitingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.findPolicyBySlug(data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to look up rate limiting policy", err.Error())
		return
	}
	if policy == nil {
		resp.Diagnostics.AddError("Rate limiting policy not found", fmt.Sprintf("No built-in rate limiting policy found for slug '%s'.", data.Slug.ValueString()))
		return
	}

	updated, err := r.modify(data, policy)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update rate limiting policy", err.Error())
		return
	}

	updateStateFromPolicy(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete cannot remove a built-in policy, so it only drops the resource from state and warns that the policy's
// last-applied settings remain on the server.
func (r *builtInRateLimitingPolicyResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Built-in rate limiting policy removed from state",
		"Built-in rate limiting policies cannot be deleted. The policy and its last-applied settings remain on the server.",
	)
}

// ImportState identifies a policy by slug (anon/user/agent).
func (r *builtInRateLimitingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, resp)
}

func (r *builtInRateLimitingPolicyResource) modify(data schemas.BuiltInRateLimitingPolicyResourceModel, policy *ratelimitingpolicies.RateLimitingPolicy) (*ratelimitingpolicies.RateLimitingPolicy, error) {
	return ratelimitingpolicies.Modify(r.Client, ratelimitingpolicies.ModifyRateLimitingPolicyCommand{
		ID:                policy.ID,
		Name:              policy.Name,
		ScopeType:         policy.ScopeType,
		IsEnabled:         data.IsEnabled.ValueBool(),
		RequestsPerMinute: int(data.RequestsPerMinute.ValueInt64()),
		BurstLimit:        int(data.BurstLimit.ValueInt64()),
		AuditMode:         data.AuditMode.ValueBool(),
	})
}
