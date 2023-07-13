terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.25"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_exception" "test" {
  name = "Test Account Exception"
}

output "test" {
  value = data.uptycs_exception.test
}

resource "uptycs_exception" "foo" {
  name        = "marc test"
  description = "marc test"
  table_name  = "aws_cloudtrail_events"
  rule        = <<EOT
{
  "and": [
    {
      "caseInsensitive": true,
      "isDate": false,
      "isVersion": false,
      "isWordMatch": false,
      "name": "account_id",
      "not": false,
      "operator": "EQUALS",
      "value": "1111111111"
    }
  ]
}
EOT
}

output "test_created" {
  value = resource.uptycs_exception.foo
}
