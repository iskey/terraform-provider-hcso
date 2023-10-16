package cce

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCCEClustersDataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_cce_clusters.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClustersDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.status", "Available"),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCCEClustersDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "hcso_cce_clusters" "test" {
  name = hcso_cce_cluster.test.name

  depends_on = [hcso_cce_cluster.test]
}
`, testAccCceCluster_config(rName))
}
