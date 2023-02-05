package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/myoung34/terraform-plugin-framework-utils/modifiers"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
	"time"
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
					modifiers.DefaultString(""),
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
			"enabled": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					modifiers.DefaultBool(false),
				},
			},
			"event_tags": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"sql_config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"interval_seconds": schema.Int64Attribute{
						Optional: true,
						PlanModifiers: []planmodifier.Int64{
							modifiers.DefaultInt(600),
						},
					},
				},
			},
			"builder_config": schema.SingleNestedAttribute{
				Optional: true,
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
							"raise_alert": schema.BoolAttribute{
								Required: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
									modifiers.DefaultBool(false),
								},
							},
							"disable_alert": schema.BoolAttribute{
								Required: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
									modifiers.DefaultBool(false),
								},
							},
							"metadata_sources": schema.StringAttribute{Optional: true},
						},
					},
				},
			},
			"alert_rule": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"rule_exceptions": schema.ListAttribute{
						ElementType: types.StringType,
						Required:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"destinations": schema.ListNestedAttribute{
						Required: true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"severity":             schema.StringAttribute{Optional: true},
								"destination_id":       schema.StringAttribute{Optional: true},
								"notify_every_alert":   schema.BoolAttribute{Optional: true},
								"close_after_delivery": schema.BoolAttribute{Optional: true},
							},
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
	}
	if eventRuleResp.Type == "builder" {
		filtersJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.Filters, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		metadataJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.AutoAlertConfig.MetadataSources, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		result.BuilderConfig = &BuilderConfig{
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
		}

		result.AlertRule = &AlertRuleLite{
			Enabled:             types.BoolValue(false),
			AlertRuleExceptions: makeListStringAttributeFn([]string{}, func(v string) (string, bool) { return v, true }),
			Destinations:        []AlertRuleDestination{},
		}
		if eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert {
			alertRuleResp, err := r.client.GetAlertRule(uptycs.AlertRule{ID: eventRuleResp.ID})
			if err == nil {

				nonGlobalRuleExceptions := make([]attr.Value, 0)
				for _, re := range alertRuleResp.AlertRuleExceptions {
					_ruleException, _ := r.client.GetException(uptycs.Exception{
						ID: re.ExceptionID,
					})
					if !_ruleException.IsGlobal {
						nonGlobalRuleExceptions = append(nonGlobalRuleExceptions, types.StringValue(_ruleException.ID))
					}
				}
				result.AlertRule.AlertRuleExceptions = types.ListValueMust(types.StringType, nonGlobalRuleExceptions)

				_destinations := make([]AlertRuleDestination, 0)
				for _, d := range alertRuleResp.Destinations {
					_destinations = append(_destinations, AlertRuleDestination{
						Severity:           types.StringValue(d.Severity),
						DestinationID:      types.StringValue(d.DestinationID),
						NotifyEveryAlert:   types.BoolValue(d.NotifyEveryAlert),
						CloseAfterDelivery: types.BoolValue(d.CloseAfterDelivery),
					})
				}

				result.AlertRule.Destinations = _destinations
			}
		}
	}

	if eventRuleResp.Type == "sql" {
		result.Rule = types.StringValue(fmt.Sprintf("%s\n", result.Rule.ValueString()))
		result.SQLConfig = &SQLConfig{
			IntervalSeconds: types.Int64Value(int64(eventRuleResp.SQLConfig.IntervalSeconds)),
		}
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

	var eventRuleResp = uptycs.EventRule{}
	var result = EventRule{}

	if plan.Type.ValueString() == "builder" {

		eventRuleResp, _ = r.client.CreateEventRule(uptycs.EventRule{
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
			BuilderConfig: &uptycs.BuilderConfig{
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

		filtersJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.Filters, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		metadataJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.AutoAlertConfig.MetadataSources, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		result = EventRule{
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
			BuilderConfig: &BuilderConfig{
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
			AlertRule: &AlertRuleLite{
				Enabled:             types.BoolValue(false),
				AlertRuleExceptions: makeListStringAttributeFn([]string{}, func(v string) (string, bool) { return v, true }),
				Destinations:        []AlertRuleDestination{},
			},
		}

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not create eventRule, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
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
			SQLConfig: &uptycs.SQLConfig{
				IntervalSeconds: int(plan.SQLConfig.IntervalSeconds.ValueInt64()),
			},
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not create eventRule, unexpected error: "+err.Error(),
			)
			return
		}
		result = EventRule{
			ID:          types.StringValue(eventRuleResp.ID),
			Enabled:     types.BoolValue(eventRuleResp.Enabled),
			Name:        types.StringValue(eventRuleResp.Name),
			Description: types.StringValue(eventRuleResp.Description),
			Code:        types.StringValue(eventRuleResp.Code),
			Type:        types.StringValue(eventRuleResp.Type),
			Grouping:    types.StringValue(eventRuleResp.Grouping),
			GroupingL2:  types.StringValue(eventRuleResp.GroupingL2),
			GroupingL3:  types.StringValue(eventRuleResp.GroupingL3),
			Score:       types.StringValue(eventRuleResp.Score),
			EventTags:   makeListStringAttribute(eventRuleResp.EventTags),
			Rule:        types.StringValue(eventRuleResp.Rule),
			SQLConfig: &SQLConfig{
				IntervalSeconds: types.Int64Value(int64(eventRuleResp.SQLConfig.IntervalSeconds)),
			},
		}
	}

	time.Sleep(1 * time.Second)

	if eventRuleResp.Type == "builder" {
		if !plan.BuilderConfig.AutoAlertConfig.RaiseAlert.ValueBool() {
			if len(plan.AlertRule.Destinations) > 0 || len(plan.AlertRule.AlertRuleExceptions.Elements()) > 0 {
				_, _ = r.client.DeleteEventRule(uptycs.EventRule{ID: eventRuleResp.ID})
				resp.Diagnostics.AddError(
					"Error creating",
					"alert_rule.destinations and alert_rule.rule_exceptions must be empty when builder_config.auto_alert_config.raise_alert is false",
				)
				return
			}

		} else {

			alertRuleResp, err := r.client.GetAlertRule(uptycs.AlertRule{ID: eventRuleResp.ID})
			if err == nil {
				if plan.AlertRule == nil {
					resp.Diagnostics.AddError(
						"Error creating",
						"Could not manage attached alertRule. alert_rule = { destinations = [], rule_exceptions = [] } is required at a minimum for event rules with type 'builder'",
					)
					_, _ = r.client.DeleteEventRule(uptycs.EventRule{ID: eventRuleResp.ID})
					return
				}

				var ruleExceptions []string
				plan.AlertRule.AlertRuleExceptions.ElementsAs(ctx, &ruleExceptions, false)

				// Gather all the rule exceptions from the plan
				// Note: this excludes global rule exceptions so you must gather those back at Update() time
				_ruleExceptions := make([]uptycs.RuleException, 0)
				for _, _re := range ruleExceptions {
					_ruleExceptions = append(_ruleExceptions, uptycs.RuleException{
						ExceptionID: _re,
					})
				}
				// Gather back the global rule exceptions so we dont remove them at Update time by leaving
				// them out of .AlertRuleExceptions[]
				for _, _re := range eventRuleResp.Exceptions {
					re, _ := r.client.GetException(uptycs.Exception{
						ID: _re.ExceptionID,
					})
					if re.IsGlobal {
						_ruleExceptions = append(_ruleExceptions, uptycs.RuleException{
							ExceptionID: re.ID,
						})
					}
				}

				_destinations := make([]uptycs.AlertRuleDestination, 0)
				for _, d := range plan.AlertRule.Destinations {
					_destinations = append(_destinations, uptycs.AlertRuleDestination{
						Severity:           d.Severity.ValueString(),
						DestinationID:      d.DestinationID.ValueString(),
						NotifyEveryAlert:   d.NotifyEveryAlert.ValueBool(),
						CloseAfterDelivery: d.CloseAfterDelivery.ValueBool(),
					})
				}

				alertRule := uptycs.AlertRule{
					ID:                  alertRuleResp.ID,
					AlertRuleExceptions: _ruleExceptions,
					Destinations:        _destinations,
					Type:                eventRuleResp.Type,
					BuilderConfig: &uptycs.BuilderConfigLite{
						ID: eventRuleResp.BuilderConfig.ID,
					},
				}
				if len(alertRule.AlertTags) == 0 {
					alertRule.AlertTags = append(alertRule.AlertTags, eventRuleResp.EventTags...)
				}

				_, err := r.client.UpdateAlertRule(alertRule)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error creating",
						"Could not update attached alertRule, unexpected error: "+err.Error(),
					)
					return
				}

				nonGlobalRuleExceptions := make([]attr.Value, 0)
				for _, re := range alertRuleResp.AlertRuleExceptions {
					_ruleException, _ := r.client.GetException(uptycs.Exception{
						ID: re.ExceptionID,
					})
					if !_ruleException.IsGlobal {
						nonGlobalRuleExceptions = append(nonGlobalRuleExceptions, types.StringValue(_ruleException.ID))
					}
				}
				result.AlertRule.AlertRuleExceptions = types.ListValueMust(types.StringType, nonGlobalRuleExceptions)
				result.AlertRule.Destinations = plan.AlertRule.Destinations
			}
		}
	} else {
		if plan.BuilderConfig != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"builder_config should not be set for event_rules with type 'sql'",
			)
			return
		}
		if plan.AlertRule != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"alert_rule should not be set for event_rules with type 'sql'",
			)
			return
		}
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

	var eventRuleResp = uptycs.EventRule{}
	var result = EventRule{}

	if plan.Type.ValueString() == "builder" {

		eventRuleResp, _ = r.client.UpdateEventRule(uptycs.EventRule{
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
			BuilderConfig: &uptycs.BuilderConfig{
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
		filtersJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.Filters, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		metadataJSON, err := json.MarshalIndent(eventRuleResp.BuilderConfig.AutoAlertConfig.MetadataSources, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		result = EventRule{
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
			BuilderConfig: &BuilderConfig{
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
			AlertRule: &AlertRuleLite{
				Enabled:             types.BoolValue(false),
				AlertRuleExceptions: makeListStringAttributeFn([]string{}, func(v string) (string, bool) { return v, true }),
				Destinations:        []AlertRuleDestination{},
			},
		}

		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating",
				"Could not update eventRule, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
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
			SQLConfig: &uptycs.SQLConfig{
				IntervalSeconds: int(plan.SQLConfig.IntervalSeconds.ValueInt64()),
			},
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating",
				"Could not update eventRule, unexpected error: "+err.Error(),
			)
			return
		}
		result = EventRule{
			ID:          types.StringValue(eventRuleResp.ID),
			Enabled:     types.BoolValue(eventRuleResp.Enabled),
			Name:        types.StringValue(eventRuleResp.Name),
			Description: types.StringValue(eventRuleResp.Description),
			Code:        types.StringValue(eventRuleResp.Code),
			Type:        types.StringValue(eventRuleResp.Type),
			Grouping:    types.StringValue(eventRuleResp.Grouping),
			GroupingL2:  types.StringValue(eventRuleResp.GroupingL2),
			GroupingL3:  types.StringValue(eventRuleResp.GroupingL3),
			Score:       types.StringValue(eventRuleResp.Score),
			EventTags:   makeListStringAttribute(eventRuleResp.EventTags),
			Rule:        types.StringValue(eventRuleResp.Rule),
			SQLConfig: &SQLConfig{
				IntervalSeconds: types.Int64Value(int64(eventRuleResp.SQLConfig.IntervalSeconds)),
			},
		}
	}

	// Do the attached alert rule
	if eventRuleResp.Type == "builder" {
		if !plan.BuilderConfig.AutoAlertConfig.RaiseAlert.ValueBool() {
			if len(plan.AlertRule.Destinations) > 0 || len(plan.AlertRule.AlertRuleExceptions.Elements()) > 0 {
				resp.Diagnostics.AddError(
					"Error creating",
					"alert_rule.destinations and alert_rule.rule_exceptions must be empty when builder_config.auto_alert_config.raise_alert is false",
				)
				return
			}
		} else {
			var ruleExceptions []string
			if plan.AlertRule != nil && len(plan.AlertRule.AlertRuleExceptions.Elements()) > 0 {
				plan.AlertRule.AlertRuleExceptions.ElementsAs(ctx, &ruleExceptions, true)
			}
			_ruleExceptions := make([]uptycs.RuleException, 0)

			// Gather all the rule exceptions from the plan
			// Note: this excludes global rule exceptions so you must gather those back at Update() time
			for _, _re := range ruleExceptions {
				_ruleExceptions = append(_ruleExceptions, uptycs.RuleException{
					ExceptionID: _re,
				})
			}
			// Gather back the global rule exceptions so we dont remove them at Update time by leaving
			// them out of .AlertRuleExceptions[]
			for _, _re := range eventRuleResp.Exceptions {
				re, _ := r.client.GetException(uptycs.Exception{
					ID: _re.ExceptionID,
				})
				if re.IsGlobal {
					_ruleExceptions = append(_ruleExceptions, uptycs.RuleException{
						ExceptionID: re.ID,
					})
				}
			}

			_destinations := make([]uptycs.AlertRuleDestination, 0)
			if plan.AlertRule != nil {
				for _, d := range plan.AlertRule.Destinations {
					_destinations = append(_destinations, uptycs.AlertRuleDestination{
						Severity:           d.Severity.ValueString(),
						DestinationID:      d.DestinationID.ValueString(),
						NotifyEveryAlert:   d.NotifyEveryAlert.ValueBool(),
						CloseAfterDelivery: d.CloseAfterDelivery.ValueBool(),
					})
				}
			}

			alertRule := uptycs.AlertRule{
				ID:                  state.ID.ValueString(),
				AlertRuleExceptions: _ruleExceptions,
				Destinations:        _destinations,
				Type:                eventRuleResp.Type,
				BuilderConfig: &uptycs.BuilderConfigLite{
					ID: eventRuleResp.BuilderConfig.ID,
				},
			}
			if len(alertRule.AlertTags) == 0 {
				alertRule.AlertTags = append(alertRule.AlertTags, eventRuleResp.EventTags...)
			}

			alertRuleResp, err := r.client.UpdateAlertRule(alertRule)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating",
					"Could not update attached alertRule, unexpected error: "+err.Error(),
				)
				return
			}

			var arResult = AlertRuleLite{
				Enabled:             types.BoolValue(alertRuleResp.Enabled),
				AlertRuleExceptions: makeListStringAttributeFn(alertRuleResp.AlertRuleExceptions, func(v uptycs.RuleException) (string, bool) { return v.ExceptionID, true }),
			}

			destinations := make([]AlertRuleDestination, 0)
			for _, d := range alertRuleResp.Destinations {
				destinations = append(destinations, AlertRuleDestination{
					Severity:           types.StringValue(d.Severity),
					DestinationID:      types.StringValue(d.DestinationID),
					NotifyEveryAlert:   types.BoolValue(d.NotifyEveryAlert),
					CloseAfterDelivery: types.BoolValue(d.CloseAfterDelivery),
				})
			}

			arResult.Destinations = destinations
			if plan.AlertRule != nil {
				result.AlertRule = &arResult
			}
		}

	} else {
		if plan.BuilderConfig != nil {
			resp.Diagnostics.AddError(
				"Error updating",
				"builder_config should not be set for event_rules with type 'sql'",
			)
			return
		}
		if plan.AlertRule != nil {
			resp.Diagnostics.AddError(
				"Error updating",
				"alert_rule should not be set for event_rules with type 'sql'",
			)
			return
		}
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
