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

type User struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Email              types.String `tfsdk:"email"`
	Phone              types.String `tfsdk:"phone"`
	Active             types.Bool   `tfsdk:"active"`
	SuperAdmin         types.Bool   `tfsdk:"super_admin"`
	Bot                types.Bool   `tfsdk:"bot"`
	Support            types.Bool   `tfsdk:"support"`
	ImageURL           types.String `tfsdk:"image_url"`
	MaxIdleTimeMins    int          `tfsdk:"max_idle_time_mins"`
	AlertHiddenColumns types.List   `tfsdk:"alert_hidden_columns"`
	Roles              types.List   `tfsdk:"roles"`
	UserObjectGroups   types.List   `tfsdk:"user_object_groups"`
}

type Role struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Permissions          types.List   `tfsdk:"permissions"`
	Custom               types.Bool   `tfsdk:"custom"`
	Hidden               types.Bool   `tfsdk:"hidden"`
	NoMinimalPermissions types.Bool   `tfsdk:"no_minimal_permissions"`
	RoleObjectGroups     types.List   `tfsdk:"role_object_groups"`
}

type ObjectGroup struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Key              types.String `tfsdk:"key"`
	Value            types.String `tfsdk:"value"`
	AssetGroupRuleID types.String `tfsdk:"asset_group_rule_id"`
	ObjectGroupID    types.String `tfsdk:"object_group_id"`
	UserID           types.String `tfsdk:"user_id"`
	RoleID           types.String `tfsdk:"role_id"`
	Description      types.String `tfsdk:"description"`
	Secret           types.String `tfsdk:"secret"`
	ObjectType       types.String `tfsdk:"object_type"`
	Custom           types.Bool   `tfsdk:"custom"`
	RetentionDays    int          `tfsdk:"retention_days"`
	Destinations     types.List   `tfsdk:"destinations"`
}
