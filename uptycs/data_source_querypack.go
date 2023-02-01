package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func QuerypackDataSource() datasource.DataSource {
	return &querypackDataSource{}
}

type querypackDataSource struct {
	client *uptycs.Client
}

func (d *querypackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_querypack"
}

func (d *querypackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *querypackDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                schema.StringAttribute{Optional: true},
			"name":              schema.StringAttribute{Optional: true},
			"description":       schema.StringAttribute{Optional: true},
			"type":              schema.StringAttribute{Optional: true},
			"additional_logger": schema.BoolAttribute{Optional: true},
			"is_internal":       schema.BoolAttribute{Optional: true},
			"resource_type":     schema.StringAttribute{Optional: true},
			"conf":              schema.StringAttribute{Optional: true},
		},
	}
}

func (d *querypackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var querypackID string
	var querypackName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &querypackID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &querypackName)

	var queryPackToLookup uptycs.Querypack

	if len(querypackID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		queryPackToLookup = uptycs.Querypack{
			Name: querypackName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		queryPackToLookup = uptycs.Querypack{
			ID: querypackID,
		}
	}

	querypackResp, err := d.client.GetQuerypack(queryPackToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get querypack with ID  "+querypackID+": "+err.Error(),
		)
		return
	}
	queryPackConfJSON, err := json.MarshalIndent(querypackResp.Conf, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Querypack{
		ID:               types.StringValue(querypackResp.ID),
		Name:             types.StringValue(querypackResp.Name),
		Description:      types.StringValue(querypackResp.Description),
		Type:             types.StringValue(querypackResp.Type),
		AdditionalLogger: types.BoolValue(querypackResp.AdditionalLogger),
		IsInternal:       types.BoolValue(querypackResp.IsInternal),
		ResourceType:     types.StringValue(querypackResp.ResourceType),
		Conf:             types.StringValue(string(queryPackConfJSON) + "\n"),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
