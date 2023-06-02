terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.23"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_lookup_table" "test" {
  id = "385d7735-9342-41bc-b660-87040313b39e"
}

output "lookup_table" {
  value = data.uptycs_lookup_table.test
}
