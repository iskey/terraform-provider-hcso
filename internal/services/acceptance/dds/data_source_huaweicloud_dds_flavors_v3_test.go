package dds

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDDSFlavorV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.hcso_dds_flavors.flavor"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSFlavorV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "mongos"),
				),
			},
		},
	})
}

var testAccDDSFlavorV3DataSource_basic = `
data "hcso_dds_flavors" "flavor" {
  engine_name = "DDS-Community"
  vcpus       = 2
  memory      = 4
  type        = "mongos"
}
`
