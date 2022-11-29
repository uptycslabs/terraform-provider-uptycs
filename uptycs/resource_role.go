package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

func RoleResource() resource.Resource {
	return &roleResource{}
}

type roleResource struct {
	client *uptycs.Client
}

func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *roleResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"description": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"permissions": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"hidden": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"no_minimal_permissions": {
				Type:     types.BoolType,
				Required: true,
			},
			"role_object_groups": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
		},
	}, nil
}

func (r roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var roleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &roleID)...)
	roleResp, err := r.client.GetRole(uptycs.Role{
		ID: roleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get role with ID  "+roleID+": "+err.Error(),
		)
		return
	}

	var result = Role{
		ID:          types.String{Value: roleResp.ID},
		Name:        types.String{Value: roleResp.Name},
		Description: types.String{Value: roleResp.Description},
		Permissions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Custom:               types.Bool{Value: roleResp.Custom},
		Hidden:               types.Bool{Value: roleResp.Hidden},
		NoMinimalPermissions: types.Bool{Value: roleResp.NoMinimalPermissions},
		RoleObjectGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, _rogid := range roleResp.RoleObjectGroups {
		result.RoleObjectGroups.Elems = append(result.RoleObjectGroups.Elems, types.String{Value: _rogid.ObjectGroupID})
	}

	for _, t := range roleResp.Permissions {
		result.Permissions.Elems = append(result.Permissions.Elems, types.String{Value: t})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Role
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map the plan permissions list of strings into a straight []string for the upcoming create
	var permissions []string
	plan.Permissions.ElementsAs(ctx, &permissions, false)

	// Need to turn the list of IDs into a specific ObjectGroup object with `ObjectGroupID` as the ID attribute
	var objectGroupIDs []string
	plan.RoleObjectGroups.ElementsAs(ctx, &objectGroupIDs, false)
	roleObjectGroups := make([]uptycs.ObjectGroup, 0)
	for _, _rog := range objectGroupIDs {
		roleObjectGroups = append(roleObjectGroups, uptycs.ObjectGroup{ObjectGroupID: _rog})
	}

	roleResp, err := r.client.CreateRole(uptycs.Role{
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Permissions:          permissions,
		Custom:               plan.Custom.Value,
		Hidden:               plan.Hidden.Value,
		NoMinimalPermissions: plan.NoMinimalPermissions.Value,
		RoleObjectGroups:     roleObjectGroups,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create role, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Role{
		ID:          types.String{Value: roleResp.ID},
		Name:        types.String{Value: roleResp.Name},
		Description: types.String{Value: roleResp.Description},
		Permissions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Custom:               types.Bool{Value: roleResp.Custom},
		Hidden:               types.Bool{Value: roleResp.Hidden},
		NoMinimalPermissions: types.Bool{Value: roleResp.NoMinimalPermissions},
		RoleObjectGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, t := range roleResp.Permissions {
		result.Permissions.Elems = append(result.Permissions.Elems, types.String{Value: t})
	}

	for _, _rogid := range roleResp.RoleObjectGroups {
		result.RoleObjectGroups.Elems = append(result.RoleObjectGroups.Elems, types.String{Value: _rogid.ObjectGroupID})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Role
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := state.ID.Value

	// Retrieve values from plan
	var plan Role
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map the plan permissions list of strings into a straight []string for the upcoming update
	var permissions []string
	plan.Permissions.ElementsAs(ctx, &permissions, false)

	// Need to turn the list of IDs into a specific ObjectGroup object with `ObjectGroupID` as the ID attribute
	var objectGroupIDs []string
	plan.RoleObjectGroups.ElementsAs(ctx, &objectGroupIDs, false)
	roleObjectGroups := make([]uptycs.ObjectGroup, 0)
	for _, _rog := range objectGroupIDs {
		roleObjectGroups = append(roleObjectGroups, uptycs.ObjectGroup{ObjectGroupID: _rog})
	}

	roleResp, err := r.client.UpdateRole(uptycs.Role{
		ID:                   roleID,
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Permissions:          permissions,
		Custom:               plan.Custom.Value,
		Hidden:               plan.Hidden.Value,
		NoMinimalPermissions: plan.NoMinimalPermissions.Value,
		RoleObjectGroups:     roleObjectGroups,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not update role, unexpected error: "+err.Error(),
		)
		return
	}
	var result = Role{
		ID:          types.String{Value: roleResp.ID},
		Name:        types.String{Value: roleResp.Name},
		Description: types.String{Value: roleResp.Description},
		Permissions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Custom:               types.Bool{Value: roleResp.Custom},
		Hidden:               types.Bool{Value: roleResp.Hidden},
		NoMinimalPermissions: types.Bool{Value: roleResp.NoMinimalPermissions},
		RoleObjectGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, t := range roleResp.Permissions {
		result.Permissions.Elems = append(result.Permissions.Elems, types.String{Value: t})
	}

	for _, _rogid := range roleResp.RoleObjectGroups {
		result.RoleObjectGroups.Elems = append(result.RoleObjectGroups.Elems, types.String{Value: _rogid.ObjectGroupID})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Role
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := state.ID.Value

	_, err := r.client.DeleteRole(uptycs.Role{
		ID: roleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete role with ID  "+roleID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
