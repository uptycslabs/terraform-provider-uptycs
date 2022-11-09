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
	_ resource.Resource                = &complianceProfileResource{}
	_ resource.ResourceWithConfigure   = &complianceProfileResource{}
	_ resource.ResourceWithImportState = &complianceProfileResource{}
)

func ComplianceProfileResource() resource.Resource {
	return &complianceProfileResource{}
}

type complianceProfileResource struct {
	client *uptycs.Client
}

func (r *complianceProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_profile"
}

func (r *complianceProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *complianceProfileResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(true)},
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"priority": {
				Type:     types.NumberType,
				Optional: true,
			},
		},
	}, nil
}

func (r *complianceProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var complianceProfileID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &complianceProfileID)...)
	complianceProfileResp, err := r.client.GetComplianceProfile(uptycs.ComplianceProfile{
		ID: complianceProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get complianceProfile with ID  "+complianceProfileID+": "+err.Error(),
		)
		return
	}

	var result = ComplianceProfile{
		ID:          types.String{Value: complianceProfileResp.ID},
		Name:        types.String{Value: complianceProfileResp.Name},
		Description: types.String{Value: complianceProfileResp.Description},
		Priority:    complianceProfileResp.Priority,
		Custom:      types.Bool{Value: complianceProfileResp.Custom},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *complianceProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ComplianceProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	complianceProfileResp, err := r.client.CreateComplianceProfile(uptycs.ComplianceProfile{
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
		Priority:    plan.Priority,
		Custom:      plan.Custom.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create complianceProfile, unexpected error: "+err.Error(),
		)
		return
	}

	var result = ComplianceProfile{
		ID:          types.String{Value: complianceProfileResp.ID},
		Name:        types.String{Value: complianceProfileResp.Name},
		Description: types.String{Value: complianceProfileResp.Description},
		Priority:    complianceProfileResp.Priority,
		Custom:      types.Bool{Value: complianceProfileResp.Custom},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *complianceProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state ComplianceProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	complianceProfileID := state.ID.Value

	// Retrieve values from plan
	var plan ComplianceProfile
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	complianceProfileResp, err := r.client.UpdateComplianceProfile(uptycs.ComplianceProfile{
		ID:          complianceProfileID,
		Name:        plan.Name.Value,
		Description: plan.Description.Value,
		Priority:    plan.Priority,
		Custom:      plan.Custom.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create complianceProfile, unexpected error: "+err.Error(),
		)
		return
	}

	var result = ComplianceProfile{
		ID:          types.String{Value: complianceProfileResp.ID},
		Name:        types.String{Value: complianceProfileResp.Name},
		Description: types.String{Value: complianceProfileResp.Description},
		Priority:    complianceProfileResp.Priority,
		Custom:      types.Bool{Value: complianceProfileResp.Custom},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *complianceProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ComplianceProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	complianceProfileID := state.ID.Value

	_, err := r.client.DeleteComplianceProfile(uptycs.ComplianceProfile{
		ID: complianceProfileID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete complianceProfile with ID  "+complianceProfileID+": "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *complianceProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
