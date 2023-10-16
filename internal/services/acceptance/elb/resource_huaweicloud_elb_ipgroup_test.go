package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/ipgroups"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccElbV3IpGroup_basic(t *testing.T) {
	var c ipgroups.IpGroup
	name := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))
	resourceName := "hcso_elb_ipgroup.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3IpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3IpGroupConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3IpGroupExists(resourceName, &c),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
				),
			},
			{
				Config: testAccElbV3IpGroupConfig_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_updated", name)),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test updated"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
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

func TestAccElbV3IpGroup_withEpsId(t *testing.T) {
	var c ipgroups.IpGroup
	name := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))
	resourceName := "hcso_elb_ipgroup.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3IpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3IpGroupConfig_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3IpGroupExists(resourceName, &c),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckElbV3IpGroupDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := cfg.ElbV3Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating ELB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_elb_ipgroup" {
			continue
		}

		_, err := ipgroups.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ipGroup still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3IpGroupExists(
	n string, c *ipgroups.IpGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		elbClient, err := cfg.ElbV3Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating ELB client: %s", err)
		}

		found, err := ipgroups.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ipGroup not found")
		}

		*c = *found

		return nil
	}
}

func testAccElbV3IpGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "hcso_elb_ipgroup" "test"{
  name        = "%s"
  description = "terraform test"

  ip_list {
    ip = "192.168.10.10"
    description = "ECS01"
  }
}
`, name)
}

func testAccElbV3IpGroupConfig_update(name string) string {
	return fmt.Sprintf(`
resource "hcso_elb_ipgroup" "test"{
  name        = "%s_updated"
  description = "terraform test updated"

  ip_list {
    ip = "192.168.10.10"
    description = "ECS01"
  }

  ip_list {
    ip = "192.168.10.11"
    description = "ECS02"
  }
}
`, name)
}

func testAccElbV3IpGroupConfig_withEpsId(name string) string {
	return fmt.Sprintf(`
resource "hcso_elb_ipgroup" "test"{
  name        = "%s"
  description = "terraform test"

  ip_list {
    ip = "192.168.10.10"
    description = "ECS01"
  }

  enterprise_project_id = "%s"
}
`, name, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST)
}
