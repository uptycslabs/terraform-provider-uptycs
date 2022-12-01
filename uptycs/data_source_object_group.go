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
	_ datasource.DataSource              = &objectGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &objectGroupDataSource{}
)

func ObjectGroupDataSource() datasource.DataSource {
	return &objectGroupDataSource{}
}

type objectGroupDataSource struct {
	client *uptycs.Client
}

func (d *objectGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_group"
}

func (d *objectGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *objectGroupDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"object_group_id": {
				Type:     types.StringType,
				Optional: true,
			},
			"key": {
				Type:     types.StringType,
				Optional: true,
			},
			"value": {
				Type:     types.StringType,
				Optional: true,
			},
			"asset_group_rule_id": {
				Type:     types.StringType,
				Optional: true,
			},
			"user_id": {
				Type:     types.StringType,
				Optional: true,
			},
			"role_id": {
				Type:     types.StringType,
				Optional: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"secret": {
				Type:     types.StringType,
				Optional: true,
			},
			"object_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"retention_days": {
				Type:     types.NumberType,
				Optional: true,
			},
			"destinations": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
		},
	}, nil
}

func (d *objectGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var objectGroupID string
	var objectGroupName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &objectGroupID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &objectGroupName)

	var objectGroupToLookup uptycs.ObjectGroup

	if len(objectGroupID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		objectGroupToLookup = uptycs.ObjectGroup{
			Name: objectGroupName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		objectGroupToLookup = uptycs.ObjectGroup{
			ID: objectGroupID,
		}
	}

	objectGroupResp, err := d.client.GetObjectGroup(objectGroupToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get objectGroup with ID  "+objectGroupID+": "+err.Error(),
		)
		return
	}

	var result = ObjectGroup{
		ID:               types.String{Value: objectGroupResp.ID},
		Name:             types.String{Value: objectGroupResp.Name},
		Key:              types.String{Value: objectGroupResp.Key},
		Value:            types.String{Value: objectGroupResp.Value},
		AssetGroupRuleID: types.String{Value: objectGroupResp.AssetGroupRuleID},
		ObjectGroupID:    types.String{Value: objectGroupResp.ObjectGroupID},
		UserID:           types.String{Value: objectGroupResp.UserID},
		RoleID:           types.String{Value: objectGroupResp.RoleID},
		Description:      types.String{Value: objectGroupResp.Description},
		Secret:           types.String{Value: objectGroupResp.Secret},
		ObjectType:       types.String{Value: objectGroupResp.ObjectType},
		RetentionDays:    0,
		Destinations: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, _dest := range objectGroupResp.Destinations {
		result.Destinations.Elems = append(result.Destinations.Elems, types.String{Value: _dest.ID})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
