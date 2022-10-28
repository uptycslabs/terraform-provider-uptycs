terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.15"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_tag" "tag_by_id" {
  id = "9471a94f-86ca-48eb-9f7f-5b3c0038d006"
}

data "uptycs_tag" "tag_by_key_val" {
  key   = "asset-group"
  value = "enrolling"
}

output "tag_by_id_id" {
  value = data.uptycs_tag.tag_by_id.id
}

output "tag_by_key_val_id" {
  value = data.uptycs_tag.tag_by_id.id
}

resource "uptycs_tag" "new_tag" {
  key   = "sometest"
  value = "marc"
  file_path_groups = [
    "FIM - Canary Baseline",
  ]
  audit_configurations   = []
  event_exclude_profiles = []
  querypacks             = []
  registry_paths         = []
  yara_group_rules       = []
}

output "new_tag_id" {
  value = uptycs_tag.new_tag
}
