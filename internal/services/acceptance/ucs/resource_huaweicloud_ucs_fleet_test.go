package ucs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getFleetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getFleet: Query the UCS Fleet detail
	var (
		getFleetHttpUrl = "v1/clustergroups/{id}"
		getFleetProduct = "ucs"
	)
	getFleetClient, err := cfg.NewServiceClient(getFleetProduct, "")
	if err != nil {
		return nil, fmt.Errorf("error creating UCS Client: %s", err)
	}

	getFleetPath := getFleetClient.Endpoint + getFleetHttpUrl
	getFleetPath = strings.ReplaceAll(getFleetPath, "{id}", state.Primary.ID)

	getFleetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getFleetResp, err := getFleetClient.Request("GET", getFleetPath, &getFleetOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Fleet: %s", err)
	}

	getFleetRespBody, err := utils.FlattenResponse(getFleetResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Fleet: %s", err)
	}

	return getFleetRespBody, nil
}

func TestAccFleet_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "hcso_ucs_fleet.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFleetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFleet_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "permissions.0.namespaces.0", "default"),
					resource.TestCheckResourceAttrPair(rName, "permissions.0.policy_ids.0",
						"hcso_ucs_policy.test1", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testFleet_update_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(rName, "permissions.0.namespaces.0", "*"),
					resource.TestCheckResourceAttr(rName, "permissions.1.namespaces.0", "default"),
					resource.TestCheckResourceAttr(rName, "permissions.1.namespaces.1", "kube-system"),
					resource.TestCheckResourceAttrPair(rName, "permissions.0.policy_ids.0",
						"hcso_ucs_policy.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "permissions.1.policy_ids.0",
						"hcso_ucs_policy.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "permissions.1.policy_ids.1",
						"hcso_ucs_policy.test2", "id"),
				),
			},
			{
				Config: testFleet_update_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "permissions.#", "0"),
				),
			},
		},
	})
}

func testFleet_basic(name string) string {
	return fmt.Sprintf(`
resource "hcso_identity_user" "test" {
  count = 1

  name     = "%[1]s-${count.index}"
  password = "Test@12345678"
}
    
resource "hcso_ucs_policy" "test1" {
  name         = "%[1]s"
  iam_user_ids = hcso_identity_user.test[*].id
  type         = "admin"
  description  = "created by terraform"
}

resource "hcso_ucs_fleet" "test" {
  name        = "%[1]s"
  description = "created by terraform"

  permissions {
    namespaces = ["default"]
    policy_ids = [hcso_ucs_policy.test1.id]
  }
}
`, name)
}

func testFleet_update_1(name string) string {
	return fmt.Sprintf(`
resource "hcso_identity_user" "test" {
  count = 2

  name     = "%[1]s-${count.index}"
  password = "Test@12345678"
}
    
resource "hcso_ucs_policy" "test1" {
  name         = "%[1]s-1"
  iam_user_ids = hcso_identity_user.test[*].id
  type         = "admin"
  description  = "created by terraform"
}

resource "hcso_ucs_policy" "test2" {
  name         = "%[1]s-2"
  iam_user_ids = hcso_identity_user.test[*].id
  type         = "custom"
  description  = "created by terraform"
  details {
    operations = ["*"]
    resources  = ["*"]
  }
}

resource "hcso_ucs_fleet" "test" {
  name        = "%[1]s"
  description = "created by terraform update"

  permissions {
    namespaces = ["*"]
    policy_ids = [hcso_ucs_policy.test1.id]
  }

  permissions {
    namespaces = ["default", "kube-system"]
    policy_ids = [hcso_ucs_policy.test1.id, hcso_ucs_policy.test2.id]
  }
}
`, name)
}

func testFleet_update_2(name string) string {
	return fmt.Sprintf(`
resource "hcso_ucs_fleet" "test" {
  name = "%[1]s"
}
`, name)
}
