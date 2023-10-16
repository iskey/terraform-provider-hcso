package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getGroupFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ApigV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apigroups.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccGroup_basic(t *testing.T) {
	var (
		group apigroups.Group

		rName      = "hcso_apig_group.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&group,
		getGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
				),
			},
			{
				Config: testAccGroup_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupImportStateFunc(),
			},
		},
	})
}

func TestAccGroup_variables(t *testing.T) {
	var (
		group apigroups.Group

		rName = "hcso_apig_group.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&group,
		getGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				// Bind two environment to group, and create some variables.
				Config: testAccGroup_variables(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "environment.#", "2"),
				),
			},
			{
				// Update the variables for two environments.
				Config: testAccGroup_variablesUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "environment.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupImportStateFunc(),
			},
		},
	})
}

func testAccGroupImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "hcso_apig_group.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{id}', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccGroup_base(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = hcso_vpc.test.id
  subnet_id             = hcso_vpc_subnet.test.id
  security_group_id     = hcso_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = [
    data.hcso_availability_zones.test.names[0],
  ]
}
`, common.TestBaseNetwork(name), name)
}

func testAccGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_group" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
  description = "Created by script"
}
`, testAccGroup_base(name), name)
}

func testAccGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_group" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
}
`, testAccGroup_base(name), name)
}

func testAccGroup_variablesBase(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_environment" "test1" {
  name        = "%[2]s_1"
  instance_id = hcso_apig_instance.test.id
  description = "Created by script"
}

resource "hcso_apig_environment" "test2" {
  name        = "%[2]s_2"
  instance_id = hcso_apig_instance.test.id
  description = "Created by script"
}
`, testAccGroup_base(name), name)
}

// Create two environments for the group, and add a total of three variables to the two environments.
// Each of the two environments has a variable with the same name and different value.
func testAccGroup_variables(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_group" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
  description = "Created by script"

  environment {
    environment_id = hcso_apig_environment.test1.id

    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
  }
  environment {
    environment_id = hcso_apig_environment.test2.id

    variable {
      name  = "TERRAFORM"
      value = "/res/terraform"
    }
    variable {
      name  = "DEMO"
      value = "/stage/demo"
    }
  }
}
`, testAccGroup_variablesBase(rName), rName)
}

func testAccGroup_variablesUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_group" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
  description = "Created by script"

  environment {
    environment_id = hcso_apig_environment.test1.id

    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
    variable {
      name  = "TEST"
      value = "/stage/test"
    }
  }
  environment {
    environment_id = hcso_apig_environment.test2.id

    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
  }
}
`, testAccGroup_variablesBase(rName), rName)
}
