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

data "uptycs_custom_profile" "example" {
  name = "ubuntu"
}

data "uptycs_custom_profile" "test" {
  id = "c6815103-33eb-41e0-bc2f-6a23cc2e1589"
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

output "test" {
  value = data.uptycs_custom_profile.test
}

output "create_test" {
  value = resource.uptycs_custom_profile.test
}
