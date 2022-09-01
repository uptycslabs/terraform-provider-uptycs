terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.7"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

resource "uptycs_event_rule" "Access_Key_Created" {
  name        = "Access Key Created"
  description = "Access Key created by an IAM user for another user using CreateAccessKey policy."
  enabled     = false
  grouping    = "ATTACK"
  grouping_l2 = "Privilege Escalation"
  grouping_l3 = "T1078"
  code        = "AWS_THREAT_PRIV_ESC_1"
  type        = "builder"
  rule        = "builder"
  event_tags = [
    "ATTACK",
    "AWS",
  ]
  builder_config = {
    table_name     = "upt_cloud_trail_events"
    added          = true
    matches_filter = true
    severity       = "low"
    key            = "upt_tenant_id"
    value_field    = "user_identity_user_name"
    auto_alert_config = {
      raise_alert   = true
      disable_alert = false
    }
    filters = <<EOT
{
  "and": [
    {
      "name": "event_name",
      "value": "CreateAccessKey",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "not": true,
      "name": "user_identity_type",
      "value": "Root",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "name": "upt_connector_type",
      "value": "aws",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "not": true,
      "name": "user_identity_account_id",
      "value": "921884229492",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "not": true,
      "name": "user_identity_user_name",
      "value": "@logdna.com",
      "operator": "CONTAINS",
      "caseInsensitive": true
    }
  ]
}
EOT
  }
}

resource "uptycs_event_rule" "a_different_rule" {
  name        = "Access Key Created"
  description = "Access Key created by an IAM user for another user using CreateAccessKey policy."
  enabled     = false
  grouping    = "ATTACK"
  grouping_l2 = "Privilege Escalation"
  grouping_l3 = "T1078"
  code        = "AWS_THREAT_PRIV_ESC_1"
  type        = "builder"
  rule        = "builder"
  event_tags = [
    "ATTACK",
    "AWS",
    "Cloud",
    "IAM",
    "Privilege Escalation",
    "T1078",
    "THREAT",
    "elastalert"
  ]
  builder_config = <<EOT
{
  "tableName": "process_open_sockets",
  "added": true,
  "matchesFilter": true,
  "filters": {
    "and": [
      {
        "not": true,
        "name": "remote_address",
        "value": "^172.(1[6-9]|2[0-9]|3[01])|^10.|^192.168.",
        "operator": "MATCHES_REGEX"
      }
    ]
  },
  "severity": "low",
  "key": "Test",
  "valueField": "pid",
  "autoAlertConfig": {},
  "addedStr": "true"
}
EOT
}
