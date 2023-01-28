package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func DestinationResource() resource.Resource {
	return &destinationResource{}
}

type destinationResource struct {
	client *uptycs.Client
}

func (r *destinationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (r *destinationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *destinationResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true},
			"name":    schema.StringAttribute{Optional: true},
			"type":    schema.StringAttribute{Optional: true},
			"address": schema.StringAttribute{Optional: true},
			"enabled": schema.BoolAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(true)},
			},
		},
	}
}

func (r *destinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var destinationID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &destinationID)...)
	destinationResp, err := r.client.GetDestination(uptycs.Destination{
		ID: destinationID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get destination with ID  "+destinationID+": "+err.Error(),
		)
		return
	}
	var result = Destination{
		ID:      types.StringValue(destinationResp.ID),
		Name:    types.StringValue(destinationResp.Name),
		Type:    types.StringValue(destinationResp.Type),
		Address: types.StringValue(destinationResp.Address),
		Enabled: types.BoolValue(destinationResp.Enabled),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *destinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Destination
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationResp, err := r.client.CreateDestination(uptycs.Destination{
		Name:    plan.Name.Value,
		Type:    plan.Type.Value,
		Address: plan.Address.Value,
		Enabled: plan.Enabled.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Destination{
		ID:      types.StringValue(destinationResp.ID),
		Name:    types.StringValue(destinationResp.Name),
		Type:    types.StringValue(destinationResp.Type),
		Address: types.StringValue(destinationResp.Address),
		Enabled: types.BoolValue(destinationResp.Enabled),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	destinationResp, err := r.client.UpdateDestination(uptycs.Destination{
		ID:      destinationID,
		Name:    plan.Name.Value,
		Type:    plan.Type.Value,
		Address: plan.Address.Value,
		Enabled: plan.Enabled.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Destination{
		ID:      types.StringValue(destinationResp.ID),
		Name:    types.StringValue(destinationResp.Name),
		Type:    types.StringValue(destinationResp.Type),
		Address: types.StringValue(destinationResp.Address),
		Enabled: types.BoolValue(destinationResp.Enabled),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *destinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Destination
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationID := state.ID.Value

	_, err := r.client.DeleteDestination(uptycs.Destination{
		ID: destinationID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete destination with ID  "+destinationID+": "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *destinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
