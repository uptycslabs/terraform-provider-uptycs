package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type resourceYaraGroupRuleType struct{}

// Alert Rule Resource schema
func (r resourceYaraGroupRuleType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"rules": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
		},
	}, nil
}

// New resource instance
func (r resourceYaraGroupRuleType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return resourceYaraGroupRule{
		p: *(p.(*Provider)),
	}, nil
}

type resourceYaraGroupRule struct {
	p Provider
}

// Read resource information
func (r resourceYaraGroupRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var yaraGroupRuleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &yaraGroupRuleID)...)
	yaraGroupRuleResp, err := r.p.client.GetYaraGroupRule(uptycs.YaraGroupRule{
		ID: yaraGroupRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get yaraGroupRule with ID  "+yaraGroupRuleID+": "+err.Error(),
		)
		return
	}
	var result = YaraGroupRule{
		ID:          types.String{Value: yaraGroupRuleResp.ID},
		Name:        types.String{Value: yaraGroupRuleResp.Name},
		Description: types.String{Value: yaraGroupRuleResp.Description},
		Rules:       types.String{Value: yaraGroupRuleResp.Rules},
		Custom:      types.Bool{Value: yaraGroupRuleResp.Custom},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Create a new resource
func (r resourceYaraGroupRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan YaraGroupRule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	yaraGroupRuleResp, err := r.p.client.CreateYaraGroupRule(uptycs.YaraGroupRule{
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
		Rules:       plan.Rules.Value,
		Custom:      plan.Custom.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create yaraGroupRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = YaraGroupRule{
		ID:          types.String{Value: yaraGroupRuleResp.ID},
		Name:        types.String{Value: yaraGroupRuleResp.Name},
		Description: types.String{Value: yaraGroupRuleResp.Description},
		Rules:       types.String{Value: yaraGroupRuleResp.Rules},
		Custom:      types.Bool{Value: yaraGroupRuleResp.Custom},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceYaraGroupRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state YaraGroupRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	yaraGroupRuleID := state.ID.Value

	// Retrieve values from plan
	var plan YaraGroupRule
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	yaraGroupRuleResp, err := r.p.client.UpdateYaraGroupRule(uptycs.YaraGroupRule{
		ID:          yaraGroupRuleID,
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
		Rules:       plan.Rules.Value,
		Custom:      plan.Custom.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create yaraGroupRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = YaraGroupRule{
		ID:          types.String{Value: yaraGroupRuleResp.ID},
		Name:        types.String{Value: yaraGroupRuleResp.Name},
		Description: types.String{Value: yaraGroupRuleResp.Description},
		Rules:       types.String{Value: yaraGroupRuleResp.Rules},
		Custom:      types.Bool{Value: yaraGroupRuleResp.Custom},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceYaraGroupRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state YaraGroupRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	yaraGroupRuleID := state.ID.Value

	_, err := r.p.client.DeleteYaraGroupRule(uptycs.YaraGroupRule{
		ID: yaraGroupRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete yaraGroupRule with ID  "+yaraGroupRuleID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceYaraGroupRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
