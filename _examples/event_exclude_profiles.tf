terraform {
  required_providers {
    uptycs = {
      source  = "uptycslabs/uptycs"
      version = "0.0.17"
    }
  }
}

provider "uptycs" {
  host        = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key     = "2222222222222222222222"
  api_secret  = "234444444444433333333333222222221111111"
}

data "uptycs_event_exclude_profile" "sample" {
  id = "2a86d4ad-3aa3-42f1-8430-6da238c82b11"
}

resource "uptycs_event_exclude_profile" "sample" {
  name        = "marc test"
  description = "a test"
  priority    = 9999
  platform    = "all"
  metadata    = <<EOT
{
  "dns_lookup_events": {},
  "user_events": {},
  "socket_events": {},
  "process_events": {
    "path": [
      "^/Library/Developer/Xcode$"
    ]
  },
  "registry_events": {},
  "process_file_events": {
    "path": [
      "^/Library/Developer/Xcode$",
      "^/Library/Application Support/JAMF$"
    ],
    "executable": [
      "^.*osqueryd\\.exe$|^.*collectguestlogs\\.exe$|^.*MsMpEng\\.exe$"
    ]
  }
}
EOT
}

output "source_name" {
  value = data.uptycs_event_exclude_profile.sample.name
}
output "source_priority" {
  value = data.uptycs_event_exclude_profile.sample.priority
}
output "source_metadata" {
  value = data.uptycs_event_exclude_profile.sample.metadata
}
