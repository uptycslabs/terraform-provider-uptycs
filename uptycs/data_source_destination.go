package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceDestinationType struct {
	p Provider
}

func (r dataSourceDestinationType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"type": {
				Type:     types.StringType,
				Optional: true,
			},
			"address": {
				Type:     types.StringType,
				Optional: true,
			},
			"enabled": {
				Type:     types.BoolType,
				Optional: true,
				//PlanModifiers: tfsdk.AttributePlanModifiers{boolDefault(true)},
			},
		},
	}, nil
}

func (r dataSourceDestinationType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceDestinationType{
		p: *(p.(*Provider)),
	}, nil
}

func (r dataSourceDestinationType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var destinationID string
	path.Root("test")
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &destinationID)...)

	destinationResp, err := r.p.client.GetDestination(uptycs.Destination{
		ID: destinationID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get destination with ID  "+destinationID+": "+err.Error(),
		)
		return
	}

	var result = Destination{
		ID:      types.String{Value: destinationResp.ID},
		Name:    types.String{Value: destinationResp.Name},
		Type:    types.String{Value: destinationResp.Type},
		Address: types.String{Value: destinationResp.Address},
		Enabled: types.Bool{Value: destinationResp.Enabled},
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
