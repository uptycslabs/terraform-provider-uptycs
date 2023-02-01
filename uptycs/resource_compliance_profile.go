package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (r *complianceProfileResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"priority":    schema.NumberAttribute{Optional: true},
		},
	}
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
		ID:          types.StringValue(complianceProfileResp.ID),
		Name:        types.StringValue(complianceProfileResp.Name),
		Description: types.StringValue(complianceProfileResp.Description),
		Priority:    complianceProfileResp.Priority,
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
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Priority:    plan.Priority,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create complianceProfile, unexpected error: "+err.Error(),
		)
		return
	}

	var result = ComplianceProfile{
		ID:          types.StringValue(complianceProfileResp.ID),
		Name:        types.StringValue(complianceProfileResp.Name),
		Description: types.StringValue(complianceProfileResp.Description),
		Priority:    complianceProfileResp.Priority,
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

	complianceProfileID := state.ID.ValueString()

	// Retrieve values from plan
	var plan ComplianceProfile
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	complianceProfileResp, err := r.client.UpdateComplianceProfile(uptycs.ComplianceProfile{
		ID:          complianceProfileID,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Priority:    plan.Priority,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create complianceProfile, unexpected error: "+err.Error(),
		)
		return
	}

	var result = ComplianceProfile{
		ID:          types.StringValue(complianceProfileResp.ID),
		Name:        types.StringValue(complianceProfileResp.Name),
		Description: types.StringValue(complianceProfileResp.Description),
		Priority:    complianceProfileResp.Priority,
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

	complianceProfileID := state.ID.ValueString()

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
