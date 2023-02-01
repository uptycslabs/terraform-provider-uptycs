package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func DestinationDataSource() datasource.DataSource {
	return &destinationDataSource{}
}

type destinationDataSource struct {
	client *uptycs.Client
}

func (d *destinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_destination"
}

func (d *destinationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *destinationDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Optional: true},
			"name":    schema.StringAttribute{Optional: true},
			"type":    schema.StringAttribute{Optional: true},
			"address": schema.StringAttribute{Optional: true},
			"enabled": schema.BoolAttribute{Optional: true}, //PlanModifiers: tfsdk.AttributePlanModifiers{modifiers.DefaultBool(true)},

		},
	}
}

func (d *destinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var destinationID string
	var destinationName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &destinationID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &destinationName)

	var destinationToLookup uptycs.Destination

	if len(destinationID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		destinationToLookup = uptycs.Destination{
			Name: destinationName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		destinationToLookup = uptycs.Destination{
			ID: destinationID,
		}
	}

	destinationResp, err := d.client.GetDestination(destinationToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get destination with ID  "+destinationID+": "+err.Error(),
		)
		return
	}

	var result = Destination{
		ID:      types.StringValue(destinationResp.ID),
		Name:    types.StringValue(destinationResp.Name),
		Type:    types.StringValue(destinationResp.Type),
		Address: types.StringValue(destinationResp.Address),
		Enabled: types.BoolValue(destinationResp.Enabled),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
