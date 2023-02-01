package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/myoung34/terraform-plugin-framework-utils/modifiers"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (r *exceptionResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"exception_type": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.DefaultString("sift"),
				},
			},
			"table_name": schema.StringAttribute{Optional: true},
			"is_global": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					modifiers.DefaultBool(true),
				},
			},
			"disabled": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					modifiers.DefaultBool(true),
				},
			},
			"close_open_alerts": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					modifiers.DefaultBool(true),
				},
			},
			"rule": schema.StringAttribute{Optional: true},
		},
	}
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
		ID:              types.StringValue(exceptionResp.ID),
		Name:            types.StringValue(exceptionResp.Name),
		Description:     types.StringValue(exceptionResp.Description),
		ExceptionType:   types.StringValue(exceptionResp.ExceptionType),
		TableName:       types.StringValue(exceptionResp.TableName),
		IsGlobal:        types.BoolValue(exceptionResp.IsGlobal),
		Disabled:        types.BoolValue(exceptionResp.Disabled),
		CloseOpenAlerts: types.BoolValue(exceptionResp.CloseOpenAlerts),
		Rule:            types.StringValue(string(ruleJSON) + "\n"),
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
		Name:            plan.Name.ValueString(),
		Description:     plan.Description.ValueString(),
		ExceptionType:   plan.ExceptionType.ValueString(),
		TableName:       plan.TableName.ValueString(),
		IsGlobal:        plan.IsGlobal.ValueBool(),
		Disabled:        plan.Disabled.ValueBool(),
		CloseOpenAlerts: plan.CloseOpenAlerts.ValueBool(),
		Rule:            uptycs.CustomJSONString(plan.Rule.ValueString()),
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
		ID:              types.StringValue(exceptionResp.ID),
		Name:            types.StringValue(exceptionResp.Name),
		Description:     types.StringValue(exceptionResp.Description),
		ExceptionType:   types.StringValue(exceptionResp.ExceptionType),
		TableName:       types.StringValue(exceptionResp.TableName),
		IsGlobal:        types.BoolValue(exceptionResp.IsGlobal),
		Disabled:        types.BoolValue(exceptionResp.Disabled),
		CloseOpenAlerts: types.BoolValue(exceptionResp.CloseOpenAlerts),
		Rule:            types.StringValue(string(ruleJSON) + "\n"),
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

	exceptionID := state.ID.ValueString()

	// Retrieve values from plan
	var plan Exception
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	exceptionResp, err := r.client.UpdateException(uptycs.Exception{
		ID:              exceptionID,
		Name:            plan.Name.ValueString(),
		Description:     plan.Description.ValueString(),
		ExceptionType:   plan.ExceptionType.ValueString(),
		TableName:       plan.TableName.ValueString(),
		IsGlobal:        plan.IsGlobal.ValueBool(),
		Disabled:        plan.Disabled.ValueBool(),
		CloseOpenAlerts: plan.CloseOpenAlerts.ValueBool(),
		Rule:            uptycs.CustomJSONString(plan.Rule.ValueString()),
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
		ID:              types.StringValue(exceptionResp.ID),
		Name:            types.StringValue(exceptionResp.Name),
		Description:     types.StringValue(exceptionResp.Description),
		ExceptionType:   types.StringValue(exceptionResp.ExceptionType),
		TableName:       types.StringValue(exceptionResp.TableName),
		IsGlobal:        types.BoolValue(exceptionResp.IsGlobal),
		Disabled:        types.BoolValue(exceptionResp.Disabled),
		CloseOpenAlerts: types.BoolValue(exceptionResp.CloseOpenAlerts),
		Rule:            types.StringValue(string(ruleJSON) + "\n"),
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

	exceptionID := state.ID.ValueString()

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
