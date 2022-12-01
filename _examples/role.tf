terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.19"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_role" "admin" {
  name = "admin"
}

output "admin_role_id" {
  value = data.uptycs_role.admin.id
}

output "admin_role_name" {
  value = data.uptycs_role.admin.name
}

output "admin_role_permissions" {
  value = data.uptycs_role.admin.permissions
}

resource "uptycs_role" "test_role" {
  name                   = "test_role"
  no_minimal_permissions = true
  permissions = [
    "ASSET:READ",
    "ALERT_RULE:READ",
    "ALERT:READ",
  ]
  role_object_groups = [
    "workstations"
  ]
}

output "test_role_permissions" {
  value = resource.uptycs_role.test_role.permissions
}
