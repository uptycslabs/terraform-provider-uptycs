package uptycs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

type resourceTagType struct{}

// Alert Rule Resource schema
func (r resourceTagType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"key": {
				Type:     types.StringType,
				Optional: true,
			},
			"flag_profile_id": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"custom_profile_id": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"compliance_profile_id": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"process_block_rule_id": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"dns_block_rule_id": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"windows_defender_preference_id": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"tag": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"custom": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"system": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(false)},
			},
			"tag_rule_id": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("")},
			},
			"status": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("active")},
			},
			"source": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("direct")},
			},
			"resource_type": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), stringDefault("asset")},
			},
			"file_path_groups": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"event_exclude_profiles": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"querypacks": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"registry_paths": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"yara_group_rules": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"audit_configurations": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
		},
	}, nil
}

// New resource instance
func (r resourceTagType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return resourceTag{
		p: *(p.(*Provider)),
	}, nil
}

type resourceTag struct {
	p Provider
}

// Read resource information
func (r resourceTag) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tagID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &tagID)...)
	tagResp, err := r.p.client.GetTag(uptycs.Tag{
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
		ID:                          types.String{Value: tagResp.ID},
		Value:                       types.String{Value: tagResp.Value},
		Key:                         types.String{Value: tagResp.Key},
		FlagProfileID:               types.String{Value: tagResp.FlagProfileID},
		CustomProfileID:             types.String{Value: tagResp.CustomProfileID},
		ComplianceProfileID:         types.String{Value: tagResp.ComplianceProfileID},
		ProcessBlockRuleID:          types.String{Value: tagResp.ProcessBlockRuleID},
		DNSBlockRuleID:              types.String{Value: tagResp.DNSBlockRuleID},
		WindowsDefenderPreferenceID: types.String{Value: tagResp.WindowsDefenderPreferenceID},
		Tag:                         types.String{Value: tagResp.Tag},
		Custom:                      types.Bool{Value: tagResp.Custom},
		System:                      types.Bool{Value: tagResp.System},
		TagRuleID:                   types.String{Value: tagResp.TagRuleID},
		Status:                      types.String{Value: tagResp.Status},
		Source:                      types.String{Value: tagResp.Source},
		ResourceType:                types.String{Value: tagResp.ResourceType},
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
	for _, t := range tagResp.FilePathGroups {
		result.FilePathGroups.Elems = append(result.FilePathGroups.Elems, types.String{Value: t.Name})
	}

	for _, eep := range tagResp.EventExcludeProfiles {
		result.EventExcludeProfiles.Elems = append(result.EventExcludeProfiles.Elems, types.String{Value: eep.Name})
	}

	for _, qp := range tagResp.Querypacks {
		result.Querypacks.Elems = append(result.Querypacks.Elems, types.String{Value: qp.Name})
	}

	for _, rp := range tagResp.RegistryPaths {
		result.RegistryPaths.Elems = append(result.RegistryPaths.Elems, types.String{Value: rp.Name})
	}

	for _, yg := range tagResp.YaraGroupRules {
		result.YaraGroupRules.Elems = append(result.YaraGroupRules.Elems, types.String{Value: yg.Name})
	}

	for _, ac := range tagResp.AuditConfigurations {
		result.AuditConfigurations.Elems = append(result.AuditConfigurations.Elems, types.String{Value: ac.Name})
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Create a new resource
func (r resourceTag) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Tag
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var filePathGroups = make([]uptycs.TagConfigurationObject, 0)
	var filePathGroupNames []string
	plan.FilePathGroups.ElementsAs(ctx, &filePathGroupNames, false)
	for _, fpg := range filePathGroupNames {
		filePathGroups = append(filePathGroups, uptycs.TagConfigurationObject{
			Name: fpg,
		})
	}

	var eventExcludeProfiles = make([]uptycs.TagConfigurationObject, 0)
	var eventExcludeProfileNames []string
	plan.EventExcludeProfiles.ElementsAs(ctx, &eventExcludeProfileNames, false)
	for _, eep := range eventExcludeProfileNames {
		eventExcludeProfiles = append(eventExcludeProfiles, uptycs.TagConfigurationObject{
			Name: eep,
		})
	}

	var registryPaths = make([]uptycs.TagConfigurationObject, 0)
	var registryPathNames []string
	plan.RegistryPaths.ElementsAs(ctx, &registryPathNames, false)
	for _, rp := range registryPathNames {
		registryPaths = append(registryPaths, uptycs.TagConfigurationObject{
			Name: rp,
		})
	}

	var queryPacks = make([]uptycs.TagConfigurationObject, 0)
	var querypackNames []string
	plan.Querypacks.ElementsAs(ctx, &querypackNames, false)
	for _, qp := range querypackNames {
		queryPacks = append(queryPacks, uptycs.TagConfigurationObject{
			Name: qp,
		})
	}

	var yaraGroupRules = make([]uptycs.TagConfigurationObject, 0)
	var yaraGroupRuleNames []string
	plan.YaraGroupRules.ElementsAs(ctx, &yaraGroupRuleNames, false)
	for _, ygr := range yaraGroupRuleNames {
		yaraGroupRules = append(yaraGroupRules, uptycs.TagConfigurationObject{
			Name: ygr,
		})
	}

	var auditConfigurations = make([]uptycs.TagConfigurationObject, 0)
	var auditConfigurationNames []string
	plan.AuditConfigurations.ElementsAs(ctx, &auditConfigurationNames, false)
	for _, ac := range auditConfigurationNames {
		auditConfigurations = append(auditConfigurations, uptycs.TagConfigurationObject{
			Name: ac,
		})
	}

	tagResp, err := r.p.client.CreateTag(uptycs.Tag{
		Value:                       plan.Value.Value,
		Key:                         plan.Key.Value,
		FlagProfileID:               plan.FlagProfileID.Value,
		CustomProfileID:             plan.CustomProfileID.Value,
		ComplianceProfileID:         plan.ComplianceProfileID.Value,
		ProcessBlockRuleID:          plan.ProcessBlockRuleID.Value,
		DNSBlockRuleID:              plan.DNSBlockRuleID.Value,
		WindowsDefenderPreferenceID: plan.WindowsDefenderPreferenceID.Value,
		Tag:                         plan.Tag.Value,
		Custom:                      plan.Custom.Value,
		System:                      plan.System.Value,
		TagRuleID:                   plan.TagRuleID.Value,
		Status:                      plan.Status.Value,
		Source:                      plan.Source.Value,
		ResourceType:                plan.ResourceType.Value,
		FilePathGroups:              filePathGroups,
		EventExcludeProfiles:        eventExcludeProfiles,
		RegistryPaths:               registryPaths,
		Querypacks:                  queryPacks,
		YaraGroupRules:              yaraGroupRules,
		AuditConfigurations:         auditConfigurations,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tag, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Tag{
		ID:                          types.String{Value: tagResp.ID},
		Value:                       types.String{Value: tagResp.Value},
		Key:                         types.String{Value: tagResp.Key},
		FlagProfileID:               types.String{Value: tagResp.FlagProfileID},
		CustomProfileID:             types.String{Value: tagResp.CustomProfileID},
		ComplianceProfileID:         types.String{Value: tagResp.ComplianceProfileID},
		ProcessBlockRuleID:          types.String{Value: tagResp.ProcessBlockRuleID},
		DNSBlockRuleID:              types.String{Value: tagResp.DNSBlockRuleID},
		WindowsDefenderPreferenceID: types.String{Value: tagResp.WindowsDefenderPreferenceID},
		Tag:                         types.String{Value: tagResp.Tag},
		Custom:                      types.Bool{Value: tagResp.Custom},
		System:                      types.Bool{Value: tagResp.System},
		TagRuleID:                   types.String{Value: tagResp.TagRuleID},
		Status:                      types.String{Value: tagResp.Status},
		Source:                      types.String{Value: tagResp.Source},
		ResourceType:                types.String{Value: tagResp.ResourceType},
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

	for _, _fpg := range tagResp.FilePathGroups {
		result.FilePathGroups.Elems = append(result.FilePathGroups.Elems, types.String{Value: _fpg.Name})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceTag) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	var filePathGroupNames []string
	plan.FilePathGroups.ElementsAs(ctx, &filePathGroupNames, false)
	for _, fpg := range filePathGroupNames {
		filePathGroups = append(filePathGroups, uptycs.TagConfigurationObject{
			Name: fpg,
		})
	}

	var eventExcludeProfiles = make([]uptycs.TagConfigurationObject, 0)
	var eventExcludeProfileNames []string
	plan.EventExcludeProfiles.ElementsAs(ctx, &eventExcludeProfileNames, false)
	for _, eep := range eventExcludeProfileNames {
		eventExcludeProfiles = append(eventExcludeProfiles, uptycs.TagConfigurationObject{
			Name: eep,
		})
	}

	var registryPaths = make([]uptycs.TagConfigurationObject, 0)
	var registryPathNames []string
	plan.RegistryPaths.ElementsAs(ctx, &registryPathNames, false)
	for _, rp := range registryPathNames {
		registryPaths = append(registryPaths, uptycs.TagConfigurationObject{
			Name: rp,
		})
	}

	var queryPacks = make([]uptycs.TagConfigurationObject, 0)
	var querypackNames []string
	plan.Querypacks.ElementsAs(ctx, &querypackNames, false)
	for _, qp := range querypackNames {
		queryPacks = append(queryPacks, uptycs.TagConfigurationObject{
			Name: qp,
		})
	}

	var yaraGroupRules = make([]uptycs.TagConfigurationObject, 0)
	var yaraGroupRuleNames []string
	plan.YaraGroupRules.ElementsAs(ctx, &yaraGroupRuleNames, false)
	for _, ygr := range yaraGroupRuleNames {
		yaraGroupRules = append(yaraGroupRules, uptycs.TagConfigurationObject{
			Name: ygr,
		})
	}

	var auditConfigurations = make([]uptycs.TagConfigurationObject, 0)
	var auditConfigurationNames []string
	plan.AuditConfigurations.ElementsAs(ctx, &auditConfigurationNames, false)
	for _, ac := range auditConfigurationNames {
		auditConfigurations = append(auditConfigurations, uptycs.TagConfigurationObject{
			Name: ac,
		})
	}

	tagResp, err := r.p.client.UpdateTag(uptycs.Tag{
		ID:                          tagID,
		Value:                       plan.Value.Value,
		Key:                         plan.Key.Value,
		FlagProfileID:               plan.FlagProfileID.Value,
		CustomProfileID:             plan.CustomProfileID.Value,
		ComplianceProfileID:         plan.ComplianceProfileID.Value,
		ProcessBlockRuleID:          plan.ProcessBlockRuleID.Value,
		DNSBlockRuleID:              plan.DNSBlockRuleID.Value,
		WindowsDefenderPreferenceID: plan.WindowsDefenderPreferenceID.Value,
		Tag:                         plan.Tag.Value,
		Custom:                      plan.Custom.Value,
		System:                      plan.System.Value,
		TagRuleID:                   plan.TagRuleID.Value,
		Status:                      plan.Status.Value,
		Source:                      plan.Source.Value,
		ResourceType:                plan.ResourceType.Value,
		FilePathGroups:              filePathGroups,
		EventExcludeProfiles:        eventExcludeProfiles,
		RegistryPaths:               registryPaths,
		Querypacks:                  queryPacks,
		YaraGroupRules:              yaraGroupRules,
		AuditConfigurations:         auditConfigurations,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create tag, unexpected error: "+err.Error(),
		)
		return
	}

	var result = Tag{
		ID:                          types.String{Value: tagResp.ID},
		Value:                       types.String{Value: tagResp.Value},
		Key:                         types.String{Value: tagResp.Key},
		FlagProfileID:               types.String{Value: tagResp.FlagProfileID},
		CustomProfileID:             types.String{Value: tagResp.CustomProfileID},
		ComplianceProfileID:         types.String{Value: tagResp.ComplianceProfileID},
		ProcessBlockRuleID:          types.String{Value: tagResp.ProcessBlockRuleID},
		DNSBlockRuleID:              types.String{Value: tagResp.DNSBlockRuleID},
		WindowsDefenderPreferenceID: types.String{Value: tagResp.WindowsDefenderPreferenceID},
		Tag:                         types.String{Value: tagResp.Tag},
		Custom:                      types.Bool{Value: tagResp.Custom},
		System:                      types.Bool{Value: tagResp.System},
		TagRuleID:                   types.String{Value: tagResp.TagRuleID},
		Status:                      types.String{Value: tagResp.Status},
		Source:                      types.String{Value: tagResp.Source},
		ResourceType:                types.String{Value: tagResp.ResourceType},
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

	for _, _fpg := range tagResp.FilePathGroups {
		result.FilePathGroups.Elems = append(result.FilePathGroups.Elems, types.String{Value: _fpg.Name})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceTag) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Tag
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID := state.ID.Value

	_, err := r.p.client.DeleteTag(uptycs.Tag{
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

// Import resource
func (r resourceTag) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
