package uptycs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptycslabs/uptycs-client-go/uptycs"
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

func (d *eventRuleDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Optional: true},
			"name":        schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"score": schema.StringAttribute{Optional: true,
				Computed: false,
			},
			"code": schema.StringAttribute{Optional: true,
				Computed: false,
			},
			"type": schema.StringAttribute{Optional: true,
				Computed: false,
			},
			"rule":        schema.StringAttribute{Optional: true},
			"grouping":    schema.StringAttribute{Optional: true},
			"grouping_l2": schema.StringAttribute{Optional: true},
			"grouping_l3": schema.StringAttribute{Optional: true},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"event_tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
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
					"filters":        schema.StringAttribute{Optional: true},
					"auto_alert_config": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"raise_alert":      schema.BoolAttribute{Optional: true},
							"disable_alert":    schema.BoolAttribute{Optional: true},
							"metadata_sources": schema.StringAttribute{Optional: true},
						},
					},
				},
			},
		},
	}
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
