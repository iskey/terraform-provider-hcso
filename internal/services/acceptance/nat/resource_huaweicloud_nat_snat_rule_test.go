package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/nat/v2/snats"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getPublicSnatRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NatGatewayClient(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating NAT v2 client: %s", err)
	}

	return snats.Get(client, state.Primary.ID)
}

func TestAccPublicSnatRule_basic(t *testing.T) {
	var (
		obj snats.Rule

		rName = "hcso_nat_snat_rule.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPublicSnatRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicSnatRule_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "hcso_nat_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "hcso_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "floating_ip_id", "hcso_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccPublicSnatRule_basic_step_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "nat_gateway_id", "hcso_nat_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPublicSnatRule_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_vpc_eip" "test" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = format("%[2]s-%%d", count.index)
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "hcso_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.hcso_images_image.test.id
  flavor_id          = data.hcso_compute_flavors.test.ids[0]
  security_group_ids = [hcso_networking_secgroup.test.id]
  availability_zone  = data.hcso_availability_zones.test.names[0]

  network {
    uuid = hcso_vpc_subnet.test.id
  }
}

resource "hcso_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "2"
  vpc_id                = hcso_vpc.test.id
  subnet_id             = hcso_vpc_subnet.test.id
  enterprise_project_id = "0"
}
`, common.TestBaseComputeResources(name), name)
}

func testAccPublicSnatRule_basic_step_1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_snat_rule" "test" {
  nat_gateway_id = hcso_nat_gateway.test.id
  subnet_id      = hcso_vpc_subnet.test.id
  floating_ip_id = hcso_vpc_eip.test[0].id
  description    = "Created by acc test"
}
`, testAccPublicSnatRule_base(name))
}

func testAccPublicSnatRule_basic_step_2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_nat_snat_rule" "test" {
  nat_gateway_id = hcso_nat_gateway.test.id
  subnet_id      = hcso_vpc_subnet.test.id
  floating_ip_id = join(",", hcso_vpc_eip.test[*].id)
}
`, testAccPublicSnatRule_base(name))
}
