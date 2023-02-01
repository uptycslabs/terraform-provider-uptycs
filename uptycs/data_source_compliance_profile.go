package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func ComplianceProfileDataSource() datasource.DataSource {
	return &complianceProfileDataSource{}
}

type complianceProfileDataSource struct {
	client *uptycs.Client
}

func (d *complianceProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compliance_profile"
}

func (d *complianceProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *complianceProfileDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"priority":    schema.Int64Attribute{Optional: true},
		},
	}
}

func (d *complianceProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var complianceProfileID string
	var complianceProfileName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &complianceProfileID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &complianceProfileName)

	var complianceProfileToLookup uptycs.ComplianceProfile

	if len(complianceProfileID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		complianceProfileToLookup = uptycs.ComplianceProfile{
			Name: complianceProfileName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		complianceProfileToLookup = uptycs.ComplianceProfile{
			ID: complianceProfileID,
		}
	}

	complianceProfileResp, err := d.client.GetComplianceProfile(complianceProfileToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get complianceProfile with ID  "+complianceProfileID+": "+err.Error(),
		)
		return
	}

	var result = ComplianceProfile{
		ID:          types.StringValue(complianceProfileResp.ID),
		Name:        types.StringValue(complianceProfileResp.Name),
		Description: types.StringValue(complianceProfileResp.Description),
		Priority:    types.Int64Value(int64(complianceProfileResp.Priority)),
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
