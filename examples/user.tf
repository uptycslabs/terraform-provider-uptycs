terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.4"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_user" "test_user" {
  name = "Marcus Young"
}

resource "uptycs_user" "new_user" {
  name = "someone"
  email = "some+test@foo.com"
  active = false
}

output "email" {
  value = data.uptycs_user.test_user.email
}