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
	_ datasource.DataSource              = &complianceProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &complianceProfileDataSource{}
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

func (d *complianceProfileDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"priority": {
				Type:     types.NumberType,
				Optional: true,
			},
		},
	}, nil
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
		ID:          types.String{Value: complianceProfileResp.ID},
		Name:        types.String{Value: complianceProfileResp.Name},
		Description: types.String{Value: complianceProfileResp.Description},
		Priority:    complianceProfileResp.Priority,
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
