terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.24"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_flag_profile" "example" {
  name = "ubuntu"
}

data "uptycs_flag_profile" "test" {
  id = "c6815103-33eb-41e0-bc2f-6a23cc2e1589"
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

output "test" {
  value = data.uptycs_flag_profile.test
}

output "create_test" {
  value = resource.uptycs_flag_profile.test
}
