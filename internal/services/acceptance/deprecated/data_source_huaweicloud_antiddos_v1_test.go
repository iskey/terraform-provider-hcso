package deprecated

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
)

func TestAccAntiDdosV1DataSource_basic(t *testing.T) {
	resourceName := "data.hcso_antiddos.antiddos"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosV1DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
		},
	})
}

func testAccCheckAntiDdosV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find defense status of EIP data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Defense status of EIP data source ID not set")
		}

		return nil
	}
}

const testAccAntiDdosV1DataSource_basic = `
resource "hcso_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
  }
}

data "hcso_antiddos" "antiddos" {
  floating_ip_id = hcso_vpc_eip.eip_1.id
}
`
