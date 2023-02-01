package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (d *auditConfigurationDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"framework":   schema.StringAttribute{Optional: true},
			"version":     schema.StringAttribute{Optional: true},
			"os_version":  schema.StringAttribute{Optional: true},
			"platform":    schema.StringAttribute{Optional: true},
			"table_name":  schema.StringAttribute{Optional: true},
			"sha256":      schema.StringAttribute{Optional: true},
			"type":        schema.StringAttribute{Optional: true},
			"checks":      schema.Int64Attribute{Optional: true},
		},
	}
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
		Checks:      types.Int64Value(int64(auditConfigurationResp.Checks)),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
