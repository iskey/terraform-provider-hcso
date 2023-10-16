package sdrs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectedinstances"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getReplicationAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HCSO_REGION_NAME
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}

	rgs := strings.Split(state.Primary.ID, "/")
	if len(rgs) != 2 {
		return nil, fmt.Errorf("invalid format specified for replication attach id," +
			" must be <protected_instance_id>/<replication_id>")
	}
	instanceID := rgs[0]
	replicationID := rgs[1]

	instance, err := protectedinstances.Get(client, instanceID).Extract()
	if err != nil {
		return nil, err
	}
	for _, attach := range instance.Attachment {
		if attach.Replication == replicationID {
			// find the target attachment
			return &attach, nil
		}
	}
	return nil, fmt.Errorf("error retrieving SDRS replication attach")
}

func TestAccReplicationAttach_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "hcso_sdrs_replication_attach.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getReplicationAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testReplicationAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "device", "/dev/vdb"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "hcso_sdrs_protected_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "replication_id", "hcso_sdrs_replication_pair.test", "id"),
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

func testReplicationAttach_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_evs_volume" "test" {
  name              = "%[2]s"
  description       = "test volume for sdrs replication pair"
  availability_zone = data.hcso_availability_zones.test.names[0]
  volume_type       = "SSD"
  size              = 100
}

resource "hcso_sdrs_replication_pair" "test" {
  name                 = "%[2]s"
  group_id             = hcso_sdrs_protection_group.test.id
  volume_id            = hcso_evs_volume.test.id
  description          = "test description"
  delete_target_volume = true
}
`, testProtectedInstance_basic(name), name)
}

func testReplicationAttach_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_sdrs_replication_attach" "test" {
  instance_id    = hcso_sdrs_protected_instance.test.id
  replication_id = hcso_sdrs_replication_pair.test.id
  device         = "/dev/vdb"
}
`, testReplicationAttach_base(name))
}
