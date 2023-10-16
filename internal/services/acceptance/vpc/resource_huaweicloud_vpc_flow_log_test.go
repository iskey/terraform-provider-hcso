package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/flowlogs"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getFlowLogResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NetworkingV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC client: %s", err)
	}
	return flowlogs.Get(client, state.Primary.ID).Extract()
}

func TestAccFlowLog_basic(t *testing.T) {
	var flowlog flowlogs.FlowLog
	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "-updated"
	resourceName := "hcso_vpc_flow_log.flow_log"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&flowlog,
		getFlowLogResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowLog_basic(rName, rName, "created by terraform testacc"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform testacc"),
					resource.TestCheckResourceAttr(resourceName, "resource_type", "network"),
					resource.TestCheckResourceAttr(resourceName, "traffic_type", "all"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccFlowLog_basic(rName, rNameUpdate, "updated by terraform testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform testacc"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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

func testAccFlowLogConfigBase(name string) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "vpc_1" {
  name = "vpc-%[1]s"
  cidr = "172.16.0.0/16"
}

resource "hcso_vpc_subnet" "subnet_1" {
  vpc_id     = hcso_vpc.vpc_1.id
  name       = "subnet-%[1]s"
  cidr       = "172.16.0.0/24"
  gateway_ip = "172.16.0.1"
}

resource "hcso_lts_group" "acc_group" {
  group_name  = "%[1]s"
  ttl_in_days = 7
}
resource "hcso_lts_stream" "acc_stream" {
  group_id    = hcso_lts_group.acc_group.id
  stream_name = "%[1]s"
}
`, name)
}

func testAccFlowLog_basic(baseName, resName, resDesc string) string {
	return fmt.Sprintf(`
%s

resource "hcso_vpc_flow_log" "flow_log" {
  name          = "%s"
  description   = "%s"
  resource_type = "network"
  resource_id   = hcso_vpc_subnet.subnet_1.id
  log_group_id  = hcso_lts_group.acc_group.id
  log_stream_id = hcso_lts_stream.acc_stream.id
}
`, testAccFlowLogConfigBase(baseName), resName, resDesc)
}
