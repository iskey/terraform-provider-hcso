package cloudtable

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cloudtable/v2/clusters"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getClusterResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CloudtableV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CloudTable v2 client: %s", err)
	}
	return clusters.Get(c, state.Primary.ID)
}

func TestAccCloudtableCluster_basic(t *testing.T) {
	var cluster clusters.Cluster
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "hcso_cloudtable_cluster.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&cluster,
		getClusterResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudtableCluster_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.hcso_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "storage_size", "10240"),
					resource.TestCheckResourceAttr(resourceName, "rs_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "hbase_version", "1.0.6"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "hcso_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id", "hcso_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"hcso_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"availability_zone", "network_id",
				},
			},
		},
	})
}

func testAccCloudtableCluster_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_cloudtable_cluster" "test" {
  availability_zone = data.hcso_availability_zones.test.names[0]
  name              = "%s"
  storage_type      = "ULTRAHIGH"
  vpc_id            = hcso_vpc.test.id
  network_id        = hcso_vpc_subnet.test.id
  security_group_id = hcso_networking_secgroup.test.id
  hbase_version     = "1.0.6"
  rs_num            = 4
}
`, common.TestBaseNetwork(rName), rName)
}
