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

type TagRule struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Query          types.String `tfsdk:"query"`
	Source         types.String `tfsdk:"source"`
	RunOnce        types.Bool   `tfsdk:"run_once"`
	Interval       int          `tfsdk:"interval"`
	OSqueryVersion types.String `tfsdk:"osquery_version"`
	Platform       types.String `tfsdk:"platform"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	System         types.Bool   `tfsdk:"system"`
	ResourceType   types.String `tfsdk:"resource_type"`
}

type TagConfiguration Tag

type Tag struct {
	ID                          types.String             `tfsdk:"id"`
	Value                       types.String             `tfsdk:"value"`
	Key                         types.String             `tfsdk:"key"`
	FlagProfileID               types.String             `tfsdk:"flag_profile_id"`
	CustomProfileID             types.String             `tfsdk:"custom_profile_id"`
	ComplianceProfileID         types.String             `tfsdk:"compliance_profile_id"`
	ProcessBlockRuleID          types.String             `tfsdk:"process_block_rule_id"`
	DNSBlockRuleID              types.String             `tfsdk:"dns_block_rule_id"`
	WindowsDefenderPreferenceID types.String             `tfsdk:"windows_defender_preference_id"`
	Tag                         types.String             `tfsdk:"tag"`
	Custom                      types.Bool               `tfsdk:"custom"`
	System                      types.Bool               `tfsdk:"system"`
	TagRuleID                   types.String             `tfsdk:"tag_rule_id"`
	Status                      types.String             `tfsdk:"status"`
	Source                      types.String             `tfsdk:"source"`
	ResourceType                types.String             `tfsdk:"resource_type"`
	FilePathGroups              []TagConfigurationObject `tfsdk:"file_path_groups"`
	EventExcludeProfiles        []TagConfigurationObject `tfsdk:"event_exclude_profiles"`
	RegistryPaths               []TagConfigurationObject `tfsdk:"registry_paths"`
	Querypacks                  []TagConfigurationObject `tfsdk:"querypacks"`
	YaraGroupRules              []TagConfigurationObject `tfsdk:"yara_group_rules"`
	AuditConfigurations         []TagConfigurationObject `tfsdk:"audit_configurations"`
	//ImageLoadExclusions # TODO: cant find any examples of this
	//AuditGroups         # TODO: cant find any examples of this
	//Destinations        # TODO: cant find any examples of this
	//Redactions          # TODO: cant find any examples of this
	//AuditRules          # TODO: cant find any examples of this
	//PrometheusTargets   # TODO: cant find any examples of this
	//AtcQueries          # TODO: cant find any examples of this
}

type TagConfigurationObjectDetails struct {
	ID                   types.String `tfsdk:"id"`
	AuditConfigurationID types.String `tfsdk:"audit_configuration_id"`
	YaraGroupRuleID      types.String `tfsdk:"yara_group_rule_id"`
	QuerypackID          types.String `tfsdk:"querypack_id"`
	RegistryPathID       types.String `tfsdk:"registry_path_id"`
	EventExcludeProfile  types.String `tfsdk:"event_exclude_profile"`
	FilePathGroupID      types.String `tfsdk:"file_path_group_id"`
	TagID                types.String `tfsdk:"tag_id"`
}

type TagConfigurationObject struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	//AuditConfigurationTag  *TagConfigurationObjectDetails `tfsdk:"AuditConfigurationTag"`
	//YaraGroupRuleTag       *TagConfigurationObjectDetails `tfsdk:"YaraGroupRuleTag"`
	//QuerypackTag           *TagConfigurationObjectDetails `tfsdk:"QuerypackTag"`
	//RegistryPathTag        *TagConfigurationObjectDetails `tfsdk:"RegistryPathTag"`
	//EventExcludeProfileTag *TagConfigurationObjectDetails `tfsdk:"EventExcludeProfileTag"`
	//FilePathGroupTag       *TagConfigurationObjectDetails `tfsdk:"FilePathGroupTag"`
}
