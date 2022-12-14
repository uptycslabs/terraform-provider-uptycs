package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
)

var (
	_ datasource.DataSource              = &eventRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &eventRuleDataSource{}
)

func EventRuleDataSource() datasource.DataSource {
	return &eventRuleDataSource{}
}

type eventRuleDataSource struct {
	client *uptycs.Client
}

func (d *eventRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_rule"
}

func (d *eventRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*uptycs.Client)
}

func (d *eventRuleDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"score": {
				Type:     types.StringType,
				Optional: true,
				Computed: false,
			},
			"code": {
				Type:     types.StringType,
				Optional: true,
				Computed: false,
			},
			"type": {
				Type:     types.StringType,
				Optional: true,
				Computed: false,
			},
			"rule": {
				Type:     types.StringType,
				Optional: true,
			},
			"grouping": {
				Type:     types.StringType,
				Optional: true,
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
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.UseStateForUnknown(), boolDefault(true)},
			},
			"event_tags": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"builder_config": {
				Optional: true,
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
						Optional: true,
						Type:     types.StringType,
					},
					"auto_alert_config": {
						Optional: true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"raise_alert": {
								Type:     types.BoolType,
								Optional: true,
							},
							"disable_alert": {
								Type:     types.BoolType,
								Optional: true,
							},
							"metadata_sources": {
								Type:     types.StringType,
								Optional: true,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (d *eventRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var eventRuleID string
	var eventRuleName string

	idAttr := req.Config.GetAttribute(ctx, path.Root("id"), &eventRuleID)
	nameAttr := req.Config.GetAttribute(ctx, path.Root("name"), &eventRuleName)

	var eventRuleToLookup uptycs.EventRule

	if len(eventRuleID) == 0 {
		resp.Diagnostics.Append(nameAttr...)
		eventRuleToLookup = uptycs.EventRule{
			Name: eventRuleName,
		}
	} else {
		resp.Diagnostics.Append(idAttr...)
		eventRuleToLookup = uptycs.EventRule{
			ID: eventRuleID,
		}
	}

	eventRuleResp, err := d.client.GetEventRule(eventRuleToLookup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read.",
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
		EventTags: types.List{
			ElemType: types.StringType,
			Elems:    make([]attr.Value, 0),
		},
		BuilderConfig: BuilderConfig{
			Filters:       types.StringValue(string([]byte(filtersJSON)) + "\n"),
			TableName:     types.StringValue(eventRuleResp.BuilderConfig.TableName),
			Added:         types.BoolValue(eventRuleResp.BuilderConfig.Added),
			MatchesFilter: types.BoolValue(eventRuleResp.BuilderConfig.MatchesFilter),
			Severity:      types.StringValue(eventRuleResp.BuilderConfig.Severity),
			Key:           types.StringValue(eventRuleResp.BuilderConfig.Key),
			ValueField:    types.StringValue(eventRuleResp.BuilderConfig.ValueField),
			AutoAlertConfig: AutoAlertConfig{
				DisableAlert:    types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.DisableAlert),
				RaiseAlert:      types.BoolValue(eventRuleResp.BuilderConfig.AutoAlertConfig.RaiseAlert),
				MetadataSources: types.StringValue(string([]byte(metadataJSON)) + "\n"),
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
