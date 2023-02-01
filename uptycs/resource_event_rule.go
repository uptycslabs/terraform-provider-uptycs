package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

func EventRuleResource() resource.Resource {
	return &eventRuleResource{}
}

type eventRuleResource struct {
	client *uptycs.Client
}

func (r *eventRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_rule"
}

func (r *eventRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*uptycs.Client)
}

func (r *eventRuleResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Required: true},
			"description": schema.StringAttribute{Required: true},
			"score": schema.StringAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringDefault(""),
				},
			},
			"code": schema.StringAttribute{Required: true,
				Computed: false,
			},
			"type": schema.StringAttribute{Required: true,
				Computed: false,
			},
			"rule":        schema.StringAttribute{Required: true},
			"grouping":    schema.StringAttribute{Required: true},
			"grouping_l2": schema.StringAttribute{Optional: true},
			"grouping_l3": schema.StringAttribute{Optional: true},
			"enabled": schema.BoolAttribute{Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolDefault(false),
				},
			},
			"event_tags": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"builder_config": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"table_name":     schema.StringAttribute{Optional: true},
					"added":          schema.BoolAttribute{Optional: true},
					"matches_filter": schema.BoolAttribute{Optional: true},
					"severity":       schema.StringAttribute{Optional: true},
					"key":            schema.StringAttribute{Optional: true},
					"value_field":    schema.StringAttribute{Optional: true},
					"filters":        schema.StringAttribute{Required: true},
					"auto_alert_config": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"raise_alert":      schema.BoolAttribute{Required: true},
							"disable_alert":    schema.BoolAttribute{Required: true},
							"metadata_sources": schema.StringAttribute{Optional: true},
						},
					},
				},
			},
		},
	}
}

func (r *eventRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var eventRuleID string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &eventRuleID)...)
	eventRuleResp, err := r.client.GetEventRule(uptycs.EventRule{
		ID: eventRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading",
			"Could not get eventRule with ID  "+eventRuleID+": "+err.Error(),
		)
		return
	}

	filtersJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.Filters, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	metadataJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.AutoAlertConfig.MetadataSources, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventRule{
		ID:          types.StringValue(eventRuleResp.ID),
		Enabled:     types.BoolValue(eventRuleResp.Enabled),
		Name:        types.StringValue(eventRuleResp.Name),
		Description: types.StringValue(eventRuleResp.Description),
		Code:        types.StringValue(eventRuleResp.Code),
		Type:        types.StringValue(eventRuleResp.Type),
		Rule:        types.StringValue(eventRuleResp.Rule),
		Grouping:    types.StringValue(eventRuleResp.Grouping),
		GroupingL2:  types.StringValue(eventRuleResp.GroupingL2),
		GroupingL3:  types.StringValue(eventRuleResp.GroupingL3),
		Score:       types.StringValue(eventRuleResp.Score),
		EventTags:   makeListStringAttribute(eventRuleResp.EventTags),
		BuilderConfig: BuilderConfig{
			Filters:       types.StringValue(string(filtersJSON) + "\n"),
			TableName:     types.StringValue(eventRuleResp.BuilderConfig.TableName),
			Added:         types.BoolValue(eventRuleResp.BuilderConfig.Added),
			MatchesFilter: types.BoolValue(eventRuleResp.BuilderConfig.MatchesFilter),
			Severity:      types.StringValue(eventRuleResp.BuilderConfig.Severity),
			Key:           types.StringValue(eventRuleResp.BuilderConfig.Key),
			ValueField:    types.StringValue(eventRuleResp.BuilderConfig.ValueField),
			AutoAlertConfig: AutoAlertConfig{
				DisableAlert:    types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.DisableAlert),
				RaiseAlert:      types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert),
				MetadataSources: types.StringValue(string(metadataJSON) + "\n"),
			},
		},
	}

	if result.Type.ValueString() == "sql" {
		result.Rule = types.StringValue(fmt.Sprintf("%s\n", result.Rule.ValueString()))
	}

	diags := resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *eventRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan EventRule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tags []string
	plan.EventTags.ElementsAs(ctx, &tags, false)

	eventRuleResp, err := r.client.CreateEventRule(uptycs.EventRule{
		Name:        plan.Name.ValueString(),
		Code:        plan.Code.ValueString(),
		Description: plan.Description.ValueString(),
		Rule:        plan.Rule.ValueString(),
		Type:        plan.Type.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
		Grouping:    plan.Grouping.ValueString(),
		GroupingL2:  plan.GroupingL2.ValueString(),
		GroupingL3:  plan.GroupingL3.ValueString(),
		EventTags:   tags,
		Score:       plan.Score.ValueString(),
		BuilderConfig: uptycs.BuilderConfig{
			Filters:       uptycs.CustomJSONString(plan.BuilderConfig.Filters.ValueString()),
			TableName:     plan.BuilderConfig.TableName.ValueString(),
			Added:         plan.BuilderConfig.Added.ValueBool(),
			MatchesFilter: plan.BuilderConfig.MatchesFilter.ValueBool(),
			Severity:      plan.BuilderConfig.Severity.ValueString(),
			Key:           plan.BuilderConfig.Key.ValueString(),
			ValueField:    plan.BuilderConfig.ValueField.ValueString(),
			AutoAlertConfig: uptycs.AutoAlertConfig{
				DisableAlert:    plan.BuilderConfig.AutoAlertConfig.DisableAlert.ValueBool(),
				RaiseAlert:      plan.BuilderConfig.AutoAlertConfig.RaiseAlert.ValueBool(),
				MetadataSources: uptycs.CustomJSONString(plan.BuilderConfig.AutoAlertConfig.MetadataSources.ValueString()),
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create eventRule, unexpected error: "+err.Error(),
		)
		return
	}

	filtersJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.Filters, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	metadataJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.AutoAlertConfig.MetadataSources, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventRule{
		ID:          types.StringValue(eventRuleResp.ID),
		Enabled:     types.BoolValue(eventRuleResp.Enabled),
		Name:        types.StringValue(eventRuleResp.Name),
		Description: types.StringValue(eventRuleResp.Description),
		Code:        types.StringValue(eventRuleResp.Code),
		Type:        types.StringValue(eventRuleResp.Type),
		Rule:        types.StringValue(eventRuleResp.Rule),
		Grouping:    types.StringValue(eventRuleResp.Grouping),
		GroupingL2:  types.StringValue(eventRuleResp.GroupingL2),
		GroupingL3:  types.StringValue(eventRuleResp.GroupingL3),
		Score:       types.StringValue(eventRuleResp.Score),
		EventTags:   makeListStringAttribute(eventRuleResp.EventTags),
		BuilderConfig: BuilderConfig{
			Filters:       types.StringValue(string(filtersJSON) + "\n"),
			TableName:     types.StringValue(eventRuleResp.BuilderConfig.TableName),
			Added:         types.BoolValue(eventRuleResp.BuilderConfig.Added),
			MatchesFilter: types.BoolValue(eventRuleResp.BuilderConfig.MatchesFilter),
			Severity:      types.StringValue(eventRuleResp.BuilderConfig.Severity),
			Key:           types.StringValue(eventRuleResp.BuilderConfig.Key),
			ValueField:    types.StringValue(eventRuleResp.BuilderConfig.ValueField),
			AutoAlertConfig: AutoAlertConfig{
				DisableAlert:    types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.DisableAlert),
				RaiseAlert:      types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert),
				MetadataSources: types.StringValue(string(metadataJSON) + "\n"),
			},
		},
	}
	if result.Type.ValueString() == "sql" {
		result.Rule = types.StringValue(fmt.Sprintf("%s\n", result.Rule.ValueString()))
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *eventRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state EventRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventRuleID := state.ID.ValueString()

	// Retrieve values from plan
	var plan EventRule
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tags []string
	plan.EventTags.ElementsAs(ctx, &tags, false)

	eventRuleResp, err := r.client.UpdateEventRule(uptycs.EventRule{
		ID:          eventRuleID,
		Name:        plan.Name.ValueString(),
		Code:        plan.Code.ValueString(),
		Description: plan.Description.ValueString(),
		Rule:        plan.Rule.ValueString(),
		Type:        plan.Type.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
		Grouping:    plan.Grouping.ValueString(),
		GroupingL2:  plan.GroupingL2.ValueString(),
		GroupingL3:  plan.GroupingL3.ValueString(),
		EventTags:   tags,
		Score:       plan.Score.ValueString(),
		BuilderConfig: uptycs.BuilderConfig{
			Filters:       uptycs.CustomJSONString(plan.BuilderConfig.Filters.ValueString()),
			TableName:     plan.BuilderConfig.TableName.ValueString(),
			Added:         plan.BuilderConfig.Added.ValueBool(),
			MatchesFilter: plan.BuilderConfig.MatchesFilter.ValueBool(),
			Severity:      plan.BuilderConfig.Severity.ValueString(),
			Key:           plan.BuilderConfig.Key.ValueString(),
			ValueField:    plan.BuilderConfig.ValueField.ValueString(),
			AutoAlertConfig: uptycs.AutoAlertConfig{
				DisableAlert:    plan.BuilderConfig.AutoAlertConfig.DisableAlert.ValueBool(),
				RaiseAlert:      plan.BuilderConfig.AutoAlertConfig.RaiseAlert.ValueBool(),
				MetadataSources: uptycs.CustomJSONString(plan.BuilderConfig.AutoAlertConfig.MetadataSources.ValueString()),
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating",
			"Could not create eventRule, unexpected error: "+err.Error(),
		)
		return
	}

	filtersJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.Filters, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	metadataJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.AutoAlertConfig.MetadataSources, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	var result = EventRule{
		ID:          types.StringValue(eventRuleResp.ID),
		Enabled:     types.BoolValue(eventRuleResp.Enabled),
		Name:        types.StringValue(eventRuleResp.Name),
		Description: types.StringValue(eventRuleResp.Description),
		Code:        types.StringValue(eventRuleResp.Code),
		Type:        types.StringValue(eventRuleResp.Type),
		Rule:        types.StringValue(eventRuleResp.Rule),
		Grouping:    types.StringValue(eventRuleResp.Grouping),
		GroupingL2:  types.StringValue(eventRuleResp.GroupingL2),
		GroupingL3:  types.StringValue(eventRuleResp.GroupingL3),
		Score:       types.StringValue(eventRuleResp.Score),
		EventTags:   makeListStringAttribute(eventRuleResp.EventTags),
		BuilderConfig: BuilderConfig{
			Filters:       types.StringValue(string(filtersJSON) + "\n"),
			TableName:     types.StringValue(eventRuleResp.BuilderConfig.TableName),
			Added:         types.BoolValue(eventRuleResp.BuilderConfig.Added),
			MatchesFilter: types.BoolValue(eventRuleResp.BuilderConfig.MatchesFilter),
			Severity:      types.StringValue(eventRuleResp.BuilderConfig.Severity),
			Key:           types.StringValue(eventRuleResp.BuilderConfig.Key),
			ValueField:    types.StringValue(eventRuleResp.BuilderConfig.ValueField),
			AutoAlertConfig: AutoAlertConfig{
				RaiseAlert:      types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert),
				DisableAlert:    types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.DisableAlert),
				MetadataSources: types.StringValue(string(metadataJSON) + "\n"),
			},
		},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *eventRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EventRule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventRuleID := state.ID.ValueString()

	_, err := r.client.DeleteEventRule(uptycs.EventRule{
		ID: eventRuleID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting",
			"Could not delete eventRule with ID  "+eventRuleID+": "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r *eventRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
