terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.20"
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

resource "uptycs_compliance_profile" "test" {
  name        = "marc test"
  description = ""
  priority    = 1337
}

resource "uptycs_flag_profile" "test" {
  name          = "marc test"
  description   = ""
  priority      = 1337
  resource_type = "asset"
  flags         = <<EOT
{
  "tls_hostname": "foo.example.com"
}
EOT
  os_flags      = <<EOT
{}
EOT
}

resource "uptycs_tag_rule" "test" {
  name        = "marcus test"
  description = "a test tag rule"
  interval    = 3601 # >3601 on realtime
  source      = "realtime"
  platform    = "something" # required if realtime source
  run_once    = false
  query       = "select 'sometest=marc' as tag from os_version where name like 'Ubuntu%';"
}

resource "uptycs_custom_profile" "test" {
  name            = "marc test"
  description     = ""
  priority        = 2
  resource_type   = "asset"
  query_schedules = <<EOT
{
  "processes": 100
}
EOT
}

resource "uptycs_tag" "new_tag" {
  key   = "sometest"
  value = "marc"
  file_path_groups = [
    "FIM - Canary Baseline",
  ]
  custom_profile         = uptycs_custom_profile.test.name
  flag_profile           = uptycs_flag_profile.test.name
  compliance_profile     = uptycs_compliance_profile.test.name
  audit_configurations   = []
  event_exclude_profiles = []
  querypacks             = []
  registry_paths         = []
  yara_group_rules       = []
}
