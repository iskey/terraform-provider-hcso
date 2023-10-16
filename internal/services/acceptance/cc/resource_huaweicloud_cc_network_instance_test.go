package cc

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

func getNetworkInstanceResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	// getNetworkInstance: Query the Network instance
	var (
		getNetworkInstanceHttpUrl = "v3/{domain_id}/ccaas/network-instances/{id}"
		getNetworkInstanceProduct = "cc"
	)
	getNetworkInstanceClient, err := config.NewServiceClient(getNetworkInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating NetworkInstance Client: %s", err)
	}

	getNetworkInstancePath := getNetworkInstanceClient.Endpoint + getNetworkInstanceHttpUrl
	getNetworkInstancePath = strings.Replace(getNetworkInstancePath, "{domain_id}", config.DomainID, -1)
	getNetworkInstancePath = strings.Replace(getNetworkInstancePath, "{id}", state.Primary.ID, -1)

	getNetworkInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getNetworkInstanceResp, err := getNetworkInstanceClient.Request("GET", getNetworkInstancePath, &getNetworkInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving NetworkInstance: %s", err)
	}
	return utils.FlattenResponse(getNetworkInstanceResp)
}

func TestAccNetworkInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_cc_network_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getNetworkInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testNetworkInstance_basic(name, 1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "vpc"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "hcso_vpc.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "region_id", "hcso_vpc.test.0", "region"),
					resource.TestCheckResourceAttrPair(rName, "cloud_connection_id", "hcso_cc_connection.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "demo_description"),
				),
			},
			{
				Config: testNetworkInstance_basic_update(name, 1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "vpc"),
					resource.TestCheckResourceAttr(rName, "description", ""),
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

func TestAccNetworkInstance_multiple(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_cc_network_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getNetworkInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testNetworkInstance_multiple(name, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttr("hcso_cc_network_instance.test.0", "type", "vpc"),
					resource.TestCheckResourceAttr("hcso_cc_network_instance.test.1", "type", "vpc"),
					resource.TestCheckResourceAttr("hcso_cc_network_instance.test.0", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("hcso_cc_network_instance.test.1", "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet("hcso_cc_network_instance.test.0", "domain_id"),
					resource.TestCheckResourceAttrSet("hcso_cc_network_instance.test.1", "domain_id"),
					resource.TestCheckResourceAttrPair("hcso_cc_network_instance.test.0", "instance_id",
						"hcso_vpc.test.0", "id"),
					resource.TestCheckResourceAttrPair("hcso_cc_network_instance.test.1", "instance_id",
						"hcso_vpc.test.1", "id"),
					resource.TestCheckResourceAttrPair("hcso_cc_network_instance.test.0", "region_id",
						"hcso_vpc.test.0", "region"),
					resource.TestCheckResourceAttrPair("hcso_cc_network_instance.test.1", "region_id",
						"hcso_vpc.test.1", "region"),
					resource.TestCheckResourceAttrPair("hcso_cc_network_instance.test.0", "cloud_connection_id",
						"hcso_cc_connection.test", "id"),
					resource.TestCheckResourceAttrPair("hcso_cc_network_instance.test.1", "cloud_connection_id",
						"hcso_cc_connection.test", "id"),
				),
			},
		},
	})
}

func testNetworkInstanceRef(name string, count int) string {
	return fmt.Sprintf(`
resource "hcso_vpc" "test" {
  count = %[1]d

  name = "%[2]s_${count.index}"
  cidr = cidrsubnet("10.12.0.0/16", 4, count.index)
}

resource "hcso_vpc_subnet" "test" {
  count = %[1]d

  name = "%[2]s_${count.index}"
  vpc_id     = hcso_vpc.test[count.index].id
  cidr       = cidrsubnet(hcso_vpc.test[count.index].cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(hcso_vpc.test[count.index].cidr, 4, 1), 1)
}

resource "hcso_cc_connection" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
  description           = "accDemo"
}
`, count, name)
}

func testNetworkInstance_basic(name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cc_network_instance" "test" {
  type                = "vpc"
  cloud_connection_id = hcso_cc_connection.test.id
  instance_id         = try(hcso_vpc.test[0].id, "")
  project_id          = "%[2]s"
  region_id           = try(hcso_vpc.test[0].region, "")
  description         = "demo_description"

  cidrs = [
    try(hcso_vpc_subnet.test[0].cidr, ""),
  ]
}
`, testNetworkInstanceRef(name, count), acceptance.HCSO_PROJECT_ID)
}

func testNetworkInstance_basic_update(name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cc_network_instance" "test" {
  type                = "vpc"
  cloud_connection_id = hcso_cc_connection.test.id
  instance_id         = try(hcso_vpc.test[0].id, "")
  project_id          = "%[2]s"
  region_id           = try(hcso_vpc.test[0].region, "")

  cidrs = [
    try(hcso_vpc_subnet.test[0].cidr, ""),
  ]
}
`, testNetworkInstanceRef(name, count), acceptance.HCSO_PROJECT_ID)
}

func testNetworkInstance_multiple(name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_cc_network_instance" "test" {
  count = %[2]d

  type                = "vpc"
  cloud_connection_id = hcso_cc_connection.test.id
  instance_id         = hcso_vpc.test[count.index].id
  project_id          = "%[3]s"
  region_id           = hcso_vpc.test[count.index].region

  cidrs = [
    hcso_vpc_subnet.test[count.index].cidr,
  ]
}
`, testNetworkInstanceRef(name, count), count, acceptance.HCSO_PROJECT_ID)
}
