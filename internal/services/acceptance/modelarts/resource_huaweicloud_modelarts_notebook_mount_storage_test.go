package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/modelarts/v1/notebook"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getNotebookMountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ModelArtsV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	notebookId, storageId, err := modelarts.ParseMountFromId(state.Primary.ID)
	if err != nil {
		return nil, err
	}

	return notebook.GetMount(client, notebookId, storageId)
}

func TestAccResourceNotebookMountStorage_basic(t *testing.T) {
	var instance notebook.MountStorageOpts
	resourceName := "hcso_modelarts_notebook_mount_storage.test"
	name := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNotebookMountResourceFunc,
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
				Config: testAccNotebookMountStorage_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "local_mount_directory", "/data/test/"),
					resource.TestCheckResourceAttr(resourceName, "type", "OBSFS"),
					resource.TestCheckResourceAttr(resourceName, "storage_path", fmt.Sprintf("obs://%s/", obsName)),
					resource.TestCheckResourceAttr(resourceName, "status", "MOUNTED"),
					resource.TestCheckResourceAttrPair(resourceName, "notebook_id",
						"hcso_modelarts_notebook.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNotebookMountStorage_basic(rName, obsName string) string {
	return fmt.Sprintf(`
resource "hcso_obs_bucket" "test" {
  bucket      = "%s"
  acl         = "private"
  multi_az    = false
  parallel_fs = true

  tags = {
    parallel_fs = "true"
  }
}

resource "hcso_modelarts_notebook" "test" {
  name      = "%s"
  flavor_id = "modelarts.vm.cpu.2u"
  image_id  = "e1a07296-22a8-4f05-8bc8-e936c8e54090"
  volume {
    type = "EFS"
  }
}

resource "hcso_modelarts_notebook_mount_storage" "test" {
  notebook_id           = hcso_modelarts_notebook.test.id
  storage_path          = "obs://${hcso_obs_bucket.test.bucket}/"
  local_mount_directory = "/data/test/"
}
`, obsName, rName)
}
