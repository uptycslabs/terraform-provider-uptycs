package uptycs

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AlertRule struct {
	ID                  types.String           `tfsdk:"id"`
	Name                types.String           `tfsdk:"name"`
	Description         types.String           `tfsdk:"description"`
	Code                types.String           `tfsdk:"code"`
	Type                types.String           `tfsdk:"type"`
	Rule                types.String           `tfsdk:"rule"`
	Grouping            types.String           `tfsdk:"grouping"`
	Enabled             types.Bool             `tfsdk:"enabled"`
	Custom              types.Bool             `tfsdk:"custom"`
	Throttled           types.Bool             `tfsdk:"throttled"`
	IsInternal          types.Bool             `tfsdk:"is_internal"`
	AlertTags           types.List             `tfsdk:"alert_tags"`
	GroupingL2          types.String           `tfsdk:"grouping_l2"`
	GroupingL3          types.String           `tfsdk:"grouping_l3"`
	Lock                types.Bool             `tfsdk:"lock"`
	AlertNotifyInterval *int                   `tfsdk:"notify_interval"`
	AlertNotifyCount    *int                   `tfsdk:"notify_count"`
	AlertRuleExceptions types.List             `tfsdk:"rule_exceptions"`
	Destinations        []AlertRuleDestination `tfsdk:"destinations"`
	SQLConfig           *SQLConfig             `tfsdk:"sql_config"`
	//ScriptConfig        *ScriptConfig          `tfsdk:"script_config"` //TODO cant find any examples
}

type ScriptConfig struct {
	ID               types.String `tfsdk:"id"`
	QueryPackID      types.String `tfsdk:"querypack_id"`
	TableName        types.String `tfsdk:"table_name"`
	EventCode        types.String `tfsdk:"event_code"`
	EventMinSeverity types.String `tfsdk:"event_min_severity"`
	Added            types.Bool   `tfsdk:"added"`
}

type AlertRuleDestination struct {
	Severity           types.String `tfsdk:"severity"`
	DestinationID      types.String `tfsdk:"destination_id"`
	NotifyEveryAlert   types.Bool   `tfsdk:"notify_every_alert"`
	CloseAfterDelivery types.Bool   `tfsdk:"close_after_delivery"`
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

type Tag struct {
	ID                          types.String `tfsdk:"id"`
	Value                       types.String `tfsdk:"value"`
	Key                         types.String `tfsdk:"key"`
	FlagProfileID               types.String `tfsdk:"flag_profile_id"`
	CustomProfileID             types.String `tfsdk:"custom_profile_id"`
	ComplianceProfileID         types.String `tfsdk:"compliance_profile_id"`
	ProcessBlockRuleID          types.String `tfsdk:"process_block_rule_id"`
	DNSBlockRuleID              types.String `tfsdk:"dns_block_rule_id"`
	WindowsDefenderPreferenceID types.String `tfsdk:"windows_defender_preference_id"`
	Tag                         types.String `tfsdk:"tag"`
	Custom                      types.Bool   `tfsdk:"custom"`
	System                      types.Bool   `tfsdk:"system"`
	TagRuleID                   types.String `tfsdk:"tag_rule_id"`
	Status                      types.String `tfsdk:"status"`
	Source                      types.String `tfsdk:"source"`
	ResourceType                types.String `tfsdk:"resource_type"`
	FilePathGroups              types.List   `tfsdk:"file_path_groups"`
	EventExcludeProfiles        types.List   `tfsdk:"event_exclude_profiles"`
	Querypacks                  types.List   `tfsdk:"querypacks"`
	RegistryPaths               types.List   `tfsdk:"registry_paths"`
	YaraGroupRules              types.List   `tfsdk:"yara_group_rules"`
	AuditConfigurations         types.List   `tfsdk:"audit_configurations"`
	//ImageLoadExclusions # TODO: cant find any examples of this
	//AuditGroups         # TODO: cant find any examples of this
	//Destinations        # TODO: cant find any examples of this
	//Redactions          # TODO: cant find any examples of this
	//AuditRules          # TODO: cant find any examples of this
	//PrometheusTargets   # TODO: cant find any examples of this
	//AtcQueries          # TODO: cant find any examples of this
}

type FilePathGroup struct {
	ID                    types.String             `tfsdk:"id"`
	Name                  types.String             `tfsdk:"name"`
	Description           types.String             `tfsdk:"description"`
	Grouping              types.String             `tfsdk:"grouping"`
	IncludePaths          types.List               `tfsdk:"include_paths"`
	IncludePathExtensions types.List               `tfsdk:"include_path_extensions"`
	ExcludePaths          types.List               `tfsdk:"exclude_paths"`
	Custom                types.Bool               `tfsdk:"custom"`
	CheckSignature        types.Bool               `tfsdk:"check_signature"`
	FileAccesses          types.Bool               `tfsdk:"file_accesses"`
	ExcludeProcessNames   types.List               `tfsdk:"exclude_process_names"`
	PriorityPaths         types.List               `tfsdk:"priority_paths"`
	Signatures            []FilePathGroupSignature `tfsdk:"signatures"`
	YaraGroupRules        []YaraGroupRule          `tfsdk:"yara_group_rules"`
	//ExcludeProcessPaths   []string                 `tfsdk:"excludeProcessPaths"` //TODO this seems broken in the API. returns null or {}
}

type FilePathGroupSignature struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Paths       types.List   `tfsdk:"paths"`
}

type YaraGroupRule struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Rules       types.String `tfsdk:"rules"`
	Custom      types.Bool   `tfsdk:"custom"`
}

type RegistryPath struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Grouping             types.String `tfsdk:"grouping"`
	IncludeRegistryPaths types.List   `tfsdk:"include_registry_paths"`
	RegAccesses          types.Bool   `tfsdk:"reg_accesses"`
	ExcludeRegistryPaths types.List   `tfsdk:"exclude_registry_paths"`
	Custom               types.Bool   `tfsdk:"custom"`
}

type Querypack struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Type             types.String `tfsdk:"type"`
	AdditionalLogger types.Bool   `tfsdk:"additional_logger"`
	Custom           types.Bool   `tfsdk:"custom"`
	IsInternal       types.Bool   `tfsdk:"is_internal"`
	ResourceType     types.String `tfsdk:"resource_type"`
	Conf             types.String `tfsdk:"conf"`
}

type Query struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Query       types.String `tfsdk:"query"`
	Removed     types.Bool   `tfsdk:"removed"`
	Version     types.String `tfsdk:"version"`
	Interval    int          `tfsdk:"interval"`
	Platform    types.String `tfsdk:"platform"`
	Snapshot    types.Bool   `tfsdk:"snapshot"`
	RunNow      types.Bool   `tfsdk:"run_now"`
	Value       types.String `tfsdk:"value"`
	QuerypackID types.String `tfsdk:"querypack_id"`
	TableName   types.String `tfsdk:"table_name"`
	DataTypes   types.String `tfsdk:"data_types"`
	Verified    types.Bool   `tfsdk:"verified"`
}

type AuditConfiguration struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Framework   types.String `tfsdk:"framework"`
	Version     types.String `tfsdk:"version"`
	OsVersion   types.String `tfsdk:"os_version"`
	Platform    types.String `tfsdk:"platform"`
	TableName   types.String `tfsdk:"table_name"`
	Sha256      types.String `tfsdk:"sha256"`
	Type        types.String `tfsdk:"type"`
	Checks      int          `tfsdk:"checks"`
}

type AssetGroupRule struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Query          types.String `tfsdk:"query"`
	Interval       int          `tfsdk:"interval"`
	OsqueryVersion types.String `tfsdk:"osquery_version"`
	Platform       types.String `tfsdk:"platform"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

type PathStruct struct {
	Path types.String `tfsdk:"path"`
}

type AtcQuery struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Query       types.String `tfsdk:"query"`
}
