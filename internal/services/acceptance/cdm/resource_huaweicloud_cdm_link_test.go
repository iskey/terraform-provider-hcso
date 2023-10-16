package cdm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/link"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
)

func getCdmLinkResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CdmV11Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDM v1 client, err=%s", err)
	}
	clusterId, linkName, err := cdm.ParseLinkInfoFromId(state.Primary.ID)
	if err != nil {
		return nil, err
	}
	return link.Get(client, clusterId, linkName)
}

// Link to OBS
func TestAccResourceCdmLink_basic(t *testing.T) {
	var obj link.LinkCreateOpts
	resourceName := "hcso_cdm_link.test"
	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmLinkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmLinkResource_basic(name, bucketName, acceptance.HCSO_ACCESS_KEY,
					acceptance.HCSO_SECRET_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "hcso_cdm_cluster.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testAccCdmLinkResource_basic(name, bucketName, ak, sk string) string {
	clusterConfig := testAccCdmCluster_basic(name)

	return fmt.Sprintf(`
%s

resource "hcso_obs_bucket" "test" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}

resource "hcso_cdm_link" "test" {
  name       = "%s"
  connector  = "obs-connector"
  cluster_id = hcso_cdm_cluster.test.id
  enabled    = true

  config = {
    "storageType" = "OBS"
    "server"      = trimprefix(hcso_obs_bucket.test.bucket_domain_name, "${hcso_obs_bucket.test.bucket}.")
    "port"        = "443"
  }

  access_key   = "%s"
  secret_key   = "%s"
}
`, clusterConfig, bucketName, name, ak, sk)
}
