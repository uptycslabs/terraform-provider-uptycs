package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func EventExcludeProfileResource() resource.Resource {
	return &eventExcludeProfileResource{}
}

type eventExcludeProfileResource struct {
	client *uptycs.Client
}

func (r *eventExcludeProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_exclude_profile"
}

func (r *eventExcludeProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *eventExcludeProfileResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"name":          schema.StringAttribute{Optional: true},
			"description":   schema.StringAttribute{Optional: true},
			"priority":      schema.Int64Attribute{Optional: true},
			"resource_type": schema.StringAttribute{Computed: true},
			"platform":      schema.StringAttribute{Optional: true},
			"metadata":      schema.StringAttribute{Optional: true},
		},
	}
}

func (r *eventExcludeProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan EventExcludeProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventExcludeProfileResp, err := r.client.CreateEventExcludeProfile(uptycs.EventExcludeProfile{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		MetadataJSON: plan.Metadata.ValueString(),
		Priority:     int(plan.Priority.ValueInt64()),
		ResourceType: plan.ResourceType.ValueString(),
		Platform:     plan.Platform.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create eventExcludeProfile, unexpected error: "+err.Error(),
		)
		return
	}

	metadataJSON, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.StringValue(eventExcludeProfileResp.ID),
		Name:         types.StringValue(eventExcludeProfileResp.Name),
		Description:  types.StringValue(eventExcludeProfileResp.Description),
		Metadata:     types.StringValue(string(metadataJSON) + "\n"),
		Priority:     types.Int64Value(int64(eventExcludeProfileResp.Priority)),
		ResourceType: types.StringValue(eventExcludeProfileResp.ResourceType),
		Platform:     types.StringValue(eventExcludeProfileResp.Platform),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *eventExcludeProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var eventExcludeProfileID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &eventExcludeProfileID)...)
	eventExcludeProfileResp, err := r.client.GetEventExcludeProfile(uptycs.EventExcludeProfile{
		ID: eventExcludeProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get eventExcludeProfile with ID  "+eventExcludeProfileID+": "+err.Error(),
		)
		return
	}

	metadataJSON, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.StringValue(eventExcludeProfileResp.ID),
		Name:         types.StringValue(eventExcludeProfileResp.Name),
		Description:  types.StringValue(eventExcludeProfileResp.Description),
		Metadata:     types.StringValue(string(metadataJSON) + "\n"),
		Priority:     types.Int64Value(int64(eventExcludeProfileResp.Priority)),
		ResourceType: types.StringValue(eventExcludeProfileResp.ResourceType),
		Platform:     types.StringValue(eventExcludeProfileResp.Platform),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *eventExcludeProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state EventExcludeProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventExcludeProfileID := state.ID.ValueString()

	// Retrieve values from plan
	var plan EventExcludeProfile
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventExcludeProfileResp, err := r.client.UpdateEventExcludeProfile(uptycs.EventExcludeProfile{
		ID:           eventExcludeProfileID,
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		MetadataJSON: plan.Metadata.ValueString(),
		Priority:     int(plan.Priority.ValueInt64()),
		ResourceType: plan.ResourceType.ValueString(),
		Platform:     plan.Platform.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating eventExcludeProfile",
			"Could not create eventExcludeProfile, unexpected error: "+err.Error(),
		)
		return
	}

	metadataJSON, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.StringValue(eventExcludeProfileResp.ID),
		Name:         types.StringValue(eventExcludeProfileResp.Name),
		Description:  types.StringValue(eventExcludeProfileResp.Description),
		Metadata:     types.StringValue(string(metadataJSON) + "\n"),
		Priority:     types.Int64Value(int64(eventExcludeProfileResp.Priority)),
		ResourceType: types.StringValue(eventExcludeProfileResp.ResourceType),
		Platform:     types.StringValue(eventExcludeProfileResp.Platform),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *eventExcludeProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EventExcludeProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventExcludeProfileID := state.ID.ValueString()

	_, err := r.client.DeleteEventExcludeProfile(uptycs.EventExcludeProfile{
		ID: eventExcludeProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete eventExcludeProfile with ID  "+eventExcludeProfileID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r eventExcludeProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
