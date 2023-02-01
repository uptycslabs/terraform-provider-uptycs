package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/myoung34/terraform-plugin-framework-utils/modifiers"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (r *yaraGroupRuleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"rules": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.DefaultString(""),
				},
			},
		},
	}
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
		ID:          types.StringValue(yaraGroupRuleResp.ID),
		Name:        types.StringValue(yaraGroupRuleResp.Name),
		Description: types.StringValue(yaraGroupRuleResp.Description),
		Rules:       types.StringValue(yaraGroupRuleResp.Rules),
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
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Rules:       plan.Rules.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create yaraGroupRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = YaraGroupRule{
		ID:          types.StringValue(yaraGroupRuleResp.ID),
		Name:        types.StringValue(yaraGroupRuleResp.Name),
		Description: types.StringValue(yaraGroupRuleResp.Description),
		Rules:       types.StringValue(yaraGroupRuleResp.Rules),
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

	yaraGroupRuleID := state.ID.ValueString()

	// Retrieve values from plan
	var plan YaraGroupRule
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	yaraGroupRuleResp, err := r.client.UpdateYaraGroupRule(uptycs.YaraGroupRule{
		ID:          yaraGroupRuleID,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Rules:       plan.Rules.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create yaraGroupRule, unexpected error: "+err.Error(),
		)
		return
	}

	var result = YaraGroupRule{
		ID:          types.StringValue(yaraGroupRuleResp.ID),
		Name:        types.StringValue(yaraGroupRuleResp.Name),
		Description: types.StringValue(yaraGroupRuleResp.Description),
		Rules:       types.StringValue(yaraGroupRuleResp.Rules),
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

	yaraGroupRuleID := state.ID.ValueString()

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
