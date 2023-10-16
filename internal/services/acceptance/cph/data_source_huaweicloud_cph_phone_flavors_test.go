package cph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
)

func TestAccDatasourcePhoneFlavors_basic(t *testing.T) {
	rName := "data.hcso_cph_phone_flavors.status_filter"
	statusFilter := acceptance.InitDataSourceCheck("data.hcso_cph_phone_flavors.status_filter")
	typeFilter := acceptance.InitDataSourceCheck("data.hcso_cph_phone_flavors.type_filter")
	vcpusFilter := acceptance.InitDataSourceCheck("data.hcso_cph_phone_flavors.vcpus_filter")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePhoneFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					statusFilter.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.server_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.disk"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.resolution"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.phone_capacity"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.status"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.type"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.extend_spec"),

					typeFilter.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),

					vcpusFilter.CheckResourceExists(),
					resource.TestCheckOutput("vcpus_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePhoneFlavors_basic() string {
	return `
data "hcso_cph_phone_flavors" "status_filter" {
  status = "1"
}
output "status_filter_is_useful" {
  value = !contains([for v in data.hcso_cph_phone_flavors.status_filter.flavors[*].status : v == "1"], "false")
}

data "hcso_cph_phone_flavors" "type_filter" {
  type = "0"
}
output "type_filter_is_useful" {
  value = !contains([for v in data.hcso_cph_phone_flavors.type_filter.flavors[*].type : v == "0"], "false")
}

data "hcso_cph_phone_flavors" "vcpus_filter" {
  vcpus = 4
}
output "vcpus_filter_is_useful" {
  value = !contains([for v in data.hcso_cph_phone_flavors.vcpus_filter.flavors[*].vcpus : v == 4], "false")
}
`
}
