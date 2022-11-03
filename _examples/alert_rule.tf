terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.16"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_alert_rule" "test_rule" {
  id = "ce371e14-9f0d-40ec-887b-abf137af716f"
}

output "test_rule" {
  value = data.uptycs_alert_rule.test_rule
}

resource "uptycs_alert_rule" "test_alert_rule" {
  alert_tags = [
    "ATTACK",
    "AWS",
    "Cloud",
    "IAM",
    "Privilege Escalation",
    "T1078",
  ]
  code        = "AWS_THREAT_PRIV_ESC_1_REDDIT_V2_MARCUS"
  description = "Access Key created by an IAM user for another user using CreateAccessKey policy."
  destinations = [
    {
      close_after_delivery = true
      destination_id       = "8196c225-a78c-4d73-b320-a37562ec7f48"
      notify_every_alert   = true
      severity             = ""
    },
  ]
  enabled         = true
  grouping        = "ATTACK"
  grouping_l2     = "Privilege Escalation"
  grouping_l3     = "T1078"
  is_internal     = false
  lock            = false
  name            = "marcus test"
  notify_count    = 0
  notify_interval = 0
  rule            = "select * from processes WHERE users.upt_time >= :from AND users.upt_time < :to;"
  rule_exceptions = [
    "ce67f12a-91d1-4a79-b0ee-a60501c5990b",
  ]
  sql_config = {
    interval_seconds = 3600
  }
  throttled = false
  type      = "sql"
}

output "test_alert_rule" {
  value = uptycs_alert_rule.test_alert_rule
}
