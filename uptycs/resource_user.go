package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func UserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *uptycs.Client
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *userResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"email": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"phone": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"active": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"super_admin": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"bot": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"support": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"image_url": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthBetween(1, 512),
				},
			},
			"max_idle_time_mins": {
				Type:     types.NumberType,
				Required: true,
			},
			"alert_hidden_columns": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"roles": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"user_object_groups": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
		},
	}, nil
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var alertHiddenColumns []string
	plan.AlertHiddenColumns.ElementsAs(ctx, &alertHiddenColumns, false)

	var roleNames []string
	plan.Roles.ElementsAs(ctx, &roleNames, false)

	roles := make([]uptycs.Role, 0)
	for _, _r := range roleNames {
		roleResp, err := r.client.GetRole(uptycs.Role{
			Name: _r,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not create user, user "+_r+" not found: "+err.Error(),
			)
			return
		}
		roles = append(roles, roleResp)
	}

	var objectGroupNames []string
	plan.UserObjectGroups.ElementsAs(ctx, &objectGroupNames, false)

	userObjectGroups := make([]uptycs.ObjectGroup, 0)
	for _, _uog := range objectGroupNames {
		uogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{
			Name: _uog,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not create user, objectGroup "+_uog+" not found: "+err.Error(),
			)
			return
		}
		userObjectGroups = append(userObjectGroups, uptycs.ObjectGroup{ObjectGroupID: uogResp.ID})
	}

	userResp, err := r.client.CreateUser(uptycs.User{
		Name:               plan.Name.Value,
		Email:              plan.Email.Value,
		Phone:              plan.Phone.Value,
		ImageURL:           plan.ImageURL.Value,
		MaxIdleTimeMins:    plan.MaxIdleTimeMins,
		Active:             plan.Active.Value,
		SuperAdmin:         plan.SuperAdmin.Value,
		Bot:                plan.Bot.Value,
		Support:            plan.Support.Value,
		AlertHiddenColumns: alertHiddenColumns,
		Roles:              roles,
		UserObjectGroups:   userObjectGroups,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	var result = User{
		ID:              types.String{Value: userResp.ID},
		Name:            types.String{Value: userResp.Name},
		Email:           types.String{Value: userResp.Email},
		Phone:           types.String{Value: userResp.Phone},
		Active:          types.Bool{Value: userResp.Active},
		SuperAdmin:      types.Bool{Value: userResp.SuperAdmin},
		Bot:             types.Bool{Value: userResp.Bot},
		Support:         types.Bool{Value: userResp.Support},
		ImageURL:        types.String{Value: userResp.ImageURL},
		MaxIdleTimeMins: userResp.MaxIdleTimeMins,
		Roles: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		AlertHiddenColumns: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		UserObjectGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, _r := range userResp.Roles {
		result.Roles.Elems = append(result.Roles.Elems, types.String{Value: _r.Name})
	}

	for _, _ahc := range userResp.AlertHiddenColumns {
		result.AlertHiddenColumns.Elems = append(result.AlertHiddenColumns.Elems, types.String{Value: _ahc})
	}

	for _, _uogid := range userResp.UserObjectGroups {
		uogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{ID: _uogid.ObjectGroupID})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read.",
				"Could not get object group with name  "+_uogid.Name+": "+err.Error(),
			)
			return
		}
		result.UserObjectGroups.Elems = append(result.UserObjectGroups.Elems, types.String{Value: uogResp.Name})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var userID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &userID)...)
	userResp, err := r.client.GetUser(uptycs.User{
		ID: userID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get user with ID  "+userID+": "+err.Error(),
		)
		return
	}
	var result = User{
		ID:              types.String{Value: userResp.ID},
		Name:            types.String{Value: userResp.Name},
		Email:           types.String{Value: userResp.Email},
		Phone:           types.String{Value: userResp.Phone},
		Active:          types.Bool{Value: userResp.Active},
		SuperAdmin:      types.Bool{Value: userResp.SuperAdmin},
		ImageURL:        types.String{Value: userResp.ImageURL},
		Bot:             types.Bool{Value: userResp.Bot},
		Support:         types.Bool{Value: userResp.Support},
		MaxIdleTimeMins: userResp.MaxIdleTimeMins,
		Roles: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		AlertHiddenColumns: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		UserObjectGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, _r := range userResp.Roles {
		result.Roles.Elems = append(result.Roles.Elems, types.String{Value: _r.Name})
	}

	for _, _ahc := range userResp.AlertHiddenColumns {
		result.AlertHiddenColumns.Elems = append(result.AlertHiddenColumns.Elems, types.String{Value: _ahc})
	}

	for _, _uogid := range userResp.UserObjectGroups {
		uogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{ID: _uogid.ObjectGroupID})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read.",
				"Could not get object group with ID  "+_uogid.ID+": "+err.Error(),
			)
			return
		}
		result.UserObjectGroups.Elems = append(result.UserObjectGroups.Elems, types.String{Value: uogResp.Name})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	var alertHiddenColumns []string
	plan.AlertHiddenColumns.ElementsAs(ctx, &alertHiddenColumns, false)

	var roleNames []string
	plan.Roles.ElementsAs(ctx, &roleNames, false)

	roles := make([]uptycs.Role, 0)
	for _, _r := range roleNames {
		roleResp, err := r.client.GetRole(uptycs.Role{
			Name: _r,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not create user, user "+_r+" not found: "+err.Error(),
			)
			return
		}
		roles = append(roles, roleResp)
	}

	var objectGroupNames []string
	plan.UserObjectGroups.ElementsAs(ctx, &objectGroupNames, false)

	userObjectGroups := make([]uptycs.ObjectGroup, 0)
	for _, _uog := range objectGroupNames {
		uogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{
			Name: _uog,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not create user, objectGroup "+_uog+" not found: "+err.Error(),
			)
			return
		}
		userObjectGroups = append(userObjectGroups, uptycs.ObjectGroup{ObjectGroupID: uogResp.ID})
	}

	userResp, err := r.client.UpdateUser(uptycs.User{
		ID:                 userID,
		Name:               plan.Name.Value,
		Email:              plan.Email.Value,
		Phone:              plan.Phone.Value,
		ImageURL:           plan.ImageURL.Value,
		MaxIdleTimeMins:    plan.MaxIdleTimeMins,
		AlertHiddenColumns: alertHiddenColumns,
		UserObjectGroups:   userObjectGroups,
		Roles:              roles,
		Active:             plan.Active.Value,
		SuperAdmin:         plan.SuperAdmin.Value,
		Bot:                plan.Bot.Value,
		Support:            plan.Support.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	var result = User{
		ID:              types.String{Value: userResp.ID},
		Name:            types.String{Value: userResp.Name},
		Email:           types.String{Value: userResp.Email},
		Phone:           types.String{Value: userResp.Phone},
		Active:          types.Bool{Value: userResp.Active},
		SuperAdmin:      types.Bool{Value: userResp.SuperAdmin},
		Bot:             types.Bool{Value: userResp.Bot},
		ImageURL:        types.String{Value: userResp.ImageURL},
		Support:         types.Bool{Value: userResp.Support},
		MaxIdleTimeMins: userResp.MaxIdleTimeMins,
		Roles: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		AlertHiddenColumns: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		UserObjectGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}
	for _, _uogid := range userResp.UserObjectGroups {
		uogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{ID: _uogid.ObjectGroupID})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read.",
				"Could not get object group with ID  "+_uogid.ObjectGroupID+": "+err.Error(),
			)
			return
		}
		result.UserObjectGroups.Elems = append(result.UserObjectGroups.Elems, types.String{Value: uogResp.Name})
	}

	for _, _r := range userResp.Roles {
		result.Roles.Elems = append(result.Roles.Elems, types.String{Value: _r.Name})
	}

	for _, _ahc := range userResp.AlertHiddenColumns {
		result.AlertHiddenColumns.Elems = append(result.AlertHiddenColumns.Elems, types.String{Value: _ahc})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state User
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := state.ID.Value

	_, err := r.client.DeleteUser(uptycs.User{
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

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
