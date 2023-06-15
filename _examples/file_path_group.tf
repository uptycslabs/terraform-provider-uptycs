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

data "uptycs_file_path_group" "fpg" {
  id = "c9dee8cc-3931-47e4-a9ca-f5b251ab44c5"
}

output "fpg" {
  value = data.uptycs_file_path_group.fpg
}

resource "uptycs_file_path_group" "new_fpg" {
  check_signature         = false
  custom                  = true
  description             = "marc testtt"
  exclude_paths           = []
  exclude_process_names   = []
  file_accesses           = true
  include_path_extensions = []
  include_paths = [
    "/tmp/%",
    "/private/tmp/%",
  ]
  name           = "marc test"
  priority_paths = []
  signatures     = []
  yara_group_rules = [
    {
      id = "c6655aac-abfd-42d4-b2bc-b0a59e98057a"
    },
    {
      name = "AmazonAccessKeyId"
    },
  ]
}

output "new_fpg" {
  value = uptycs_file_path_group.new_fpg
}
