package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
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

func (d *userDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
			},
			"email": {
				Type:     types.StringType,
				Optional: true,
			},
			"phone": {
				Type:     types.StringType,
				Optional: true,
			},
			"active": {
				Type:     types.BoolType,
				Optional: true,
			},
			"super_admin": {
				Type:     types.BoolType,
				Optional: true,
			},
			"bot": {
				Type:     types.BoolType,
				Optional: true,
			},
			"support": {
				Type:     types.BoolType,
				Optional: true,
			},
			"image_url": {
				Type:     types.StringType,
				Optional: true,
			},
			"max_idle_time_mins": {
				Type:     types.NumberType,
				Optional: true,
			},
			"alert_hidden_columns": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"roles": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"user_object_groups": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
		},
	}, nil
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
