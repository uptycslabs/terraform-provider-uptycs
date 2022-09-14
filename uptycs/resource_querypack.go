package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type resourceQuerypackType struct{}

// Alert Rule Resource schema
func (r resourceQuerypackType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Required: true,
			},
			"type": {
				Type:        types.StringType,
				Description: "Should be one of: compliance default hardware incident system vulnerability",
				Required:    true,
			},
			"additional_logger": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"is_internal": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"resource_type": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("asset")},
			},
			"conf": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

// New resource instance
func (r resourceQuerypackType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return resourceQuerypack{
		p: *(p.(*Provider)),
	}, nil
}

type resourceQuerypack struct {
	p Provider
}

// Read resource information
func (r resourceQuerypack) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var queryPackID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &queryPackID)...)
	querypackResp, err := r.p.client.GetQuerypack(uptycs.Querypack{
		ID: queryPackID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get queryPack with ID  "+queryPackID+": "+err.Error(),
		)
		return
	}

	queryPackConfJSON, err := json.MarshalIndent(querypackResp.Conf, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Querypack{
		ID:               types.String{Value: querypackResp.ID},
		Name:             types.String{Value: querypackResp.Name},
		Description:      types.String{Value: querypackResp.Description},
		Type:             types.String{Value: querypackResp.Type},
		AdditionalLogger: types.Bool{Value: querypackResp.AdditionalLogger},
		Custom:           types.Bool{Value: querypackResp.Custom},
		IsInternal:       types.Bool{Value: querypackResp.IsInternal},
		ResourceType:     types.String{Value: querypackResp.ResourceType},
		Conf:             types.String{Value: string([]byte(queryPackConfJSON)) + "\n"},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Create a new resource
func (r resourceQuerypack) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Querypack
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//r.p.client.HostURL = "http://localhost:8000"
	querypackResp, err := r.p.client.CreateQuerypack(uptycs.Querypack{
		Name:             plan.Name.Value,
		Description:      plan.Description.Value,
		Type:             plan.Type.Value,
		AdditionalLogger: plan.AdditionalLogger.Value,
		Custom:           plan.Custom.Value,
		IsInternal:       plan.IsInternal.Value,
		ResourceType:     plan.ResourceType.Value,
		Queries:          []uptycs.Query{},
		Conf:             uptycs.CustomJSONString(plan.Conf.Value),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create queryPack, unexpected error: "+err.Error(),
		)
		return
	}

	queryPackConfJSON, err := json.MarshalIndent(querypackResp.Conf, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Querypack{
		ID:               types.String{Value: querypackResp.ID},
		Name:             types.String{Value: querypackResp.Name},
		Description:      types.String{Value: querypackResp.Description},
		Type:             types.String{Value: querypackResp.Type},
		AdditionalLogger: types.Bool{Value: querypackResp.AdditionalLogger},
		Custom:           types.Bool{Value: querypackResp.Custom},
		IsInternal:       types.Bool{Value: querypackResp.IsInternal},
		ResourceType:     types.String{Value: querypackResp.ResourceType},
		Conf:             types.String{Value: string([]byte(queryPackConfJSON)) + "\n"},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceQuerypack) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Querypack
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryPackID := state.ID.Value

	// Retrieve values from plan
	var plan Querypack
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	querypackResp, err := r.p.client.UpdateQuerypack(uptycs.Querypack{
		ID: queryPackID,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create queryPack, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Querypack{
		ID: types.String{Value: querypackResp.ID},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceQuerypack) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Querypack
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryPackID := state.ID.Value

	_, err := r.p.client.DeleteQuerypack(uptycs.Querypack{
		ID: queryPackID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete queryPack with ID  "+queryPackID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceQuerypack) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
