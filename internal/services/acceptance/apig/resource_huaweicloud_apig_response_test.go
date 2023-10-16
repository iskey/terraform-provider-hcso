package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/responses"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getResponseFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return responses.Get(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["group_id"],
		state.Primary.ID).Extract()
}

func TestAccResponse_basic(t *testing.T) {
	var (
		resp responses.Response

		rName       = "hcso_apig_response.test"
		name        = acceptance.RandomAccResourceName()
		updateName  = acceptance.RandomAccResourceName()
		basicConfig = testAccResponse_base(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&resp,
		getResponseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResponse_basic(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				Config: testAccResponse_basic(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResponseImportStateFunc(),
			},
		},
	})
}

func TestAccResponse_customRules(t *testing.T) {
	var (
		resp responses.Response

		rName       = "hcso_apig_response.test"
		name        = acceptance.RandomAccResourceName()
		basicConfig = testAccResponse_base(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&resp,
		getResponseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResponse_basic(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				// Add two custom rule.
				Config: testAccResponse_rules(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule.#", "2"),
				),
			},
			{
				// Remove one and update another.
				Config: testAccResponse_rulesUpdate(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResponseImportStateFunc(),
			},
		},
	})
}

func testAccResponseImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "hcso_apig_response.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{group_id}/{name}', but '%s/%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["group_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["name"]), nil
	}
}

func testAccResponse_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_availability_zones" "test" {}

resource "hcso_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = hcso_vpc.test.id
  subnet_id             = hcso_vpc_subnet.test.id
  security_group_id     = hcso_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = [
    data.hcso_availability_zones.test.names[0],
  ]
}

resource "hcso_apig_group" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testAccResponse_basic(relatedConfig string, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_response" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
  group_id    = hcso_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, relatedConfig, name)
}

func testAccResponse_rules(relatedConfig string, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_response" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
  group_id    = hcso_apig_group.test.id

  rule {
    error_type  = "ACCESS_DENIED"
    body        = "{\"error_code\":\"$context.error.code\",\"error_msg\":\"$context.error.message\"}"
    status_code = 400
  }
  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, relatedConfig, name)
}

func testAccResponse_rulesUpdate(relatedConfig string, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_response" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
  group_id    = hcso_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 403
  }
}
`, relatedConfig, name)
}
