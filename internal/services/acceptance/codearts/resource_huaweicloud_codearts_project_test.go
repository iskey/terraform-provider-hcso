package codearts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getProjectResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getProject: Query the Project
	var (
		getProjectHttpUrl = "v4/projects/{id}"
		getProjectProduct = "projectman"
	)
	getProjectClient, err := config.NewServiceClient(getProjectProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Project Client: %s", err)
	}

	getProjectPath := getProjectClient.Endpoint + getProjectHttpUrl
	getProjectPath = strings.Replace(getProjectPath, "{id}", state.Primary.ID, -1)

	getProjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getProjectResp, err := getProjectClient.Request("GET", getProjectPath, &getProjectOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Project: %s", err)
	}
	return utils.FlattenResponse(getProjectResp)
}

func TestAccProject_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_codearts_project.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProjectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProject_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "scrum"),
				),
			},
			{
				Config: testProject_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "scrum"),
					resource.TestCheckResourceAttr(rName, "description", "demo_description"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"description",
				},
			},
		},
	})
}

func testProject_basic(name string) string {
	return fmt.Sprintf(`
resource "hcso_codearts_project" "test" {
  name = "%s"
  type = "scrum"
}
`, name)
}

func testProject_basic_update(name string) string {
	return fmt.Sprintf(`
resource "hcso_codearts_project" "test" {
  name = "%s"
  type = "scrum"
  description = "demo_description"
}
`, name)
}
