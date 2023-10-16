package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourcePools_basic(t *testing.T) {
	rName := "data.hcso_elb_pools.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pools.0.name", name),
					resource.TestCheckResourceAttrPair(rName, "pools.0.id",
						"hcso_elb_pool.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.description",
						"hcso_elb_pool.test", "description"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.protocol",
						"hcso_elb_pool.test", "protocol"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.lb_method",
						"hcso_elb_pool.test", "lb_method"),
					resource.TestCheckResourceAttr(rName, "pools.0.type", "instance"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.vpc_id",
						"hcso_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "pools.0.protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(rName, "pools.0.slow_start_enabled", "false"),
				),
			},
		},
	})
}

func testAccDatasourcePools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_elb_pools" "test" {
  name = "%s"

  depends_on = [
    hcso_elb_pool.test
  ]
}
`, testAccElbV3PoolConfig_basic(name), name)
}
