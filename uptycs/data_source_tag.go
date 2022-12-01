package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &tagDataSource{}
	_ datasource.DataSourceWithConfigure = &tagDataSource{}
)

func TagDataSource() datasource.DataSource {
	return &tagDataSource{}
}

type tagDataSource struct {
	client *uptycs.Client
}

func (d *tagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (d *tagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *tagDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
			},
			"value": {
				Type:     types.StringType,
				Optional: true,
			},
			"key": {
				Type:     types.StringType,
				Optional: true,
			},
			"flag_profile": {
				Type:     types.StringType,
				Optional: true,
			},
			"custom_profile": {
				Type:     types.StringType,
				Optional: true,
			},
			"compliance_profiled": {
				Type:     types.StringType,
				Optional: true,
			},
			"process_block_rule": {
				Type:     types.StringType,
				Optional: true,
			},
			"dns_block_rule": {
				Type:     types.StringType,
				Optional: true,
			},
			"windows_defender_preference": {
				Type:     types.StringType,
				Optional: true,
			},
			"tag": {
				Type:     types.StringType,
				Optional: true,
			},
			"system": {
				Type:     types.BoolType,
				Optional: true,
			},
			"tag_rule": {
				Type:     types.StringType,
				Optional: true,
			},
			"status": {
				Type:     types.StringType,
				Optional: true,
			},
			"source": {
				Type:     types.StringType,
				Optional: true,
			},
			"resource_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"file_path_groups": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"event_exclude_profiles": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"querypacks": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"registry_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"yara_group_rules": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"audit_configurations": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
		},
	}, nil
}

func (d *tagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tagID string
	var tagKey string
	var tagValue string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &tagID)
	keyAttr := req.Config.GetAttribute(ctx, path.Root("key"), &tagKey)
	valueAttr := req.Config.GetAttribute(ctx, path.Root("value"), &tagValue)

	var tagToLookup uptycs.Tag

	if len(tagID) == 0 {
		resp.Diagnostics.Append(keyAttr...)
		resp.Diagnostics.Append(valueAttr...)
		tagToLookup = uptycs.Tag{
			Key:   tagKey,
			Value: tagValue,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		tagToLookup = uptycs.Tag{
			ID: tagID,
		}
	}

	tagResp, err := d.client.GetTag(tagToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
			"Could not get tag with ID  "+tagID+": "+err.Error(),
		)
		return
	}

	var result = Tag{
		ID:                        types.String{Value: tagResp.ID},
		Value:                     types.String{Value: tagResp.Value},
		Key:                       types.String{Value: tagResp.Key},
		FlagProfile:               types.String{Value: tagResp.FlagProfileID},
		CustomProfile:             types.String{Value: tagResp.CustomProfileID},
		ComplianceProfile:         types.String{Value: tagResp.ComplianceProfileID},
		ProcessBlockRule:          types.String{Value: tagResp.ProcessBlockRuleID},
		DNSBlockRule:              types.String{Value: tagResp.DNSBlockRuleID},
		WindowsDefenderPreference: types.String{Value: tagResp.WindowsDefenderPreferenceID},
		Tag:                       types.String{Value: tagResp.Tag},
		System:                    types.Bool{Value: tagResp.System},
		TagRule:                   types.String{Value: tagResp.TagRuleID},
		Status:                    types.String{Value: tagResp.Status},
		Source:                    types.String{Value: tagResp.Source},
		ResourceType:              types.String{Value: tagResp.ResourceType},
		FilePathGroups: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		EventExcludeProfiles: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		Querypacks: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		RegistryPaths: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		YaraGroupRules: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		AuditConfigurations: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
	}

	for _, fpg := range tagResp.FilePathGroups {
		result.FilePathGroups.Elems = append(result.FilePathGroups.Elems, types.String{Value: fpg.ID})
	}

	for _, eep := range tagResp.EventExcludeProfiles {
		result.EventExcludeProfiles.Elems = append(result.EventExcludeProfiles.Elems, types.String{Value: eep.ID})
	}

	for _, qp := range tagResp.Querypacks {
		result.Querypacks.Elems = append(result.Querypacks.Elems, types.String{Value: qp.ID})
	}

	for _, rp := range tagResp.RegistryPaths {
		result.RegistryPaths.Elems = append(result.RegistryPaths.Elems, types.String{Value: rp.ID})
	}

	for _, yg := range tagResp.YaraGroupRules {
		result.YaraGroupRules.Elems = append(result.YaraGroupRules.Elems, types.String{Value: yg.ID})
	}

	for _, ac := range tagResp.AuditConfigurations {
		result.AuditConfigurations.Elems = append(result.AuditConfigurations.Elems, types.String{Value: ac.ID})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
