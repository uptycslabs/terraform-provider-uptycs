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
	_ resource.Resource                = &exceptionResource{}
	_ resource.ResourceWithConfigure   = &exceptionResource{}
	_ resource.ResourceWithImportState = &exceptionResource{}
)

func ExceptionResource() resource.Resource {
	return &exceptionResource{}
}

type exceptionResource struct {
	client *uptycs.Client
}

func (r *exceptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception"
}

func (r *exceptionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *exceptionResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"exception_type": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("sift")},
			},
			"table_name": {
				Type:     types.StringType,
				Optional: true,
			},
			"is_global": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"disabled": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"close_open_alerts": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"rule": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (r *exceptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var exceptionID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &exceptionID)...)
	exceptionResp, err := r.client.GetException(uptycs.Exception{
		ID: exceptionID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get exception with ID  "+exceptionID+": "+err.Error(),
		)
		return
	}

	ruleJSON, err := json.MarshalIndent(exceptionResp.Rule, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Exception{
		ID:              types.String{Value: exceptionResp.ID},
		Name:            types.String{Value: exceptionResp.Name},
		Description:     types.String{Value: exceptionResp.Description},
		ExceptionType:   types.String{Value: exceptionResp.ExceptionType},
		TableName:       types.String{Value: exceptionResp.TableName},
		IsGlobal:        types.Bool{Value: exceptionResp.IsGlobal},
		Disabled:        types.Bool{Value: exceptionResp.Disabled},
		CloseOpenAlerts: types.Bool{Value: exceptionResp.CloseOpenAlerts},
		Rule:            types.String{Value: string([]byte(ruleJSON)) + "\n"},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *exceptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Exception
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exceptionResp, err := r.client.CreateException(uptycs.Exception{
		Name:            plan.Name.Value,
		Description:     plan.Description.Value,
		ExceptionType:   plan.ExceptionType.Value,
		TableName:       plan.TableName.Value,
		IsGlobal:        plan.IsGlobal.Value,
		Disabled:        plan.Disabled.Value,
		CloseOpenAlerts: plan.CloseOpenAlerts.Value,
		Rule:            uptycs.CustomJSONString(plan.Rule.Value),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create exception, unexpected error: "+err.Error(),
		)
		return
	}

	ruleJSON, err := json.MarshalIndent(exceptionResp.Rule, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Exception{
		ID:              types.String{Value: exceptionResp.ID},
		Name:            types.String{Value: exceptionResp.Name},
		Description:     types.String{Value: exceptionResp.Description},
		ExceptionType:   types.String{Value: exceptionResp.ExceptionType},
		TableName:       types.String{Value: exceptionResp.TableName},
		IsGlobal:        types.Bool{Value: exceptionResp.IsGlobal},
		Disabled:        types.Bool{Value: exceptionResp.Disabled},
		CloseOpenAlerts: types.Bool{Value: exceptionResp.CloseOpenAlerts},
		Rule:            types.String{Value: string([]byte(ruleJSON)) + "\n"},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *exceptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Exception
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exceptionID := state.ID.Value

	// Retrieve values from plan
	var plan Exception
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exceptionResp, err := r.client.UpdateException(uptycs.Exception{
		ID:              exceptionID,
		Name:            plan.Name.Value,
		Description:     plan.Description.Value,
		ExceptionType:   plan.ExceptionType.Value,
		TableName:       plan.TableName.Value,
		IsGlobal:        plan.IsGlobal.Value,
		Disabled:        plan.Disabled.Value,
		CloseOpenAlerts: plan.CloseOpenAlerts.Value,
		Rule:            uptycs.CustomJSONString(plan.Rule.Value),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create exception, unexpected error: "+err.Error(),
		)
		return
	}

	ruleJSON, err := json.MarshalIndent(exceptionResp.Rule, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Exception{
		ID:              types.String{Value: exceptionResp.ID},
		Name:            types.String{Value: exceptionResp.Name},
		Description:     types.String{Value: exceptionResp.Description},
		ExceptionType:   types.String{Value: exceptionResp.ExceptionType},
		TableName:       types.String{Value: exceptionResp.TableName},
		IsGlobal:        types.Bool{Value: exceptionResp.IsGlobal},
		Disabled:        types.Bool{Value: exceptionResp.Disabled},
		CloseOpenAlerts: types.Bool{Value: exceptionResp.CloseOpenAlerts},
		Rule:            types.String{Value: string([]byte(ruleJSON)) + "\n"},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *exceptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Exception
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exceptionID := state.ID.Value

	_, err := r.client.DeleteException(uptycs.Exception{
		ID: exceptionID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete exception with ID  "+exceptionID+": "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *exceptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
