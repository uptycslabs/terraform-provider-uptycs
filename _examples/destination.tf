terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.10"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_destination" "foo" {
  id = "4c0dee1f-c19a-45fe-bf5d-fd031d6f694f"
}

resource "uptycs_destination" "test" {
  name    = "marc test"
  address = "marcus.young@foo.com"
  type    = "email"
}

output "foo" {
  value = data.uptycs_destination.foo.name
}
