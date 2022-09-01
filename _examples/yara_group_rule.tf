terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.8"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_yara_group_rule" "test" {
  id = "9a5a3262-ee74-417c-ade0-c1948ec8bc27"
}

resource "uptycs_yara_group_rule" "test" {
  name        = "marc testttt"
  description = "another marc test"
}

output "test" {
  value = data.uptycs_yara_group_rule.test
}

output "new_test" {
  value = resource.uptycs_yara_group_rule.test
}
