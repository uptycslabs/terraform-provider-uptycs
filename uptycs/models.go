package uptycs

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AlertRule struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Code        types.String `tfsdk:"code"`
	Type        types.String `tfsdk:"type"`
	Rule        types.String `tfsdk:"rule"`
	Grouping    types.String `tfsdk:"grouping"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	GroupingL2  types.String `tfsdk:"grouping_l2"`
	GroupingL3  types.String `tfsdk:"grouping_l3"`
	SQLConfig   SQLConfig    `tfsdk:"sql_config"`
}

type SQLConfig struct {
	IntervalSeconds int `tfsdk:"interval_seconds"`
}

type EventRule struct {
	ID            types.String  `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	Description   types.String  `tfsdk:"description"`
	Code          types.String  `tfsdk:"code"`
	Type          types.String  `tfsdk:"type"`
	Rule          types.String  `tfsdk:"rule"`
	Grouping      types.String  `tfsdk:"grouping"`
	GroupingL2    types.String  `tfsdk:"grouping_l2"`
	GroupingL3    types.String  `tfsdk:"grouping_l3"`
	Enabled       types.Bool    `tfsdk:"enabled"`
	EventTags     types.List    `tfsdk:"event_tags"`
	BuilderConfig BuilderConfig `tfsdk:"builder_config"`
}

type BuilderConfig struct {
	TableName       types.String    `tfsdk:"table_name"`
	Added           types.Bool      `tfsdk:"added"`
	MatchesFilter   types.Bool      `tfsdk:"matches_filter"`
	Filters         types.String    `tfsdk:"filters"`
	Severity        types.String    `tfsdk:"severity"`
	Key             types.String    `tfsdk:"key"`
	ValueField      types.String    `tfsdk:"value_field"`
	AutoAlertConfig AutoAlertConfig `tfsdk:"auto_alert_config"`
}

type AutoAlertConfig struct {
	RaiseAlert   types.Bool `tfsdk:"raise_alert"`
	DisableAlert types.Bool `tfsdk:"disable_alert"`
}

type Destination struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Address types.String `tfsdk:"address"`
	//Config TODO
	//"config": {
	//  "sender": null
	//},
	Enabled types.Bool `tfsdk:"enabled"`
	//Template TODO
}

type EventExcludeProfile struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Priority     int          `tfsdk:"priority"`
	Metadata     types.String `tfsdk:"metadata"`
	ResourceType types.String `tfsdk:"resource_type"`
	Platform     types.String `tfsdk:"platform"`
}
