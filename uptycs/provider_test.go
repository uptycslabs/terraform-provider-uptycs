package uptycs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func makeProviderFactoryMap(name string, prov *provider) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		name: providerserver.NewProtocol6WithError(prov),
	}
}

func TestUptycs(t *testing.T) {
	const testConfig = // language=hcl
	`
provider "uptycs" {
}

resource "uptycs_alert_rule" "test" {
  name        = "marcus test 2"
  description = "marcus test"
  enabled     = true
  grouping    = "MITRE"
  grouping_l2 = "Impact"
  grouping_l3 = "T1560"
  sql_config = {
    interval_seconds : 3600,
  }
  code = "test_marc"
  type = "sql"
  rule = "select * from processes limit 2 :to;"
}

data "uptycs_destination" "test_destination" {
  id = "4c0dee1f-c19a-45fe-bf5d-fd031d6f694f"
}

output "email" {
  value = data.uptycs_destination.test_destination.name
}

resource "uptycs_destination" "test" {
  name    = "marc test"
  address = "marcus.young@foo.com"
  type    = "email"
}

data "uptycs_event_exclude_profile" "sample" {
  id = "2a86d4ad-3aa3-42f1-8430-6da238c82b11"
}
  
resource "uptycs_event_exclude_profile" "sample" {
  name = "marc test"
  description = "a test"
  priority = 9999
  platform = "all"
  metadata = <<EOT
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

resource "uptycs_event_rule" "Access_Key_Created" {
  name        = "Access Key Created"
  description = "Access Key created by an IAM user for another user using CreateAccessKey policy."
  enabled     = false
  grouping    = "ATTACK"
  grouping_l2 = "Privilege Escalation"
  grouping_l3 = "T1078"
  code        = "AWS_THREAT_PRIV_ESC_1"
  type        = "builder"
  rule        = "builder"
  event_tags = [
    "ATTACK",
    "AWS",
  ]
  builder_config = {
    table_name     = "upt_cloud_trail_events"
    added          = true
    matches_filter = true
    severity       = "low"
    key            = "upt_tenant_id"
    value_field    = "user_identity_user_name"
    auto_alert_config = {
      raise_alert   = true
      disable_alert = false
    }
    filters = <<EOT
{
  "and": [
    {
      "name": "event_name",
      "value": "CreateAccessKey",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "not": true,
      "name": "user_identity_type",
      "value": "Root",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "name": "upt_connector_type",
      "value": "aws",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "not": true,
      "name": "user_identity_account_id",
      "value": "921884229492",
      "operator": "EQUALS",
      "caseInsensitive": true
    },
    {
      "not": true,
      "name": "user_identity_user_name",
      "value": "@logdna.com",
      "operator": "CONTAINS",
      "caseInsensitive": true
    }
  ]
}
EOT
  }
}

resource "uptycs_user" "new_user" {
  name  = "someone"
  email = "some+test@foo.com"
}

data "uptycs_user" "by_name" {
  name = "Marcus Young"
}

data "uptycs_user" "by_id" {
  id = "f48f4c40-9c4a-47bb-9e3f-797d4deca92a"
}
 
output "email_by_name" {
  value = data.uptycs_user.by_name.email
}

output "email_by_id" {
  value = data.uptycs_user.by_id.email
}
`

	var (
		prov = new(provider)
	)

	resource.Test(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: makeProviderFactoryMap("uptycs", prov),
			Steps: []resource.TestStep{
				{
					Config: testConfig,
				},
			},
		},
	)
}
