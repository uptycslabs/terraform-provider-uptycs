package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func TagResource() resource.Resource {
	return &tagResource{}
}

type tagResource struct {
	client *uptycs.Client
}

func (r *tagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *tagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *tagResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Optional: true,
				Computed: true,
			},
			"value": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"key": schema.StringAttribute{Optional: true},
			"flag_profile": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"custom_profile": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"compliance_profile": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"process_block_rule": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"dns_block_rule": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"windows_defender_preference": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"tag": schema.StringAttribute{Optional: true,
				Computed: true,
			},
			"system": schema.BoolAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"tag_rule": schema.StringAttribute{Optional: true,
				Computed: true,
			},
			"status": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("active")},
			},
			"source": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("direct")},
			},
			"resource_type": schema.StringAttribute{Optional: true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("asset")},
			},
			"file_path_groups": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"event_exclude_profiles": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"querypacks": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"registry_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"yara_group_rules": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"audit_configurations": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tagID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &tagID)...)
	tagResp, err := r.client.GetTag(uptycs.Tag{
		ID: tagID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
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
		TagRule:                   types.StringValue(tagResp.TagRuleID),
		Tag:                       types.StringValue(tagResp.Tag),
		System:                    types.BoolValue(tagResp.System),
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

func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Tag
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var filePathGroups = make([]uptycs.TagConfigurationObject, 0)
	var filePathGroupIDs []string
	plan.FilePathGroups.ElementsAs(ctx, &filePathGroupIDs, false)
	for _, fpg := range filePathGroupIDs {
		filePathGroups = append(filePathGroups, uptycs.TagConfigurationObject{
			ID: fpg,
		})
	}

	var eventExcludeProfiles = make([]uptycs.TagConfigurationObject, 0)
	var eventExcludeProfileIDs []string
	plan.EventExcludeProfiles.ElementsAs(ctx, &eventExcludeProfileIDs, false)
	for _, eep := range eventExcludeProfileIDs {
		eventExcludeProfiles = append(eventExcludeProfiles, uptycs.TagConfigurationObject{
			ID: eep,
		})
	}

	var registryPaths = make([]uptycs.TagConfigurationObject, 0)
	var registryPathIDs []string
	plan.RegistryPaths.ElementsAs(ctx, &registryPathIDs, false)
	for _, rp := range registryPathIDs {
		registryPaths = append(registryPaths, uptycs.TagConfigurationObject{
			ID: rp,
		})
	}

	var queryPacks = make([]uptycs.TagConfigurationObject, 0)
	var querypackIDs []string
	plan.Querypacks.ElementsAs(ctx, &querypackIDs, false)
	for _, qp := range querypackIDs {
		queryPacks = append(queryPacks, uptycs.TagConfigurationObject{
			ID: qp,
		})
	}

	var yaraGroupRules = make([]uptycs.TagConfigurationObject, 0)
	var yaraGroupRuleIDs []string
	plan.YaraGroupRules.ElementsAs(ctx, &yaraGroupRuleIDs, false)
	for _, ygr := range yaraGroupRuleIDs {
		yaraGroupRules = append(yaraGroupRules, uptycs.TagConfigurationObject{
			ID: ygr,
		})
	}

	var auditConfigurations = make([]uptycs.TagConfigurationObject, 0)
	var auditConfigurationIDs []string
	plan.AuditConfigurations.ElementsAs(ctx, &auditConfigurationIDs, false)
	for _, ac := range auditConfigurationIDs {
		auditConfigurations = append(auditConfigurations, uptycs.TagConfigurationObject{
			ID: ac,
		})
	}

	tagResp, err := r.client.CreateTag(uptycs.Tag{
		Value:                     plan.Value.Value,
		Key:                       plan.Key.Value,
		FlagProfile:               plan.FlagProfile.Value,
		CustomProfile:             plan.CustomProfile.Value,
		ComplianceProfile:         plan.ComplianceProfile.Value,
		ProcessBlockRule:          plan.ProcessBlockRule.Value,
		DNSBlockRule:              plan.DNSBlockRule.Value,
		WindowsDefenderPreference: plan.WindowsDefenderPreference.Value,
		TagRule:                   plan.TagRule.Value,
		Tag:                       plan.Tag.Value,
		System:                    plan.System.Value,
		Status:                    plan.Status.Value,
		Source:                    plan.Source.Value,
		ResourceType:              plan.ResourceType.Value,
		FilePathGroups:            filePathGroups,
		EventExcludeProfiles:      eventExcludeProfiles,
		RegistryPaths:             registryPaths,
		Querypacks:                queryPacks,
		YaraGroupRules:            yaraGroupRules,
		AuditConfigurations:       auditConfigurations,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tag, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Tag{
		ID:                        types.StringValue(tagResp.ID),
		Value:                     types.StringValue(tagResp.Value),
		Key:                       types.StringValue(tagResp.Key),
		ComplianceProfile:         types.StringValue(tagResp.ComplianceProfileID),
		ProcessBlockRule:          types.StringValue(tagResp.ProcessBlockRuleID),
		FlagProfile:               types.StringValue(tagResp.FlagProfileID),
		CustomProfile:             types.StringValue(tagResp.CustomProfileID),
		DNSBlockRule:              types.StringValue(tagResp.DNSBlockRuleID),
		WindowsDefenderPreference: types.StringValue(tagResp.WindowsDefenderPreferenceID),
		TagRule:                   types.StringValue(tagResp.TagRuleID),
		Tag:                       types.StringValue(tagResp.Tag),
		System:                    types.BoolValue(tagResp.System),
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

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state Tag
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID := state.ID.Value

	// Retrieve values from plan
	var plan Tag
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var filePathGroups = make([]uptycs.TagConfigurationObject, 0)
	var filePathGroupIDs []string
	plan.FilePathGroups.ElementsAs(ctx, &filePathGroupIDs, false)
	for _, fpg := range filePathGroupIDs {
		filePathGroups = append(filePathGroups, uptycs.TagConfigurationObject{
			ID: fpg,
		})
	}

	var eventExcludeProfiles = make([]uptycs.TagConfigurationObject, 0)
	var eventExcludeProfileIDs []string
	plan.EventExcludeProfiles.ElementsAs(ctx, &eventExcludeProfileIDs, false)
	for _, eep := range eventExcludeProfileIDs {
		eventExcludeProfiles = append(eventExcludeProfiles, uptycs.TagConfigurationObject{
			ID: eep,
		})
	}

	var registryPaths = make([]uptycs.TagConfigurationObject, 0)
	var registryPathIDs []string
	plan.RegistryPaths.ElementsAs(ctx, &registryPathIDs, false)
	for _, rp := range registryPathIDs {
		registryPaths = append(registryPaths, uptycs.TagConfigurationObject{
			ID: rp,
		})
	}

	var queryPacks = make([]uptycs.TagConfigurationObject, 0)
	var querypackIDs []string
	plan.Querypacks.ElementsAs(ctx, &querypackIDs, false)
	for _, qp := range querypackIDs {
		queryPacks = append(queryPacks, uptycs.TagConfigurationObject{
			ID: qp,
		})
	}

	var yaraGroupRules = make([]uptycs.TagConfigurationObject, 0)
	var yaraGroupRuleIDs []string
	plan.YaraGroupRules.ElementsAs(ctx, &yaraGroupRuleIDs, false)
	for _, ygr := range yaraGroupRuleIDs {
		yaraGroupRules = append(yaraGroupRules, uptycs.TagConfigurationObject{
			ID: ygr,
		})
	}

	var auditConfigurations = make([]uptycs.TagConfigurationObject, 0)
	var auditConfigurationIDs []string
	plan.AuditConfigurations.ElementsAs(ctx, &auditConfigurationIDs, false)
	for _, ac := range auditConfigurationIDs {
		auditConfigurations = append(auditConfigurations, uptycs.TagConfigurationObject{
			ID: ac,
		})
	}

	tagResp, err := r.client.UpdateTag(uptycs.Tag{
		ID:                        tagID,
		Value:                     plan.Value.Value,
		Key:                       plan.Key.Value,
		FlagProfile:               plan.FlagProfile.Value,
		CustomProfile:             plan.CustomProfile.Value,
		ComplianceProfile:         plan.ComplianceProfile.Value,
		ProcessBlockRule:          plan.ProcessBlockRule.Value,
		DNSBlockRule:              plan.DNSBlockRule.Value,
		WindowsDefenderPreference: plan.WindowsDefenderPreference.Value,
		TagRule:                   plan.TagRule.Value,
		Tag:                       plan.Tag.Value,
		System:                    plan.System.Value,
		FilePathGroups:            filePathGroups,
		EventExcludeProfiles:      eventExcludeProfiles,
		RegistryPaths:             registryPaths,
		Querypacks:                queryPacks,
		YaraGroupRules:            yaraGroupRules,
		AuditConfigurations:       auditConfigurations,
		//ResourceType:                plan.ResourceType.Value, //│ {"error":{"status":400,"code":"INVALID_OR_REQUIRED_FIELD","message":{"brief":"","detail":"\"resourceType\" is not allowed","developer":""}}}
		//Status:                      plan.Status.Value,  // {"error":{"status":400,"code":"INVALID_OR_REQUIRED_FIELD","message":{"brief":"","detail":"\"status\" is│ not allowed","developer":""}}}
		//Source:                      plan.Source.Value,  // {"error":{"status":400,"code":"INVALID_OR_REQUIRED_FIELD","message":{"brief":"","detail":"\"source\" is│ not allowed","developer":""}}}
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tag, unexpected error: "+err.Error(),
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

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Tag
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID := state.ID.Value

	_, err := r.client.DeleteTag(uptycs.Tag{
		ID: tagID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete tag with ID  "+tagID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
