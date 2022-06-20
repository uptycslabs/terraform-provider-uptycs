terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.2"
    }
  }
}

provider "uptycs" {
  host = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key = "2222222222222222222222"
  api_secret = "234444444444433333333333222222221111111"
}

resource "uptycs_alert_rule" "test" {
  name        = "marcus test 2"
  description = "marcus test"
  enabled     = true
  grouping    = "MITRE"
  grouping_l2 = "Impact"
  grouping_l3 = "T1560"
  sql_config = {
    interval_seconds : 3600,
  }
  code = "test_marc"
  type = "sql"
  rule = "select * from processes limit 2 :to;"
}
