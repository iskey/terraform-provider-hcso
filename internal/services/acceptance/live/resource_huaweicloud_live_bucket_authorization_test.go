package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccLiveBucketAuthorization_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcso_live_bucket_authorization.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testLiveBucketAuthorization_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "bucket", name),
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

func testLiveBucketAuthorization_basic(name string) string {
	return fmt.Sprintf(`
resource "hcso_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "hcso_live_bucket_authorization" "test" {
  depends_on = [hcso_obs_bucket.test]

  bucket = "%[1]s"
}
`, name)
}
