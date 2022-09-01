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

data "uptycs_audit_configuration" "foo" {
  id = "7d51a844-f28e-4dbf-8831-e4a063e16156"
}

output "foo" {
  value = data.uptycs_audit_configuration.foo.name
}
