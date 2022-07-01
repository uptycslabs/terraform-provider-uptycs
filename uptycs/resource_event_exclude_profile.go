package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type resourceEventExcludeProfileType struct{}

// Resource schema
func (r resourceEventExcludeProfileType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

// New resource instance
func (r resourceEventExcludeProfileType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceEventExcludeProfile{
		p: *(p.(*provider)),
	}, nil
}

type resourceEventExcludeProfile struct {
	p provider
}

// Create a new resource
func (r resourceEventExcludeProfile) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan EventExcludeProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventExcludeProfileResp, err := r.p.client.CreateEventExcludeProfile(uptycs.EventExcludeProfile{
		Name:         plan.Name.Value,
		Description:  plan.Description.Value,
		MetadataJson: plan.Metadata.Value,
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

	metadataJson, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJson)) + "\n"},
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

// Read resource information
func (r resourceEventExcludeProfile) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var eventExcludeProfileId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &eventExcludeProfileId)...)
	eventExcludeProfileResp, err := r.p.client.GetEventExcludeProfile(uptycs.EventExcludeProfile{
		ID: eventExcludeProfileId,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get eventExcludeProfile with ID  "+eventExcludeProfileId+": "+err.Error(),
		)
		return
	}

	metadataJson, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJson)) + "\n"},
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

// Update resource
func (r resourceEventExcludeProfile) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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

	eventExcludeProfileResp, err := r.p.client.UpdateEventExcludeProfile(uptycs.EventExcludeProfile{
		ID:           eventExcludeProfileID,
		Name:         plan.Name.Value,
		Description:  plan.Description.Value,
		MetadataJson: plan.Metadata.Value,
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

	metadataJson, err := json.MarshalIndent(eventExcludeProfileResp.Metadata, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventExcludeProfile{
		ID:           types.String{Value: eventExcludeProfileResp.ID},
		Name:         types.String{Value: eventExcludeProfileResp.Name},
		Description:  types.String{Value: eventExcludeProfileResp.Description},
		Metadata:     types.String{Value: string([]byte(metadataJson)) + "\n"},
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

// Delete resource
func (r resourceEventExcludeProfile) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state EventExcludeProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventExcludeProfileID := state.ID.Value

	_, err := r.p.client.DeleteEventExcludeProfile(uptycs.EventExcludeProfile{
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
func (r resourceEventExcludeProfile) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
