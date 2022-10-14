terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.14"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_tag_rule" "tr" {
  id = "56a72960-1673-418a-ac51-dbead6069742"
}

output "tr" {
  value = data.uptycs_tag_rule.tr
}

resource "uptycs_tag_rule" "new_tr" {
  name        = "marcus test"
  description = "a test tag rule"
  interval    = 3601 # >3601 on realtime
  source      = "realtime"
  platform    = "something" # required if realtime source
  run_once    = false
  query       = "select owner_id as tag, upt_asset_id from upt_cloud_instance_inventory;"
}

output "new_tr" {
  value = uptycs_tag_rule.new_tr
}
