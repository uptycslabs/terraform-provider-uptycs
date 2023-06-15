package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func LookupTableDataSource() datasource.DataSource {
	return &lookupTableDataSource{}
}

type lookupTableDataSource struct {
	client *uptycs.Client
}

func (d *lookupTableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lookup_table"
}

func (d *lookupTableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *lookupTableDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"id_field":    schema.StringAttribute{Optional: true},
			"data_rows": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (d *lookupTableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var lookupTableID string
	var lookupTableName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &lookupTableID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &lookupTableName)

	var lookupTableToLookup uptycs.LookupTable

	if len(lookupTableID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		lookupTableToLookup = uptycs.LookupTable{
			Name: lookupTableName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		lookupTableToLookup = uptycs.LookupTable{
			ID: lookupTableID,
		}
	}

	lookupTableResp, err := d.client.GetLookupTable(lookupTableToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get lookupTable with ID  "+lookupTableID+": "+err.Error(),
		)
		return
	}

	var result = LookupTable{
		ID:          types.StringValue(lookupTableResp.ID),
		Name:        types.StringValue(lookupTableResp.Name),
		Description: types.StringValue(lookupTableResp.Description),
		IDField:     types.StringValue(lookupTableResp.IDField),
		DataRows:    makeListStringAttributeFn(lookupTableResp.DataRows, func(v uptycs.LookupTableDataRow) (string, bool) { return string(v.Data), true }),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
