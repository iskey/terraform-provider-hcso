package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourceIdentityCenterUsers_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.hcso_identitycenter_users.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterUsers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.0.id"),
					resource.TestCheckResourceAttrSet(rName, "users.0.user_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.family_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.given_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.display_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.email"),
					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
					resource.TestCheckOutput("family_name_filter_is_useful", "true"),
					resource.TestCheckOutput("given_name_filter_is_useful", "true"),
					resource.TestCheckOutput("display_name_filter_is_useful", "true"),
					resource.TestCheckOutput("email_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenterUsers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_identitycenter_users" "test" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  user_name         = hcso_identitycenter_user.test.user_name
}

data "hcso_identitycenter_users" "user_name_filter" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  user_name         = hcso_identitycenter_user.test.user_name
}

data "hcso_identitycenter_users" "family_name_filter" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  family_name       = hcso_identitycenter_user.test.family_name
}

data "hcso_identitycenter_users" "given_name_filter" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  given_name        = hcso_identitycenter_user.test.given_name
}

data "hcso_identitycenter_users" "display_name_filter" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  display_name      = hcso_identitycenter_user.test.display_name
}

data "hcso_identitycenter_users" "email_filter" {
  identity_store_id = data.hcso_identitycenter_instance.test.identity_store_id
  email             = hcso_identitycenter_user.test.email
}

locals {
  user_name_filter_result = [for v in data.hcso_identitycenter_users.user_name_filter.users[*].user_name:
  v == data.hcso_identitycenter_users.test.users.0.user_name]
  family_name_filter_result = [for v in data.hcso_identitycenter_users.family_name_filter.users[*].family_name:
  v == data.hcso_identitycenter_users.test.users.0.family_name]
  given_name_filter_result = [for v in data.hcso_identitycenter_users.given_name_filter.users[*].given_name:
  v == data.hcso_identitycenter_users.test.users.0.given_name]
  display_name_filter_result = [for v in data.hcso_identitycenter_users.display_name_filter.users[*].display_name:
  v == data.hcso_identitycenter_users.test.users.0.display_name]
  email_filter_filter_result = [for v in data.hcso_identitycenter_users.email_filter.users[*].email:
  v == data.hcso_identitycenter_users.test.users.0.email]
}

output "user_name_filter_is_useful" {
  value = alltrue(local.user_name_filter_result) && length(local.user_name_filter_result) > 0
}

output "family_name_filter_is_useful" {
  value = alltrue(local.family_name_filter_result) && length(local.family_name_filter_result) > 0
}

output "given_name_filter_is_useful" {
  value = alltrue(local.given_name_filter_result) && length(local.given_name_filter_result) > 0
}

output "display_name_filter_is_useful" {
  value = alltrue(local.display_name_filter_result) && length(local.display_name_filter_result) > 0
}

output "email_filter_is_useful" {
  value = alltrue(local.email_filter_filter_result) && length(local.email_filter_filter_result) > 0
}
`, testIdentityCenterUser_basic(name))
}
