package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/myoung34/terraform-plugin-framework-utils/modifiers"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (r *userResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Required: true},
			"email": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.DefaultString(""),
				},
			},
			"phone": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.DefaultString(""),
				},
			},
			"active": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					modifiers.DefaultBool(true),
				},
			},
			"super_admin": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(false),
				},
			},
			"bot": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					modifiers.DefaultBool(false),
				},
			},
			"support": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					modifiers.DefaultBool(false),
				},
			},
			"image_url": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					modifiers.DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 512),
				},
			},
			"max_idle_time_mins": schema.Int64Attribute{Required: true},
			"alert_hidden_columns": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"user_object_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
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
		ID:                 types.StringValue(userResp.ID),
		Name:               types.StringValue(userResp.Name),
		Email:              types.StringValue(userResp.Email),
		Phone:              types.StringValue(userResp.Phone),
		Active:             types.BoolValue(userResp.Active),
		SuperAdmin:         types.BoolValue(userResp.SuperAdmin),
		ImageURL:           types.StringValue(userResp.ImageURL),
		Bot:                types.BoolValue(userResp.Bot),
		Support:            types.BoolValue(userResp.Support),
		MaxIdleTimeMins:    types.Int64Value(int64(userResp.MaxIdleTimeMins)),
		Roles:              makeListStringAttributeFn(userResp.Roles, func(f uptycs.Role) (string, bool) { return f.ID, true }),
		AlertHiddenColumns: makeListStringAttribute(userResp.AlertHiddenColumns),
		UserObjectGroups:   makeListStringAttributeFn(userResp.UserObjectGroups, func(g uptycs.ObjectGroup) (string, bool) { return g.ObjectGroupID, true }),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	// Need to turn the list of IDs into a specific Role object with `ID` as the ID attribute
	var roleIDs []string
	plan.Roles.ElementsAs(ctx, &roleIDs, false)
	roles := make([]uptycs.Role, 0)
	for _, _r := range roleIDs {
		roles = append(roles, uptycs.Role{
			ID: _r,
		})
	}

	// Need to turn the list of IDs into a specific ObjectGroup object with `ObjectGroupID` as the ID attribute
	var objectGroupIDs []string
	plan.UserObjectGroups.ElementsAs(ctx, &objectGroupIDs, false)
	userObjectGroups := make([]uptycs.ObjectGroup, 0)
	for _, _uog := range objectGroupIDs {
		userObjectGroups = append(userObjectGroups, uptycs.ObjectGroup{ObjectGroupID: _uog})
	}

	userResp, err := r.client.CreateUser(uptycs.User{
		Name:               plan.Name.ValueString(),
		Email:              plan.Email.ValueString(),
		Phone:              plan.Phone.ValueString(),
		ImageURL:           plan.ImageURL.ValueString(),
		MaxIdleTimeMins:    int(plan.MaxIdleTimeMins.ValueInt64()),
		Active:             plan.Active.ValueBool(),
		SuperAdmin:         plan.SuperAdmin.ValueBool(),
		Bot:                plan.Bot.ValueBool(),
		Support:            plan.Support.ValueBool(),
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
		ID:                 types.StringValue(userResp.ID),
		Name:               types.StringValue(userResp.Name),
		Email:              types.StringValue(userResp.Email),
		Phone:              types.StringValue(userResp.Phone),
		Active:             types.BoolValue(userResp.Active),
		SuperAdmin:         types.BoolValue(userResp.SuperAdmin),
		Bot:                types.BoolValue(userResp.Bot),
		Support:            types.BoolValue(userResp.Support),
		ImageURL:           types.StringValue(userResp.ImageURL),
		MaxIdleTimeMins:    types.Int64Value(int64(userResp.MaxIdleTimeMins)),
		Roles:              makeListStringAttributeFn(userResp.Roles, func(f uptycs.Role) (string, bool) { return f.ID, true }),
		AlertHiddenColumns: makeListStringAttribute(userResp.AlertHiddenColumns),
		UserObjectGroups:   makeListStringAttributeFn(userResp.UserObjectGroups, func(g uptycs.ObjectGroup) (string, bool) { return g.ObjectGroupID, true }),
	}

	diags = resp.State.Set(ctx, result)
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

	userID := state.ID.ValueString()

	// Retrieve values from plan
	var plan User
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var alertHiddenColumns []string
	plan.AlertHiddenColumns.ElementsAs(ctx, &alertHiddenColumns, false)

	// Need to turn the list of IDs into a specific Role object with `ID` as the ID attribute
	var roleIDs []string
	plan.Roles.ElementsAs(ctx, &roleIDs, false)
	roles := make([]uptycs.Role, 0)
	for _, _r := range roleIDs {
		roles = append(roles, uptycs.Role{
			ID: _r,
		})
	}

	// Need to turn the list of IDs into a specific ObjectGroup object with `ObjectGroupID` as the ID attribute
	var objectGroupIDs []string
	plan.UserObjectGroups.ElementsAs(ctx, &objectGroupIDs, false)
	userObjectGroups := make([]uptycs.ObjectGroup, 0)
	for _, _uog := range objectGroupIDs {
		userObjectGroups = append(userObjectGroups, uptycs.ObjectGroup{ObjectGroupID: _uog})
	}

	userResp, err := r.client.UpdateUser(uptycs.User{
		ID:                 userID,
		Name:               plan.Name.ValueString(),
		Email:              plan.Email.ValueString(),
		Phone:              plan.Phone.ValueString(),
		ImageURL:           plan.ImageURL.ValueString(),
		MaxIdleTimeMins:    int(plan.MaxIdleTimeMins.ValueInt64()),
		AlertHiddenColumns: alertHiddenColumns,
		UserObjectGroups:   userObjectGroups,
		Roles:              roles,
		Active:             plan.Active.ValueBool(),
		SuperAdmin:         plan.SuperAdmin.ValueBool(),
		Bot:                plan.Bot.ValueBool(),
		Support:            plan.Support.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	var result = User{
		ID:                 types.StringValue(userResp.ID),
		Name:               types.StringValue(userResp.Name),
		Email:              types.StringValue(userResp.Email),
		Phone:              types.StringValue(userResp.Phone),
		Active:             types.BoolValue(userResp.Active),
		SuperAdmin:         types.BoolValue(userResp.SuperAdmin),
		Bot:                types.BoolValue(userResp.Bot),
		ImageURL:           types.StringValue(userResp.ImageURL),
		Support:            types.BoolValue(userResp.Support),
		MaxIdleTimeMins:    types.Int64Value(int64(userResp.MaxIdleTimeMins)),
		Roles:              makeListStringAttributeFn(userResp.Roles, func(f uptycs.Role) (string, bool) { return f.ID, true }),
		AlertHiddenColumns: makeListStringAttribute(userResp.AlertHiddenColumns),
		UserObjectGroups:   makeListStringAttributeFn(userResp.UserObjectGroups, func(g uptycs.ObjectGroup) (string, bool) { return g.ObjectGroupID, true }),
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

	userID := state.ID.ValueString()

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
