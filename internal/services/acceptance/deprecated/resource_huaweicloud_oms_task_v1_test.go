package deprecated

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/maas/v1/task"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccMaasTask_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckMaas(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMaasTaskV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaasTaskV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMaasTaskV1Exists("hcso_oms_task.task_1"),
					resource.TestCheckResourceAttr("hcso_oms_task.task_1", "description", "migration task"),
				),
			},
		},
	})
}

func testAccCheckMaasTaskV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	maasClient, err := config.MaasV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud maas client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_oms_task" {
			continue
		}

		_, err := task.Get(maasClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Maas task still exists")
		}
	}

	return nil
}

func testAccCheckMaasTaskV1Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		maasClient, err := config.MaasV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud maas client: %s", err)
		}

		found, err := task.Get(maasClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if strconv.FormatInt(found.ID, 10) != rs.Primary.ID {
			return fmtp.Errorf("Task not found")
		}

		return nil
	}
}

var testAccMaasTaskV1_basic = fmt.Sprintf(`
resource "hcso_oms_task" "task_1" {
  description = "migration task"
  enable_kms = false
  thread_num = 1
  src_node {
    region = "cn-beijing"
	ak = "%s"
	sk = "%s"
    object_key = "123.txt"
    bucket = "oms-bucket"
  }
  dst_node {
    region = "%s"
	ak = "%s"
	sk = "%s"
    object_key = "oms"
    bucket = "oms-test"
  }
}
`, acceptance.HCSO_SRC_ACCESS_KEY, acceptance.HCSO_SRC_SECRET_KEY, acceptance.HCSO_REGION_NAME, acceptance.HCSO_ACCESS_KEY, acceptance.HCSO_SECRET_KEY)
