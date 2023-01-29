package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (r *roleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Required: true},
			"description": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringDefault(""),
				},
			},
			"permissions": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"hidden": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolDefault(false),
				},
			},
			"no_minimal_permissions": schema.BoolAttribute{Required: true},
			"role_object_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
		ID:                   types.StringValue(roleResp.ID),
		Name:                 types.StringValue(roleResp.Name),
		Description:          types.StringValue(roleResp.Description),
		Permissions:          makeListStringAttributeFn(roleResp.RoleObjectGroups, func(g uptycs.ObjectGroup) (string, bool) { return g.ObjectGroupID, true }),
		Hidden:               types.BoolValue(roleResp.Hidden),
		NoMinimalPermissions: types.BoolValue(roleResp.NoMinimalPermissions),
		RoleObjectGroups:     makeListStringAttribute(roleResp.Permissions),
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
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Permissions:          permissions,
		Hidden:               plan.Hidden.ValueBool(),
		NoMinimalPermissions: plan.NoMinimalPermissions.ValueBool(),
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
		ID:                   types.StringValue(roleResp.ID),
		Name:                 types.StringValue(roleResp.Name),
		Description:          types.StringValue(roleResp.Description),
		Permissions:          makeListStringAttribute(roleResp.Permissions),
		Hidden:               types.BoolValue(roleResp.Hidden),
		NoMinimalPermissions: types.BoolValue(roleResp.NoMinimalPermissions),
		RoleObjectGroups:     makeListStringAttributeFn(roleResp.RoleObjectGroups, func(g uptycs.ObjectGroup) (string, bool) { return g.ObjectGroupID, true }),
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

	roleID := state.ID.ValueString()

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
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Permissions:          permissions,
		Hidden:               plan.Hidden.ValueBool(),
		NoMinimalPermissions: plan.NoMinimalPermissions.ValueBool(),
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
		ID:          types.StringValue(roleResp.ID),
		Name:        types.StringValue(roleResp.Name),
		Description: types.StringValue(roleResp.Description),
		Permissions: makeListStringAttribute(roleResp.Permissions), Hidden: types.BoolValue(roleResp.Hidden),
		NoMinimalPermissions: types.BoolValue(roleResp.NoMinimalPermissions),
		RoleObjectGroups:     makeListStringAttributeFn(roleResp.RoleObjectGroups, func(g uptycs.ObjectGroup) (string, bool) { return g.ObjectGroupID, true }),
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

	roleID := state.ID.ValueString()

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

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
