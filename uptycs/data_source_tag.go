package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (d *tagDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                          schema.StringAttribute{Optional: true},
			"value":                       schema.StringAttribute{Optional: true},
			"key":                         schema.StringAttribute{Optional: true},
			"flag_profile":                schema.StringAttribute{Optional: true},
			"custom_profile":              schema.StringAttribute{Optional: true},
			"compliance_profiled":         schema.StringAttribute{Optional: true},
			"process_block_rule":          schema.StringAttribute{Optional: true},
			"dns_block_rule":              schema.StringAttribute{Optional: true},
			"windows_defender_preference": schema.StringAttribute{Optional: true},
			"tag":                         schema.StringAttribute{Optional: true},
			"system":                      schema.BoolAttribute{Optional: true},
			"tag_rule":                    schema.StringAttribute{Optional: true},
			"status":                      schema.StringAttribute{Optional: true},
			"source":                      schema.StringAttribute{Optional: true},
			"resource_type":               schema.StringAttribute{Optional: true},
			"file_path_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"event_exclude_profiles": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"querypacks": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"registry_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"yara_group_rules": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"audit_configurations": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
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
		ID:                        types.StringValue(tagResp.ID),
		Value:                     types.StringValue(tagResp.Value),
		Key:                       types.StringValue(tagResp.Key),
		FlagProfile:               types.StringValue(tagResp.FlagProfileID),
		CustomProfile:             types.StringValue(tagResp.CustomProfileID),
		ComplianceProfile:         types.StringValue(tagResp.ComplianceProfileID),
		ProcessBlockRule:          types.StringValue(tagResp.ProcessBlockRuleID),
		DNSBlockRule:              types.StringValue(tagResp.DNSBlockRuleID),
		WindowsDefenderPreference: types.StringValue(tagResp.WindowsDefenderPreferenceID),
		Tag:                       types.StringValue(tagResp.Tag),
		System:                    types.BoolValue(tagResp.System),
		TagRule:                   types.StringValue(tagResp.TagRuleID),
		Status:                    types.StringValue(tagResp.Status),
		Source:                    types.StringValue(tagResp.Source),
		ResourceType:              types.StringValue(tagResp.ResourceType),
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
