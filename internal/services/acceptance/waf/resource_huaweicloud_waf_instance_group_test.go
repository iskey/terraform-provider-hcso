package waf

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getWafInstanceGroupFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.WafDedicatedV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud WAF dedicated client: %s", err)
	}
	return pools.Get(client, state.Primary.ID)
}

func TestAccWafInstanceGroup_basic(t *testing.T) {
	var group pools.Pool
	resourceName := "hcso_waf_instance_group.group_1"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getWafInstanceGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafInstanceGroup_conf(name, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccWafInstanceGroup_conf(name, name+"_updated"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_updated"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccWafInstanceGroup_conf(baseName, groupName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_waf_dedicated_instance" "instance_1" {
  name               = "%s"
  available_zone     = data.hcso_availability_zones.test.names[0]
  specification_code = "waf.instance.professional"
  ecs_flavor         = data.hcso_compute_flavors.test.ids[0]
  vpc_id             = hcso_vpc.test.id
  subnet_id          = hcso_vpc_subnet.test.id
  
  security_group = [
    hcso_networking_secgroup.test.id
  ]
}

resource "hcso_waf_instance_group" "group_1" {
  name   = "%s"
  vpc_id = hcso_vpc.test.id

  depends_on = [hcso_waf_dedicated_instance.instance_1]
}
`, common.TestBaseComputeResources(baseName), baseName, groupName)
}
