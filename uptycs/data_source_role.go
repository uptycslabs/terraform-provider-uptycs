package uptycs

import (
	"context"
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
		ID:                   types.StringValue(roleResp.ID),
		Name:                 types.StringValue(roleResp.Name),
		Description:          types.StringValue(roleResp.Description),
		Hidden:               types.BoolValue(roleResp.Hidden),
		NoMinimalPermissions: types.BoolValue(roleResp.NoMinimalPermissions),
	}

	var diags diag.Diagnostics
	result.Permissions, diags = types.ListValueFrom(ctx, types.StringType, roleResp.Permissions)
	result.RoleObjectGroups, diags = types.ListValueFrom(
		ctx,
		types.StringType,
		StringsFromFn(func(_rogid ObjectGroup) string { return _rogid.ObjectGroupID.String() }, roleResp.RoleObjectGroups),
	)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
