terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.9"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_querypack" "qp" {
  id = "4fe27193-c548-460c-b198-d44a3c96050a"
}

output "qp_test" {
  value = data.uptycs_querypack.qp
}

resource "uptycs_querypack" "new_qp" {
  description = "a test"
  name        = "marc_test"
  type        = "vulnerability"
  conf        = <<-EOT
{
  "queries": {
    "linux_baseline": {
      "description": "",
      "query": "SELECT\n    path,\n    filename,\n    symlink\nFROM\n    file\nWHERE\n    (\n        path like '/usr/lib/%%'\n        OR path like '/lib64/%%'\n        OR path like '/bin/%%'\n        OR path like '/sbin/%%'\n        OR path like '/usr/bin/%%'\n        OR path like '/usr/sbin/%%'\n        OR path like '/usr/local/bin/%%'\n        OR path like '/usr/local/sbin/%%'\n    )\n    and filename != '.'",
      "removed": true,
      "version": null,
      "interval": 86400,
      "platform": "linux",
      "snapshot": true,
      "runNow": false,
      "value": ""
    },
    "linux_baseline_lib_directory": {
      "description": "",
      "query": "SELECT\n    path,\n    directory,\n    filename,\n    symlink\nFROM\n    file\nWHERE path like '/lib/%%'\n  and filename != '.'",
      "removed": true,
      "version": null,
      "interval": 86400,
      "platform": "linux",
      "snapshot": true,
      "runNow": false,
      "value": ""
    }
  }
}
EOT
}

output "new_qp" {
  value = resource.uptycs_querypack.new_qp
}
