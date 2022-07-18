package uptycs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceUserType struct {
	p provider
}

func (r dataSourceUserType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Type:     types.BoolType,
				Optional: true,
			},
			"super_admin": {
				Type:     types.BoolType,
				Optional: true,
			},
			"bot": {
				Type:     types.BoolType,
				Computed: true,
				Optional: true,
			},
			"support": {
				Type:     types.BoolType,
				Optional: true,
			},
			//"max_idle_time_mins": {
			//	Type:     types.NumberType,
			//	Optional: true,
			//},
		},
	}, nil
}

func (d dataSourceUserType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceUserType{
		p: *(p.(*provider)),
	}, nil
}

func (d dataSourceUserType) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var userId string
	var userName string

	idAttr := req.Config.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &userId)

	var userToLookup uptycs.User

	if len(userId) == 0 {
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("name"), &userName)...)
		userToLookup = uptycs.User{
			Name: userName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		userToLookup = uptycs.User{
			ID: userId,
		}
	}

	userResp, err := d.p.client.GetUser(userToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get user "+userId+" "+userName+": "+err.Error(),
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
