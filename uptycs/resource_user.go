package uptycs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type resourceUserType struct{}

// Alert Rule Resource schema
func (r resourceUserType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"email": {
				Type:     types.StringType,
				Optional: true,
			},
			"active": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(true)},
			},
			"super_admin": {
				Type:     types.BoolType,
				Computed: true,
			},
			"bot": {
				Type:          types.BoolType,
				Computed:      true,
				Optional:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
			},
			"support": {
				Type:          types.BoolType,
				Computed:      true,
				Optional:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(false)},
			},
			//"max_idle_time_mins": {
			//	Type:     types.NumberType,
			//	Optional: true,
			//},
		},
	}, nil
}

// New resource instance
func (r resourceUserType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceUser{
		p: *(p.(*provider)),
	}, nil
}

type resourceUser struct {
	p provider
}

// Create a new resource
func (r resourceUser) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userResp, err := r.p.client.CreateUser(uptycs.User{
		Name:       plan.Name.Value,
		Email:      plan.Email.Value,
		Active:     plan.Active.Value,
		SuperAdmin: plan.SuperAdmin.Value,
		Bot:        plan.Bot.Value,
		Support:    plan.Support.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	var result = User{
		ID:         types.String{Value: userResp.ID},
		Name:       types.String{Value: userResp.Name},
		Email:      types.String{Value: userResp.Email},
		Active:     types.Bool{Value: userResp.Active},
		SuperAdmin: types.Bool{Value: userResp.SuperAdmin},
		Bot:        types.Bool{Value: userResp.Bot},
		Support:    types.Bool{Value: userResp.Support},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceUser) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var userId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &userId)...)
	userResp, err := r.p.client.GetUser(uptycs.User{
		ID: userId,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get user with ID  "+userId+": "+err.Error(),
		)
		return
	}
	var result = User{
		ID:         types.String{Value: userResp.ID},
		Name:       types.String{Value: userResp.Name},
		Email:      types.String{Value: userResp.Email},
		Active:     types.Bool{Value: userResp.Active},
		SuperAdmin: types.Bool{Value: userResp.SuperAdmin},
		Bot:        types.Bool{Value: userResp.Bot},
		Support:    types.Bool{Value: userResp.Support},
		//MaxIdleTimeMins: userResp.MaxIdleTimeMins,
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update resource
func (r resourceUser) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state User
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := state.ID.Value

	// Retrieve values from plan
	var plan User
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userResp, err := r.p.client.UpdateUser(uptycs.User{
		ID:         userID,
		Name:       plan.Name.Value,
		Email:      plan.Email.Value,
		Active:     plan.Active.Value,
		SuperAdmin: plan.SuperAdmin.Value,
		Bot:        plan.Bot.Value,
		Support:    plan.Support.Value,
		//MaxIdleTimeMins: plan.MaxIdleTimeMins,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	var result = User{
		ID:         types.String{Value: userResp.ID},
		Name:       types.String{Value: userResp.Name},
		Email:      types.String{Value: userResp.Email},
		Active:     types.Bool{Value: userResp.Active},
		SuperAdmin: types.Bool{Value: userResp.SuperAdmin},
		Bot:        types.Bool{Value: userResp.Bot},
		Support:    types.Bool{Value: userResp.Support},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceUser) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state User
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := state.ID.Value

	_, err := r.p.client.DeleteUser(uptycs.User{
		ID: userID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete user with ID  "+userID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceUser) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
