terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.26"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_registry_path" "rp" {
  id = "ce064913-0c00-4b14-8df3-b1dd90372f04"
}

resource "uptycs_registry_path" "new_rp" {
  name = "marc test"
  include_registry_paths = [
    "/foo/bar/**wut"
  ]
  exclude_registry_paths = []
}

output "rp" {
  value = data.uptycs_registry_path.rp
}

output "new_rp" {
  value = resource.uptycs_registry_path.new_rp
}
