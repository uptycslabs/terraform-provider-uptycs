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
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

func RoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

type roleDataSource struct {
	client *uptycs.Client
}

func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *roleDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"permissions": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"custom": {
				Type:     types.BoolType,
				Optional: true,
			},
			"hidden": {
				Type:     types.BoolType,
				Optional: true,
			},
			"no_minimal_permissions": {
				Type:     types.BoolType,
				Optional: true,
			},
			"role_object_groups": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
		},
	}, nil
}

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var roleID string
	var roleName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &roleID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &roleName)

	var roleToLookup uptycs.Role

	if len(roleID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		roleToLookup = uptycs.Role{
			Name: roleName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		roleToLookup = uptycs.Role{
			ID: roleID,
		}
	}

	roleResp, err := d.client.GetRole(roleToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
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

	for _, t := range roleResp.Permissions {
		result.Permissions.Elems = append(result.Permissions.Elems, types.String{Value: t})
	}

	for _, _rogid := range roleResp.RoleObjectGroups {
		result.RoleObjectGroups.Elems = append(result.RoleObjectGroups.Elems, types.String{Value: _rogid.ObjectGroupID})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
