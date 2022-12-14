package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &auditConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &auditConfigurationDataSource{}
)

func AuditConfigurationDataSource() datasource.DataSource {
	return &auditConfigurationDataSource{}
}

type auditConfigurationDataSource struct {
	client *uptycs.Client
}

func (d *auditConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_configuration"
}

func (d *auditConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *auditConfigurationDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (d *auditConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var auditConfigurationID string
	var auditConfigurationName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &auditConfigurationID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &auditConfigurationName)

	var auditConfigurationToLookup uptycs.AuditConfiguration

	if len(auditConfigurationID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		auditConfigurationToLookup = uptycs.AuditConfiguration{
			Name: auditConfigurationName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		auditConfigurationToLookup = uptycs.AuditConfiguration{
			ID: auditConfigurationID,
		}
	}

	auditConfigurationResp, err := d.client.GetAuditConfiguration(auditConfigurationToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get auditConfiguration with ID  "+auditConfigurationID+": "+err.Error(),
		)
		return
	}

	var result = AuditConfiguration{
		ID:          types.StringValue(auditConfigurationResp.ID),
		Name:        types.StringValue(auditConfigurationResp.Name),
		Description: types.StringValue(auditConfigurationResp.Description),
		Framework:   types.StringValue(auditConfigurationResp.Framework),
		Version:     types.StringValue(auditConfigurationResp.Version),
		OsVersion:   types.StringValue(auditConfigurationResp.OsVersion),
		Platform:    types.StringValue(auditConfigurationResp.Platform),
		TableName:   types.StringValue(auditConfigurationResp.TableName),
		Sha256:      types.StringValue(auditConfigurationResp.Sha256),
		Type:        types.StringValue(auditConfigurationResp.Type),
		Checks:      auditConfigurationResp.Checks,
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
