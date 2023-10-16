package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getELBPoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	elbClient, err := cfg.ElbV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return pools.Get(elbClient, state.Primary.ID).Extract()
}

func TestAccElbV3Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"hcso_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.type", "APP_COOKIE"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.cookie_name", "testCookie"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"hcso_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason",
						"test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.type", "APP_COOKIE"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.cookie_name", "testCookie"),
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

func TestAccElbV3Pool_basic_with_loadbalancer(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_loadbalancer(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id",
						"hcso_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_loadbalancer(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id",
						"hcso_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"hcso_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
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

func TestAccElbV3Pool_basic_with_listener(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_listener(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"hcso_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_listener(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"hcso_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
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

func TestAccElbV3Pool_basic_with_type_ip(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_type_ip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_type_ip(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
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

func testAccElbV3PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name      = "%s"
  protocol  = "HTTP"
  lb_method = "ROUND_ROBIN"
  type      = "instance"
  vpc_id    = hcso_vpc.test.id

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3PoolConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name      = "%s"
  protocol  = "HTTP"
  lb_method = "LEAST_CONNECTIONS"
  type      = "instance"
  vpc_id    = hcso_vpc.test.id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
`, common.TestVpc(rName), rNameUpdate)
}

func testAccElbV3PoolConfig_basic_with_loadbalancer(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = hcso_elb_loadbalancer.test.id
}
`, testAccElbV3LoadBalancerConfig_basic(rName), rName)
}

func testAccElbV3PoolConfig_update_with_loadbalancer(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = hcso_elb_loadbalancer.test.id
  type            = "instance"
  vpc_id          = hcso_vpc.test.id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"
}
`, testAccElbV3LoadBalancerConfig_basic(rName), rNameUpdate)
}

func testAccElbV3PoolConfig_basic_with_listener(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = hcso_elb_listener.test.id
}
`, testAccElbV3ListenerConfig_basic(rName), rName)
}

func testAccElbV3PoolConfig_update_with_listener(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "LEAST_CONNECTIONS"
  listener_id = hcso_elb_listener.test.id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"
}
`, testAccElbV3ListenerConfig_basic(rName), rNameUpdate)
}

func testAccElbV3PoolConfig_basic_with_type_ip(rName string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name      = "%s"
  protocol  = "HTTP"
  lb_method = "ROUND_ROBIN"
  type      = "ip"
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3PoolConfig_update_with_type_ip(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "hcso_elb_pool" "test" {
  name      = "%s"
  protocol  = "HTTP"
  lb_method = "LEAST_CONNECTIONS"
  type      = "ip"

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"
}
`, common.TestVpc(rName), rNameUpdate)
}
