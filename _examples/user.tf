terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.15"
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
  #id = "baeb925d-ea1f-44ab-a92a-cc0a5a985cb9"
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

data "uptycs_user" "test_user" {
  id = "c052635d-3e51-4f02-9b9e-b4ad20df7cbd"
}

output "test_user_name" {
  value = data.uptycs_user.test_user.name
}

output "test_user_roles" {
  value = data.uptycs_user.test_user.roles
}

output "test_user_user_object_groups" {
  value = data.uptycs_user.test_user.user_object_groups
}

resource "uptycs_user" "new_user" {
  name               = "someone"
  email              = "some+test@foo.com"
  phone              = "888-867-5309"
  image_url          = "42"
  max_idle_time_mins = 30
  alert_hidden_columns = [
    "id",
  ]
  active = true
  roles = [
    "user",
  ]
  user_object_groups = [
    "workstations"
  ]
}

output "new_user_id" {
  value = resource.uptycs_user.new_user.id
}
