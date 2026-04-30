package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const subscriptionURITemplate = "/api/{spaceId}/subscriptions{/id}{?skip,take,ids,partialName,spaces}"

type subscriptionApiModel struct {
	ID                            string                                `json:"Id,omitempty"`
	Name                          string                                `json:"Name"`
	Type                          string                                `json:"Type"`
	IsDisabled                    bool                                  `json:"IsDisabled"`
	SpaceID                       string                                `json:"SpaceId,omitempty"`
	EventNotificationSubscription eventNotificationSubscriptionApiModel `json:"EventNotificationSubscription"`
}

type eventNotificationSubscriptionApiModel struct {
	Filter                     subscriptionFilterApiModel `json:"Filter"`
	EmailTeams                 []string                   `json:"EmailTeams"`
	EmailFrequencyPeriod       string                     `json:"EmailFrequencyPeriod,omitempty"`
	EmailPriority              string                     `json:"EmailPriority,omitempty"`
	EmailShowDatesInTimeZoneId string                     `json:"EmailShowDatesInTimeZoneId,omitempty"`
	WebhookURI                 string                     `json:"WebhookURI,omitempty"`
	WebhookTeams               []string                   `json:"WebhookTeams"`
	WebhookTimeout             string                     `json:"WebhookTimeout,omitempty"`
	WebhookHeaderKey           string                     `json:"WebhookHeaderKey,omitempty"`
	WebhookHeaderValue         string                     `json:"WebhookHeaderValue,omitempty"`
}

type subscriptionFilterApiModel struct {
	Users           []string `json:"Users"`
	Projects        []string `json:"Projects"`
	ProjectGroups   []string `json:"ProjectGroups"`
	Environments    []string `json:"Environments"`
	EventGroups     []string `json:"EventGroups"`
	EventCategories []string `json:"EventCategories"`
	EventAgents     []string `json:"EventAgents"`
	Tenants         []string `json:"Tenants"`
	Tags            []string `json:"Tags"`
	DocumentTypes   []string `json:"DocumentTypes"`
}

type subscriptionModel struct {
	Name                          types.String                        `tfsdk:"name"`
	SpaceID                       types.String                        `tfsdk:"space_id"`
	IsDisabled                    types.Bool                          `tfsdk:"is_disabled"`
	EventNotificationSubscription *eventNotificationSubscriptionModel `tfsdk:"event_notification_subscription"`
	schemas.ResourceModel
}

type eventNotificationSubscriptionModel struct {
	Filter                     *subscriptionFilterModel `tfsdk:"filter"`
	EmailTeams                 types.List               `tfsdk:"email_teams"`
	EmailFrequencyPeriod       types.String             `tfsdk:"email_frequency_period"`
	EmailPriority              types.String             `tfsdk:"email_priority"`
	EmailShowDatesInTimeZoneId types.String             `tfsdk:"email_show_dates_in_timezone_id"`
	WebhookURI                 types.String             `tfsdk:"webhook_uri"`
	WebhookTeams               types.List               `tfsdk:"webhook_teams"`
	WebhookTimeout             types.String             `tfsdk:"webhook_timeout"`
	WebhookHeaderKey           types.String             `tfsdk:"webhook_header_key"`
	WebhookHeaderValue         types.String             `tfsdk:"webhook_header_value"`
}

type subscriptionFilterModel struct {
	Users           types.List `tfsdk:"users"`
	Projects        types.List `tfsdk:"projects"`
	ProjectGroups   types.List `tfsdk:"project_groups"`
	Environments    types.List `tfsdk:"environments"`
	EventGroups     types.List `tfsdk:"event_groups"`
	EventCategories types.List `tfsdk:"event_categories"`
	EventAgents     types.List `tfsdk:"event_agents"`
	Tenants         types.List `tfsdk:"tenants"`
	Tags            types.List `tfsdk:"tags"`
	DocumentTypes   types.List `tfsdk:"document_types"`
}

type subscriptionResource struct {
	*Config
}

var _ resource.Resource = &subscriptionResource{}
var _ resource.ResourceWithImportState = &subscriptionResource{}

func NewSubscriptionResource() resource.Resource {
	return &subscriptionResource{}
}

func (r *subscriptionResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.SubscriptionResourceName)
}

func (r *subscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.SubscriptionSchema{}.GetResourceSchema()
}

func (r *subscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *subscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *subscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan subscriptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiModel := expandSubscription(&plan)

	created, err := newclient.Add[subscriptionApiModel](r.Config.Client, subscriptionURITemplate, plan.SpaceID.ValueString(), apiModel)
	if err != nil {
		resp.Diagnostics.AddError("error creating subscription", err.Error())
		return
	}

	flattenSubscription(created, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *subscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state subscriptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscription, err := newclient.GetByID[subscriptionApiModel](r.Config.Client, subscriptionURITemplate, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "subscription"); err != nil {
			resp.Diagnostics.AddError("error reading subscription", err.Error())
		}
		return
	}

	flattenSubscription(subscription, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *subscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan subscriptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiModel := expandSubscription(&plan)

	updated, err := newclient.Update[subscriptionApiModel](r.Config.Client, subscriptionURITemplate, plan.SpaceID.ValueString(), plan.ID.ValueString(), apiModel)
	if err != nil {
		resp.Diagnostics.AddError("error updating subscription", err.Error())
		return
	}

	flattenSubscription(updated, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *subscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state subscriptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := newclient.DeleteByID(r.Config.Client, subscriptionURITemplate, state.SpaceID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("error deleting subscription", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func expandSubscription(model *subscriptionModel) *subscriptionApiModel {
	n := model.EventNotificationSubscription

	apiNotification := eventNotificationSubscriptionApiModel{
		Filter:                     expandSubscriptionFilter(n.Filter),
		EmailFrequencyPeriod:       n.EmailFrequencyPeriod.ValueString(),
		EmailPriority:              n.EmailPriority.ValueString(),
		EmailShowDatesInTimeZoneId: n.EmailShowDatesInTimeZoneId.ValueString(),
		WebhookURI:                 n.WebhookURI.ValueString(),
		WebhookTimeout:             n.WebhookTimeout.ValueString(),
		WebhookHeaderKey:           n.WebhookHeaderKey.ValueString(),
		WebhookHeaderValue:         n.WebhookHeaderValue.ValueString(),
		EmailTeams:                 util.ExpandStringList(n.EmailTeams),
		WebhookTeams:               util.ExpandStringList(n.WebhookTeams),
	}

	return &subscriptionApiModel{
		ID:                            model.ID.ValueString(),
		Name:                          model.Name.ValueString(),
		Type:                          "Event",
		IsDisabled:                    model.IsDisabled.ValueBool(),
		SpaceID:                       model.SpaceID.ValueString(),
		EventNotificationSubscription: apiNotification,
	}
}

func expandSubscriptionFilter(model *subscriptionFilterModel) subscriptionFilterApiModel {
	if model == nil {
		return subscriptionFilterApiModel{}
	}

	return subscriptionFilterApiModel{
		Users:           util.ExpandStringList(model.Users),
		Projects:        util.ExpandStringList(model.Projects),
		ProjectGroups:   util.ExpandStringList(model.ProjectGroups),
		Environments:    util.ExpandStringList(model.Environments),
		EventGroups:     util.ExpandStringList(model.EventGroups),
		EventCategories: util.ExpandStringList(model.EventCategories),
		EventAgents:     util.ExpandStringList(model.EventAgents),
		Tenants:         util.ExpandStringList(model.Tenants),
		Tags:            util.ExpandStringList(model.Tags),
		DocumentTypes:   util.ExpandStringList(model.DocumentTypes),
	}
}

func flattenSubscription(api *subscriptionApiModel, model *subscriptionModel) {
	model.ID = types.StringValue(api.ID)
	model.Name = types.StringValue(api.Name)
	model.SpaceID = types.StringValue(api.SpaceID)
	model.IsDisabled = types.BoolValue(api.IsDisabled)

	if model.EventNotificationSubscription == nil {
		model.EventNotificationSubscription = &eventNotificationSubscriptionModel{}
	}

	n := model.EventNotificationSubscription
	apiN := api.EventNotificationSubscription

	n.EmailFrequencyPeriod = types.StringValue(apiN.EmailFrequencyPeriod)
	n.EmailPriority = types.StringValue(apiN.EmailPriority)
	n.EmailShowDatesInTimeZoneId = types.StringValue(apiN.EmailShowDatesInTimeZoneId)
	n.WebhookTimeout = types.StringValue(apiN.WebhookTimeout)

	n.WebhookURI = util.StringOrNull(apiN.WebhookURI)
	n.WebhookHeaderKey = util.StringOrNull(apiN.WebhookHeaderKey)

	n.EmailTeams = util.FlattenStringList(apiN.EmailTeams)
	n.WebhookTeams = util.FlattenStringList(apiN.WebhookTeams)
	n.Filter = flattenSubscriptionFilter(apiN.Filter, n.Filter)
}

func flattenSubscriptionFilter(api subscriptionFilterApiModel, existing *subscriptionFilterModel) *subscriptionFilterModel {
	allEmpty := len(api.Users) == 0 &&
		len(api.Projects) == 0 &&
		len(api.ProjectGroups) == 0 &&
		len(api.Environments) == 0 &&
		len(api.EventGroups) == 0 &&
		len(api.EventCategories) == 0 &&
		len(api.EventAgents) == 0 &&
		len(api.Tenants) == 0 &&
		len(api.Tags) == 0 &&
		len(api.DocumentTypes) == 0

	if allEmpty && existing == nil {
		return nil
	}

	if existing == nil {
		existing = &subscriptionFilterModel{}
	}

	existing.Users = util.FlattenStringList(api.Users)
	existing.Projects = util.FlattenStringList(api.Projects)
	existing.ProjectGroups = util.FlattenStringList(api.ProjectGroups)
	existing.Environments = util.FlattenStringList(api.Environments)
	existing.EventGroups = util.FlattenStringList(api.EventGroups)
	existing.EventCategories = util.FlattenStringList(api.EventCategories)
	existing.EventAgents = util.FlattenStringList(api.EventAgents)
	existing.Tenants = util.FlattenStringList(api.Tenants)
	existing.Tags = util.FlattenStringList(api.Tags)
	existing.DocumentTypes = util.FlattenStringList(api.DocumentTypes)

	return existing
}
