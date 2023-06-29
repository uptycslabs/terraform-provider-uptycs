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

data "uptycs_destination" "securitymonitoring" {
  name = "#security-monitoring"
}

data "uptycs_destination" "foo" {
  id = "4c0dee1f-c19a-45fe-bf5d-fd031d6f694f"
}

resource "uptycs_destination" "test" {
  name    = "marc test"
  address = "marcus.young@foo.com"
  type    = "email"
}

resource "uptycs_destination" "test2" {
  address = "https://hooks.example.com/hooks/123456"
  enabled = true
  name    = "Testing"
  type    = "http"
  config = {
    sender           = ""
    token            = ""
    password         = ""
    data_key         = ""
    headers          = <<EOT
{}
EOT
    method           = "POST"
    slack_attachment = false
    username         = ""
  }
  template = ""
  lifecycle {
    ignore_changes = [config]
  }
}
