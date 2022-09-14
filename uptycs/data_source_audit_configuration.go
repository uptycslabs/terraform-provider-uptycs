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

type dataSourceAuditConfigurationType struct {
	p Provider
}

func (r dataSourceAuditConfigurationType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"framework": {
				Type:     types.StringType,
				Optional: true,
			},
			"version": {
				Type:     types.StringType,
				Optional: true,
			},
			"os_version": {
				Type:     types.StringType,
				Optional: true,
			},
			"platform": {
				Type:     types.StringType,
				Optional: true,
			},
			"table_name": {
				Type:     types.StringType,
				Optional: true,
			},
			"sha256": {
				Type:     types.StringType,
				Optional: true,
			},
			"type": {
				Type:     types.StringType,
				Optional: true,
			},
			"checks": {
				Type:     types.NumberType,
				Optional: true,
			},
		},
	}, nil
}

func (r dataSourceAuditConfigurationType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return dataSourceAuditConfigurationType{
		p: *(p.(*Provider)),
	}, nil
}

func (r dataSourceAuditConfigurationType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var auditConfigurationID string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &auditConfigurationID)...)

	auditConfigurationResp, err := r.p.client.GetAuditConfiguration(uptycs.AuditConfiguration{
		ID: auditConfigurationID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get auditConfiguration with ID  "+auditConfigurationID+": "+err.Error(),
		)
		return
	}

	var result = AuditConfiguration{
		ID:          types.String{Value: auditConfigurationResp.ID},
		Name:        types.String{Value: auditConfigurationResp.Name},
		Description: types.String{Value: auditConfigurationResp.Description},
		Framework:   types.String{Value: auditConfigurationResp.Framework},
		Version:     types.String{Value: auditConfigurationResp.Version},
		OsVersion:   types.String{Value: auditConfigurationResp.OsVersion},
		Platform:    types.String{Value: auditConfigurationResp.Platform},
		TableName:   types.String{Value: auditConfigurationResp.TableName},
		Sha256:      types.String{Value: auditConfigurationResp.Sha256},
		Type:        types.String{Value: auditConfigurationResp.Type},
		Checks:      auditConfigurationResp.Checks,
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
