package uptycs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func makeProviderFactoryMap(name string, prov *UptycsProvider) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		name: providerserver.NewProtocol6WithError(prov),
	}
}

func TestUptycs(t *testing.T) {
	const testConfig = // language=hcl
	`
provider "uptycs" {
  host = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key = "2222222222222222222222"
  api_secret = "234444444444433333333333222222221111111"
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
`

	var (
		prov = new(UptycsProvider)
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
