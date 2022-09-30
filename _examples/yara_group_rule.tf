terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.12"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_yara_group_rule" "ygr" {
  id = "9a5a3262-ee74-417c-ade0-c1948ec8bc27"
}

resource "uptycs_yara_group_rule" "new_ygr" {
  name        = "marc testttt"
  description = "another marc test"
}

output "ygr" {
  value = data.uptycs_yara_group_rule.ygr
}

output "new_ygr" {
  value = resource.uptycs_yara_group_rule.new_ygr
}
