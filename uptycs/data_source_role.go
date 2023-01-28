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

func (d *roleDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"permissions": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"hidden":                 schema.BoolAttribute{Optional: true},
			"no_minimal_permissions": schema.BoolAttribute{Optional: true},
			"role_object_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
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
		ID:          types.StringValue(roleResp.ID),
		Name:        types.StringValue(roleResp.Name),
		Description: types.StringValue(roleResp.Description),
		Permissions: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Hidden:               types.BoolValue(roleResp.Hidden),
		NoMinimalPermissions: types.BoolValue(roleResp.NoMinimalPermissions),
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
