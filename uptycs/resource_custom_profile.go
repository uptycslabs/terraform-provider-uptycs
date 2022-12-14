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
	_ resource.Resource                = &customProfileResource{}
	_ resource.ResourceWithConfigure   = &customProfileResource{}
	_ resource.ResourceWithImportState = &customProfileResource{}
)

func CustomProfileResource() resource.Resource {
	return &customProfileResource{}
}

type customProfileResource struct {
	client *uptycs.Client
}

func (r *customProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_profile"
}

func (r *customProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *customProfileResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"query_schedules": {
				Type:     types.StringType,
				Optional: true,
			},
			"priority": {
				Type:     types.NumberType,
				Optional: true,
			},
			"resource_type": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (r *customProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var customProfileID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &customProfileID)...)
	customProfileResp, err := r.client.GetCustomProfile(uptycs.CustomProfile{
		ID: customProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get customProfile with ID  "+customProfileID+": "+err.Error(),
		)
		return
	}

	queryScheduleJSON, err := json.MarshalIndent(customProfileResp.QuerySchedules, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = CustomProfile{
		ID:             types.StringValue(customProfileResp.ID),
		Name:           types.StringValue(customProfileResp.Name),
		Description:    types.StringValue(customProfileResp.Description),
		QuerySchedules: types.StringValue(string([]byte(queryScheduleJSON)) + "\n"),
		Priority:       customProfileResp.Priority,
		ResourceType:   types.StringValue(customProfileResp.ResourceType),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *customProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CustomProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	customProfileResp, err := r.client.CreateCustomProfile(uptycs.CustomProfile{
		Name:           plan.Name.Value,
		Description:    plan.Description.Value,
		QuerySchedules: uptycs.CustomJSONString(plan.QuerySchedules.Value),
		Priority:       plan.Priority,
		ResourceType:   plan.ResourceType.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create customProfile, unexpected error: "+err.Error(),
		)
		return
	}

	queryScheduleJSON, err := json.MarshalIndent(customProfileResp.QuerySchedules, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = CustomProfile{
		ID:             types.StringValue(customProfileResp.ID),
		Name:           types.StringValue(customProfileResp.Name),
		Description:    types.StringValue(customProfileResp.Description),
		QuerySchedules: types.StringValue(string([]byte(queryScheduleJSON)) + "\n"),
		Priority:       customProfileResp.Priority,
		ResourceType:   types.StringValue(customProfileResp.ResourceType),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state CustomProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	customProfileID := state.ID.Value

	// Retrieve values from plan
	var plan CustomProfile
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	customProfileResp, err := r.client.UpdateCustomProfile(uptycs.CustomProfile{
		ID:             customProfileID,
		Name:           plan.Name.Value,
		Description:    plan.Description.Value,
		QuerySchedules: uptycs.CustomJSONString(plan.QuerySchedules.Value),
		Priority:       plan.Priority,
		ResourceType:   plan.ResourceType.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create customProfile, unexpected error: "+err.Error(),
		)
		return
	}

	queryScheduleJSON, err := json.MarshalIndent(customProfileResp.QuerySchedules, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = CustomProfile{
		ID:             types.StringValue(customProfileResp.ID),
		Name:           types.StringValue(customProfileResp.Name),
		Description:    types.StringValue(customProfileResp.Description),
		QuerySchedules: types.StringValue(string([]byte(queryScheduleJSON)) + "\n"),
		Priority:       customProfileResp.Priority,
		ResourceType:   types.StringValue(customProfileResp.ResourceType),
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	customProfileID := state.ID.Value

	_, err := r.client.DeleteCustomProfile(uptycs.CustomProfile{
		ID: customProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete customProfile with ID  "+customProfileID+": "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *customProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
