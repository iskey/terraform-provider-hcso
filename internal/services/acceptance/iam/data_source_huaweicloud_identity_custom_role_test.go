package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccIdentityCustomRoleDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_identity_custom_role.role_1"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCustomRoleDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
				),
			},
		},
	})
}

func testAccIdentityCustomRoleDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "hcso_identity_role" test {
  name        = "%s"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ]
    }
  ]
}
EOF
}

data "hcso_identity_custom_role" "role_1" {
  name = hcso_identity_role.test.name

  depends_on = [
	hcso_identity_role.test
  ]
}
`, rName)
}
