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

// Create a new resource
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

	// Map the plan object group names list of strings into a straight []string for the upcoming create
	// We will need to get the IDs of these however so that we can do the create (needs the ID of the object group)
	var objectGroupNames []string
	plan.RoleObjectGroups.ElementsAs(ctx, &objectGroupNames, false)

	// init an empty []uptycs.ObjectGroup of size 0
	roleObjectGroups := make([]uptycs.ObjectGroup, 0)
	//iterate the object_group_names provided
	for _, _rog := range objectGroupNames {
		//Attempt to GET the object group by Name provided in the terraform plan
		rogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{
			Name: _rog,
		})
		// Couldnt get the object group the user provided in object_group_names, error out
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not create role, objectGroup "+_rog+" not found: "+err.Error(),
			)
			return
		}
		// Successful GET, so build up an uptycs.ObjectGroup{} with the objectGroupID being the ID of the objectGroup we validated
		roleObjectGroups = append(roleObjectGroups, uptycs.ObjectGroup{ObjectGroupID: rogResp.ID})
	}

	roleResp, err := r.client.CreateRole(uptycs.Role{
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Permissions:          permissions,
		Custom:               plan.Custom.Value,
		Hidden:               plan.Hidden.Value,
		NoMinimalPermissions: plan.NoMinimalPermissions.Value,
		RoleObjectGroups:     roleObjectGroups, // Now we have an []uptycs.ObjectGroup{} with ObjectGroupID being validated ObjectGroups
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

	// Iterate the response permissions and fill up the list with it
	for _, t := range roleResp.Permissions {
		result.Permissions.Elems = append(result.Permissions.Elems, types.String{Value: t})
	}

	// Iterate the roleObjectGroups in the GET response
	for _, _rogid := range roleResp.RoleObjectGroups {
		//Attempt to GET the object group. Note: the objectGroupID attribute is the ID to GET by
		rogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{ID: _rogid.ObjectGroupID})
		if err != nil {
			// Couldnt find the object group, give an error
			resp.Diagnostics.AddError(
				"Failed to read.",
				"Could not get object group with name  "+_rogid.Name+": "+err.Error(),
			)
			return
		}
		// build up the state object to be the list of strings of objectGroupNames (friendly to the user)
		result.RoleObjectGroups.Elems = append(result.RoleObjectGroups.Elems, types.String{Value: rogResp.Name})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
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

	// Iterate the roleObjectGroups in the GET response
	for _, _rogid := range roleResp.RoleObjectGroups {
		//Attempt to GET the object group. Note: the objectGroupID attribute is the ID to GET by
		rogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{ID: _rogid.ObjectGroupID})
		if err != nil {
			// Couldnt find the object group, give an error
			resp.Diagnostics.AddError(
				"Failed to read.",
				"Could not get object group with ID  "+_rogid.ID+": "+err.Error(),
			)
			return
		}
		// build up the state object to be the list of strings of objectGroupNames (friendly to the user)
		result.RoleObjectGroups.Elems = append(result.RoleObjectGroups.Elems, types.String{Value: rogResp.Name})
	}

	// Iterate the response permissions and fill up the list with it
	for _, t := range roleResp.Permissions {
		result.Permissions.Elems = append(result.Permissions.Elems, types.String{Value: t})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update resource
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

	// Map the plan object group names list of strings into a straight []string for the upcoming update
	// We will need to get the IDs of these however so that we can do the update (needs the ID of the object group)
	var objectGroupNames []string
	plan.RoleObjectGroups.ElementsAs(ctx, &objectGroupNames, false)

	// init an empty []uptycs.ObjectGroup of size 0
	roleObjectGroups := make([]uptycs.ObjectGroup, 0)
	//iterate the object_group_names provided
	for _, _rog := range objectGroupNames {
		//Attempt to GET the object group by Name provided in the terraform plan
		rogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{
			Name: _rog,
		})
		// Couldnt get the object group the user provided in object_group_names, error out
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not update role, objectGroup "+_rog+" not found: "+err.Error(),
			)
			return
		}
		// Successful GET, so build up an uptycs.ObjectGroup{} with the objectGroupID being the ID of the objectGroup we validated
		roleObjectGroups = append(roleObjectGroups, uptycs.ObjectGroup{ObjectGroupID: rogResp.ID})
	}

	roleResp, err := r.client.UpdateRole(uptycs.Role{
		ID:                   roleID,
		Name:                 plan.Name.Value,
		Description:          plan.Description.Value,
		Permissions:          permissions,
		Custom:               plan.Custom.Value,
		Hidden:               plan.Hidden.Value,
		NoMinimalPermissions: plan.NoMinimalPermissions.Value,
		RoleObjectGroups:     roleObjectGroups, // Now we have an []uptycs.ObjectGroup{} with ObjectGroupID being validated ObjectGroups
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

	// Iterate the response permissions and fill up the list with it
	for _, t := range roleResp.Permissions {
		result.Permissions.Elems = append(result.Permissions.Elems, types.String{Value: t})
	}

	// Iterate the roleObjectGroups in the GET response
	for _, _rogid := range roleResp.RoleObjectGroups {
		//Attempt to GET the object group. Note: the objectGroupID attribute is the ID to GET by
		rogResp, err := r.client.GetObjectGroup(uptycs.ObjectGroup{ID: _rogid.ObjectGroupID})
		if err != nil {
			// Couldnt find the object group, give an error
			resp.Diagnostics.AddError(
				"Failed to read.",
				"Could not get object group with name  "+_rogid.Name+": "+err.Error(),
			)
			return
		}
		// build up the state object to be the list of strings of objectGroupNames (friendly to the user)
		result.RoleObjectGroups.Elems = append(result.RoleObjectGroups.Elems, types.String{Value: rogResp.Name})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
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

// Import resource
func (r roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
