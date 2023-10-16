package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/er/v3/propagations"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/er"
)

func getPropagationResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ErV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return er.QueryPropagationById(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["route_table_id"], state.Primary.ID)
}

func TestAccPropagation_basic(t *testing.T) {
	var (
		obj propagations.Propagation

		rName    = "hcso_er_propagation.test"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPropagationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPropagation_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"hcso_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id",
						"hcso_er_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "attachment_id",
						"hcso_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "vpc"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPropagationImportStateFunc(),
			},
		},
	})
}

func testAccPropagationImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId, propagationId string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "hcso_er_propagation" {
				instanceId = rs.Primary.Attributes["instance_id"]
				routeTableId = rs.Primary.Attributes["route_table_id"]
				propagationId = rs.Primary.ID
			}
		}
		if instanceId == "" || routeTableId == "" || propagationId == "" {
			return "", fmt.Errorf("some import IDs are missing, want "+
				"'<instance_id>/<route_table_id>/<propagation_id>', but '%s/%s/%s'",
				instanceId, routeTableId, propagationId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, routeTableId, propagationId), nil
	}
}

func testAccPropagation_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "hcso_vpc_subnet" "test" {
  vpc_id = hcso_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(hcso_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(hcso_vpc.test.cidr, 4, 1), 1)
}

resource "hcso_er_instance" "test" {
  availability_zones = ["%[2]s"]

  name = "%[1]s"
  asn  = %[3]d
}

resource "hcso_er_vpc_attachment" "test" {
  instance_id = hcso_er_instance.test.id
  vpc_id      = hcso_vpc.test.id
  subnet_id   = hcso_vpc_subnet.test.id

  name                   = "%[1]s"
  auto_create_vpc_routes = true
}

resource "hcso_er_route_table" "test" {
  instance_id = hcso_er_instance.test.id

  name = "%[1]s"
}
`, name, acceptance.HCSO_AVAILABILITY_ZONE, bgpAsNum)
}

func testAccPropagation_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_er_propagation" "test" {
  instance_id    = hcso_er_instance.test.id
  route_table_id = hcso_er_route_table.test.id
  attachment_id  = hcso_er_vpc_attachment.test.id
}
`, testAccPropagation_base(name, bgpAsNum))
}
