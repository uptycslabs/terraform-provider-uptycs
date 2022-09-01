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

data "uptycs_tag_rule" "foo" {
  id = "56a72960-1673-418a-ac51-dbead6069742"
}

output "tag_rule_foo_name" {
  value = data.uptycs_tag_rule.foo.name
}
output "tag_rule_foo_query" {
  value = data.uptycs_tag_rule.foo.query
}
