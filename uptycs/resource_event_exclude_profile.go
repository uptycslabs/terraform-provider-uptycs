package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &eventExcludeProfileResource{}
	_ resource.ResourceWithConfigure   = &eventExcludeProfileResource{}
	_ resource.ResourceWithImportState = &eventExcludeProfileResource{}
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

func (r *eventExcludeProfileResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"priority": {
				Type:     types.NumberType,
				Optional: true,
			},
			"resource_type": {
				Type:     types.StringType,
				Computed: true,
			},
			"platform": {
				Type:     types.StringType,
				Optional: true,
			},
			"metadata": {
				Optional: true,
				Type:     types.StringType,
			},
		},
	}, nil
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
		Name:         plan.Name.Value,
		Description:  plan.Description.Value,
		MetadataJSON: plan.Metadata.Value,
		Priority:     plan.Priority,
		ResourceType: plan.ResourceType.Value,
		Platform:     plan.Platform.Value,
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
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJSON)) + "\n"},
		Priority:     eventExcludeProfileResp.Priority,
		ResourceType: types.String{Value: eventExcludeProfileResp.ResourceType},
		Platform:     types.String{Value: eventExcludeProfileResp.Platform},
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
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJSON)) + "\n"},
		Priority:     eventExcludeProfileResp.Priority,
		ResourceType: types.String{Value: eventExcludeProfileResp.ResourceType},
		Platform:     types.String{Value: eventExcludeProfileResp.Platform},
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

	eventExcludeProfileID := state.ID.Value

	// Retrieve values from plan
	var plan EventExcludeProfile
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventExcludeProfileResp, err := r.client.UpdateEventExcludeProfile(uptycs.EventExcludeProfile{
		ID:           eventExcludeProfileID,
		Name:         plan.Name.Value,
		Description:  plan.Description.Value,
		MetadataJSON: plan.Metadata.Value,
		Priority:     plan.Priority,
		ResourceType: plan.ResourceType.Value,
		Platform:     plan.Platform.Value,
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
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJSON)) + "\n"},
		Priority:     eventExcludeProfileResp.Priority,
		ResourceType: types.String{Value: eventExcludeProfileResp.ResourceType},
		Platform:     types.String{Value: eventExcludeProfileResp.Platform},
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

	eventExcludeProfileID := state.ID.Value

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
