package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ resource.Resource                = &eventRuleResource{}
	_ resource.ResourceWithConfigure   = &eventRuleResource{}
	_ resource.ResourceWithImportState = &eventRuleResource{}
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

func (r *eventRuleResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Required: true,
			},
			"code": {
				Type:     types.StringType,
				Required: true,
				Computed: false,
			},
			"type": {
				Type:     types.StringType,
				Required: true,
				Computed: false,
			},
			"rule": {
				Type:     types.StringType,
				Required: true,
			},
			"grouping": {
				Type:     types.StringType,
				Required: true,
			},
			"grouping_l2": {
				Type:     types.StringType,
				Optional: true,
			},
			"grouping_l3": {
				Type:     types.StringType,
				Optional: true,
			},
			"enabled": {
				Type:     types.BoolType,
				Optional: true,
			},
			"event_tags": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"builder_config": {
				Required: true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"table_name": {
						Type:     types.StringType,
						Optional: true,
					},
					"added": {
						Type:     types.BoolType,
						Optional: true,
					},
					"matches_filter": {
						Type:     types.BoolType,
						Optional: true,
					},
					"severity": {
						Type:     types.StringType,
						Optional: true,
					},
					"key": {
						Type:     types.StringType,
						Optional: true,
					},
					"value_field": {
						Type:     types.StringType,
						Optional: true,
					},
					"filters": {
						Required: true,
						Type:     types.StringType,
					},
					"auto_alert_config": {
						Required: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"raise_alert": {
								Type:     types.BoolType,
								Required: true,
							},
							"disable_alert": {
								Type:     types.BoolType,
								Required: true,
							},
						}),
					},
				}),
			},
		},
	}, nil
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
		Name:        plan.Name.Value,
		Code:        plan.Code.Value,
		Description: plan.Description.Value,
		Rule:        plan.Rule.Value,
		Type:        plan.Type.Value,
		Enabled:     plan.Enabled.Value,
		Custom:      true,
		Grouping:    plan.Grouping.Value,
		GroupingL2:  plan.GroupingL2.Value,
		GroupingL3:  plan.GroupingL3.Value,
		EventTags:   tags,
		BuilderConfig: uptycs.BuilderConfig{
			Filters:       uptycs.CustomJSONString(plan.BuilderConfig.Filters.Value),
			TableName:     plan.BuilderConfig.TableName.Value,
			Added:         plan.BuilderConfig.Added.Value,
			MatchesFilter: plan.BuilderConfig.MatchesFilter.Value,
			Severity:      plan.BuilderConfig.Severity.Value,
			Key:           plan.BuilderConfig.Key.Value,
			ValueField:    plan.BuilderConfig.ValueField.Value,
			AutoAlertConfig: uptycs.AutoAlertConfig{
				DisableAlert: plan.BuilderConfig.AutoAlertConfig.DisableAlert.Value,
				RaiseAlert:   plan.BuilderConfig.AutoAlertConfig.RaiseAlert.Value,
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

	var result = EventRule{
		ID:          types.String{Value: eventRuleResp.ID},
		Enabled:     types.Bool{Value: eventRuleResp.Enabled},
		Name:        types.String{Value: eventRuleResp.Name},
		Description: types.String{Value: eventRuleResp.Description},
		Code:        types.String{Value: eventRuleResp.Code},
		Type:        types.String{Value: eventRuleResp.Type},
		Rule:        types.String{Value: eventRuleResp.Rule},
		Grouping:    types.String{Value: eventRuleResp.Grouping},
		GroupingL2:  types.String{Value: eventRuleResp.GroupingL2},
		GroupingL3:  types.String{Value: eventRuleResp.GroupingL3},
		EventTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		BuilderConfig: BuilderConfig{
			Filters:       types.String{Value: string([]byte(filtersJSON)) + "\n"},
			TableName:     types.String{Value: eventRuleResp.BuilderConfig.TableName},
			Added:         types.Bool{Value: eventRuleResp.BuilderConfig.Added},
			MatchesFilter: types.Bool{Value: eventRuleResp.BuilderConfig.MatchesFilter},
			Severity:      types.String{Value: eventRuleResp.BuilderConfig.Severity},
			Key:           types.String{Value: eventRuleResp.BuilderConfig.Key},
			ValueField:    types.String{Value: eventRuleResp.BuilderConfig.ValueField},
			AutoAlertConfig: AutoAlertConfig{
				DisableAlert: types.Bool{Value: eventRuleResp.BuilderConfig.AutoAlertConfig.DisableAlert},
				RaiseAlert:   types.Bool{Value: eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert},
			},
		},
	}
	if result.Type.Value == "sql" {
		result.Rule.Value += "\n"
	}

	for _, t := range eventRuleResp.EventTags {
		result.EventTags.Elems = append(result.EventTags.Elems, types.String{Value: t})
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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

	var result = EventRule{
		ID:          types.String{Value: eventRuleResp.ID},
		Enabled:     types.Bool{Value: eventRuleResp.Enabled},
		Name:        types.String{Value: eventRuleResp.Name},
		Description: types.String{Value: eventRuleResp.Description},
		Code:        types.String{Value: eventRuleResp.Code},
		Type:        types.String{Value: eventRuleResp.Type},
		Rule:        types.String{Value: eventRuleResp.Rule},
		Grouping:    types.String{Value: eventRuleResp.Grouping},
		GroupingL2:  types.String{Value: eventRuleResp.GroupingL2},
		GroupingL3:  types.String{Value: eventRuleResp.GroupingL3},
		EventTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		BuilderConfig: BuilderConfig{
			Filters:       types.String{Value: string([]byte(filtersJSON)) + "\n"},
			TableName:     types.String{Value: eventRuleResp.BuilderConfig.TableName},
			Added:         types.Bool{Value: eventRuleResp.BuilderConfig.Added},
			MatchesFilter: types.Bool{Value: eventRuleResp.BuilderConfig.MatchesFilter},
			Severity:      types.String{Value: eventRuleResp.BuilderConfig.Severity},
			Key:           types.String{Value: eventRuleResp.BuilderConfig.Key},
			ValueField:    types.String{Value: eventRuleResp.BuilderConfig.ValueField},
			AutoAlertConfig: AutoAlertConfig{
				DisableAlert: types.Bool{Value: eventRuleResp.BuilderConfig.AutoAlertConfig.DisableAlert},
				RaiseAlert:   types.Bool{Value: eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert},
			},
		},
	}

	if result.Type.Value == "sql" {
		result.Rule.Value += "\n"
	}

	for _, t := range eventRuleResp.EventTags {
		result.EventTags.Elems = append(result.EventTags.Elems, types.String{Value: t})
	}

	diags := resp.State.Set(ctx, result)
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

	eventRuleID := state.ID.Value

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
		Name:        plan.Name.Value,
		Code:        plan.Code.Value,
		Custom:      true,
		Description: plan.Description.Value,
		Rule:        plan.Rule.Value,
		Type:        plan.Type.Value,
		Enabled:     plan.Enabled.Value,
		Grouping:    plan.Grouping.Value,
		GroupingL2:  plan.GroupingL2.Value,
		GroupingL3:  plan.GroupingL3.Value,
		EventTags:   tags,
		BuilderConfig: uptycs.BuilderConfig{
			Filters:       uptycs.CustomJSONString(plan.BuilderConfig.Filters.Value),
			TableName:     plan.BuilderConfig.TableName.Value,
			Added:         plan.BuilderConfig.Added.Value,
			MatchesFilter: plan.BuilderConfig.MatchesFilter.Value,
			Severity:      plan.BuilderConfig.Severity.Value,
			Key:           plan.BuilderConfig.Key.Value,
			ValueField:    plan.BuilderConfig.ValueField.Value,
			AutoAlertConfig: uptycs.AutoAlertConfig{
				DisableAlert: plan.BuilderConfig.AutoAlertConfig.DisableAlert.Value,
				RaiseAlert:   plan.BuilderConfig.AutoAlertConfig.RaiseAlert.Value,
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

	var result = EventRule{
		ID:          types.String{Value: eventRuleResp.ID},
		Enabled:     types.Bool{Value: eventRuleResp.Enabled},
		Name:        types.String{Value: eventRuleResp.Name},
		Description: types.String{Value: eventRuleResp.Description},
		Code:        types.String{Value: eventRuleResp.Code},
		Type:        types.String{Value: eventRuleResp.Type},
		Rule:        types.String{Value: eventRuleResp.Rule},
		Grouping:    types.String{Value: eventRuleResp.Grouping},
		GroupingL2:  types.String{Value: eventRuleResp.GroupingL2},
		GroupingL3:  types.String{Value: eventRuleResp.GroupingL3},
		EventTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		BuilderConfig: BuilderConfig{
			Filters:       types.String{Value: string([]byte(filtersJSON)) + "\n"},
			TableName:     types.String{Value: eventRuleResp.BuilderConfig.TableName},
			Added:         types.Bool{Value: eventRuleResp.BuilderConfig.Added},
			MatchesFilter: types.Bool{Value: eventRuleResp.BuilderConfig.MatchesFilter},
			Severity:      types.String{Value: eventRuleResp.BuilderConfig.Severity},
			Key:           types.String{Value: eventRuleResp.BuilderConfig.Key},
			ValueField:    types.String{Value: eventRuleResp.BuilderConfig.ValueField},
			AutoAlertConfig: AutoAlertConfig{
				DisableAlert: types.Bool{Value: eventRuleResp.BuilderConfig.AutoAlertConfig.DisableAlert},
				RaiseAlert:   types.Bool{Value: eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert},
			},
		},
	}

	for _, t := range eventRuleResp.EventTags {
		result.EventTags.Elems = append(result.EventTags.Elems, types.String{Value: t})
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

	eventRuleID := state.ID.Value

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
