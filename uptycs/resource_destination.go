package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/myoung34/uptycs-client-go/uptycs"
)

type resourceDestinationType struct{}

// Alert Rule Resource schema
func (r resourceDestinationType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"type": {
				Type:     types.StringType,
				Optional: true,
			},
			"address": {
				Type:     types.StringType,
				Optional: true,
			},
			"enabled": {
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(true)},
			},
		},
	}, nil
}

// New resource instance
func (r resourceDestinationType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDestination{
		p: *(p.(*provider)),
	}, nil
}

type resourceDestination struct {
	p provider
}

// Create a new resource
func (r resourceDestination) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Destination
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationResp, err := r.p.client.CreateDestination(uptycs.Destination{
		Name:    plan.Name.Value,
		Type:    plan.Type.Value,
		Address: plan.Address.Value,
		Enabled: plan.Enabled.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Destination{
		ID:      types.String{Value: destinationResp.ID},
		Name:    types.String{Value: destinationResp.Name},
		Type:    types.String{Value: destinationResp.Type},
		Address: types.String{Value: destinationResp.Address},
		Enabled: types.Bool{Value: destinationResp.Enabled},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceDestination) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var destinationId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &destinationId)...)
	destinationResp, err := r.p.client.GetDestination(uptycs.Destination{
		ID: destinationId,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading order",
			"Could not get destination with ID  "+destinationId+": "+err.Error(),
		)
		return
	}
	var result = Destination{
		ID:      types.String{Value: destinationResp.ID},
		Name:    types.String{Value: destinationResp.Name},
		Type:    types.String{Value: destinationResp.Type},
		Address: types.String{Value: destinationResp.Address},
		Enabled: types.Bool{Value: destinationResp.Enabled},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update resource
func (r resourceDestination) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state Destination
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationID := state.ID.Value

	// Retrieve values from plan
	var plan Destination
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationResp, err := r.p.client.UpdateDestination(uptycs.Destination{
		ID:      destinationID,
		Name:    plan.Name.Value,
		Type:    plan.Type.Value,
		Address: plan.Address.Value,
		Enabled: plan.Enabled.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Destination{
		ID:      types.String{Value: destinationResp.ID},
		Name:    types.String{Value: destinationResp.Name},
		Type:    types.String{Value: destinationResp.Type},
		Address: types.String{Value: destinationResp.Address},
		Enabled: types.Bool{Value: destinationResp.Enabled},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceDestination) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Destination
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationID := state.ID.Value

	_, err := r.p.client.DeleteDestination(uptycs.Destination{
		ID: destinationID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting order",
			"Could not delete destination with ID  "+destinationID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceDestination) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
