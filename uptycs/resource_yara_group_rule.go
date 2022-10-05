package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &yaraGroupRuleResource{}
	_ resource.ResourceWithConfigure   = &yaraGroupRuleResource{}
	_ resource.ResourceWithImportState = &yaraGroupRuleResource{}
)

func YaraGroupRuleResource() resource.Resource {
	return &yaraGroupRuleResource{}
}

type yaraGroupRuleResource struct {
	client *uptycs.Client
}

func (r *yaraGroupRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_yara_group_rule"
}

func (r *yaraGroupRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *yaraGroupRuleResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r *yaraGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var yaraGroupRuleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &yaraGroupRuleID)...)
	yaraGroupRuleResp, err := r.client.GetYaraGroupRule(uptycs.YaraGroupRule{
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

func (r *yaraGroupRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan YaraGroupRule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	yaraGroupRuleResp, err := r.client.CreateYaraGroupRule(uptycs.YaraGroupRule{
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

func (r *yaraGroupRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	yaraGroupRuleResp, err := r.client.UpdateYaraGroupRule(uptycs.YaraGroupRule{
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

func (r *yaraGroupRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state YaraGroupRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	yaraGroupRuleID := state.ID.Value

	_, err := r.client.DeleteYaraGroupRule(uptycs.YaraGroupRule{
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

func (r *yaraGroupRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
