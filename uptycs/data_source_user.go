package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func UserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *uptycs.Client
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *userDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                 schema.StringAttribute{Optional: true},
			"name":               schema.StringAttribute{Optional: true},
			"email":              schema.StringAttribute{Optional: true},
			"phone":              schema.StringAttribute{Optional: true},
			"active":             schema.BoolAttribute{Optional: true},
			"super_admin":        schema.BoolAttribute{Optional: true},
			"bot":                schema.BoolAttribute{Optional: true},
			"support":            schema.BoolAttribute{Optional: true},
			"image_url":          schema.StringAttribute{Optional: true},
			"max_idle_time_mins": schema.NumberAttribute{Optional: true},
			"alert_hidden_columns": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"user_object_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var userID string
	var userName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &userID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &userName)

	var userToLookup uptycs.User

	if len(userID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		userToLookup = uptycs.User{
			Name: userName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		userToLookup = uptycs.User{
			ID: userID,
		}
	}

	userResp, err := d.client.GetUser(userToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get user with ID  "+userID+": "+err.Error(),
		)
		return
	}

	var result = User{
		ID:              types.StringValue(userResp.ID),
		Name:            types.StringValue(userResp.Name),
		Email:           types.StringValue(userResp.Email),
		Phone:           types.StringValue(userResp.Phone),
		Active:          types.BoolValue(userResp.Active),
		SuperAdmin:      types.BoolValue(userResp.SuperAdmin),
		Bot:             types.BoolValue(userResp.Bot),
		Support:         types.BoolValue(userResp.Support),
		ImageURL:        types.StringValue(userResp.ImageURL),
		MaxIdleTimeMins: userResp.MaxIdleTimeMins,
		AlertHiddenColumns: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Roles: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		UserObjectGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, _r := range userResp.Roles {
		result.Roles.Elems = append(result.Roles.Elems, types.String{Value: _r.ID})
	}

	for _, _ahc := range userResp.AlertHiddenColumns {
		result.AlertHiddenColumns.Elems = append(result.AlertHiddenColumns.Elems, types.String{Value: _ahc})
	}

	for _, _uogid := range userResp.UserObjectGroups {
		result.UserObjectGroups.Elems = append(result.UserObjectGroups.Elems, types.String{Value: _uogid.ObjectGroupID})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
