package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/nat/v3/transitips"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getPrivateTransitIpResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v3 client: %s", err)
	}

	return transitips.Get(client, state.Primary.ID)
}

func TestAccPrivateTransitIp_basic(t *testing.T) {
	var (
		obj transitips.TransitIp

		rName1 = "hcso_nat_private_transit_ip.test"
		rName2 = "hcso_nat_private_transit_ip.random_ip_address"
		name   = acceptance.RandomAccResourceNameWithDash()
	)

	rc1 := acceptance.InitResourceCheck(
		rName1,
		&obj,
		getPrivateTransitIpResourceFunc,
	)
	rc2 := acceptance.InitResourceCheck(
		rName2,
		&obj,
		getPrivateTransitIpResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateTransitIp_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName1, "subnet_id",
						"hcso_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.68"),
					resource.TestCheckResourceAttr(rName1, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName1, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName1, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName1, "created_at"),
					resource.TestCheckResourceAttrSet(rName1, "updated_at"),
					resource.TestCheckResourceAttrPair(rName2, "subnet_id",
						"hcso_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(rName2, "ip_address"),
					resource.TestCheckResourceAttr(rName2, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName2, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName2, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName2, "created_at"),
					resource.TestCheckResourceAttrSet(rName2, "updated_at"),
				),
			},
			{
				Config: testAccPrivateTransitIp_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.88"),
					resource.TestCheckResourceAttr(rName1, "tags.foo", "baaar"),
					resource.TestCheckResourceAttr(rName1, "tags.newkey", "value"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "tags.foo", "baaar"),
					resource.TestCheckResourceAttr(rName2, "tags.newkey", "value"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPrivateTransitIp_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_private_transit_ip" "test" {
  subnet_id             = hcso_vpc_subnet.test.id
  ip_address            = "192.168.0.68"
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "hcso_nat_private_transit_ip" "random_ip_address" {
  subnet_id             = hcso_vpc_subnet.test.id
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccPrivateTransitIp_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_private_transit_ip" "test" {
  subnet_id             = hcso_vpc_subnet.test.id
  ip_address            = "192.168.0.88"
  enterprise_project_id = "0"

  tags = {
    foo    = "baaar"
    newkey = "value"
  }
}

resource "hcso_nat_private_transit_ip" "random_ip_address" {
  subnet_id             = hcso_vpc_subnet.test.id
  enterprise_project_id = "0"

  tags = {
    foo    = "baaar"
    newkey = "value"
  }
}
`, common.TestBaseNetwork(name), name)
}
