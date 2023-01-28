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

func ExceptionDataSource() datasource.DataSource {
	return &exceptionDataSource{}
}

type exceptionDataSource struct {
	client *uptycs.Client
}

func (d *exceptionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception"
}

func (d *exceptionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *exceptionDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                schema.StringAttribute{Optional: true},
			"name":              schema.StringAttribute{Optional: true},
			"description":       schema.StringAttribute{Optional: true},
			"exception_type":    schema.StringAttribute{Optional: true},
			"table_name":        schema.StringAttribute{Optional: true},
			"is_global":         schema.BoolAttribute{Optional: true},
			"disabled":          schema.BoolAttribute{Optional: true},
			"close_open_alerts": schema.BoolAttribute{Optional: true},
			"rule":              schema.StringAttribute{Optional: true},
		},
	}
}

func (d *exceptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var exceptionID string
	var exceptionName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &exceptionID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &exceptionName)

	var exceptionToLookup uptycs.Exception

	if len(exceptionID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		exceptionToLookup = uptycs.Exception{
			Name: exceptionName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		exceptionToLookup = uptycs.Exception{
			ID: exceptionID,
		}
	}
	exceptionResp, err := d.client.GetException(exceptionToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get exception with ID  "+exceptionID+": "+err.Error(),
		)
		return
	}

	ruleJSON, err := json.MarshalIndent(exceptionResp.Rule, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = Exception{
		ID:              types.StringValue(exceptionResp.ID),
		Name:            types.StringValue(exceptionResp.Name),
		Description:     types.StringValue(exceptionResp.Description),
		ExceptionType:   types.StringValue(exceptionResp.ExceptionType),
		TableName:       types.StringValue(exceptionResp.TableName),
		IsGlobal:        types.BoolValue(exceptionResp.IsGlobal),
		Disabled:        types.BoolValue(exceptionResp.Disabled),
		CloseOpenAlerts: types.BoolValue(exceptionResp.CloseOpenAlerts),
		Rule:            types.StringValue(string([]byte(ruleJSON)) + "\n"),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
