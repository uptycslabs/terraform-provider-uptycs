package uptycs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type dataSourceDestinationType struct {
	p provider
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

func (d dataSourceDestinationType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceDestinationType{
		p: *(p.(*provider)),
	}, nil
}

func (d dataSourceDestinationType) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var destinationId string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), &destinationId)...)

	destinationResp, err := d.p.client.GetDestination(uptycs.Destination{
		ID: destinationId,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading order",
			"Could not get destination with ID  "+destinationId+": "+err.Error(),
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
