package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (d *objectGroupDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Optional: true},
			"name":                schema.StringAttribute{Optional: true},
			"object_group_id":     schema.StringAttribute{Optional: true},
			"key":                 schema.StringAttribute{Optional: true},
			"value":               schema.StringAttribute{Optional: true},
			"asset_group_rule_id": schema.StringAttribute{Optional: true},
			"user_id":             schema.StringAttribute{Optional: true},
			"role_id":             schema.StringAttribute{Optional: true},
			"description":         schema.StringAttribute{Optional: true},
			"secret":              schema.StringAttribute{Optional: true},
			"object_type":         schema.StringAttribute{Optional: true},
			"retention_days":      schema.NumberAttribute{Optional: true},
			"destinations": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
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
		ID:               types.StringValue(objectGroupResp.ID),
		Name:             types.StringValue(objectGroupResp.Name),
		Key:              types.StringValue(objectGroupResp.Key),
		Value:            types.StringValue(objectGroupResp.Value),
		AssetGroupRuleID: types.StringValue(objectGroupResp.AssetGroupRuleID),
		ObjectGroupID:    types.StringValue(objectGroupResp.ObjectGroupID),
		UserID:           types.StringValue(objectGroupResp.UserID),
		RoleID:           types.StringValue(objectGroupResp.RoleID),
		Description:      types.StringValue(objectGroupResp.Description),
		Secret:           types.StringValue(objectGroupResp.Secret),
		ObjectType:       types.StringValue(objectGroupResp.ObjectType),
		RetentionDays:    0,
		Destinations:     makeListStringAttributeFn(objectGroupResp.Destinations, func(d uptycs.Destination) (string, bool) { return d.ID, true }),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
