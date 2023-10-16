package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/appauths"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getAppAuthFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := appauths.ListOpts{
		InstanceId: state.Primary.Attributes["instance_id"],
		AppId:      state.Primary.Attributes["application_id"],
	}
	resp, err := appauths.ListAuthorized(client, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccAppAuth_basic(t *testing.T) {
	var (
		authApis []appauths.ApiAuthInfo

		rName      = "hcso_apig_application_authorization.test"
		rc         = acceptance.InitResourceCheck(rName, &authApis, getAppAuthFunc)
		baseConfig = testAccAppAuth_base()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppAuth_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccAppAuth_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppAuthImportIdFunc(rName),
			},
		},
	})
}

func testAccAppAuthImportIdFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		resourceId := rs.Primary.ID
		if instanceId == "" || resourceId == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<id>' (the format of resource ID is "+
				"'<env_id>/<application_id>'), but got '%s/%s'", instanceId, resourceId)
		}
		return fmt.Sprintf("%s/%s", instanceId, resourceId), nil
	}
}

func testAccAppAuth_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "hcso_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.hcso_images_image.test.id
  flavor_id          = data.hcso_compute_flavors.test.ids[0]
  security_group_ids = [hcso_networking_secgroup.test.id]
  availability_zone  = data.hcso_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = hcso_vpc_subnet.test.id
  }
}

resource "hcso_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = hcso_vpc.test.id
  subnet_id             = hcso_vpc_subnet.test.id
  security_group_id     = hcso_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.hcso_availability_zones.test.names, 0, 1), null)
}

resource "hcso_apig_group" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
}

resource "hcso_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = hcso_apig_instance.test.id
  port        = 80
  algorithm   = "WRR"
  protocol    = "HTTP"
  path        = "/"
  http_code   = "201"

  members {
    id = hcso_compute_instance.test.id
  }
}

resource "hcso_apig_api" "test" {
  count = 3

  instance_id             = hcso_apig_instance.test.id
  group_id                = hcso_apig_group.test.id
  name                    = "%[2]s_${count.index}"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/${count.index}"
  security_authentication = "APP"
  matching                = "Exact"

  web {
    path             = "/getUserAge/${count.index}"
    vpc_channel_id   = hcso_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }
}

resource "hcso_apig_environment" "test" {
  instance_id = hcso_apig_instance.test.id
  name        = "%[2]s"
}

resource "hcso_apig_api_publishment" "test" {
  count = 3

  instance_id = hcso_apig_instance.test.id
  api_id      = hcso_apig_api.test[count.index].id
  env_id      = hcso_apig_environment.test.id
}

resource "hcso_apig_application" "test" {
  instance_id = hcso_apig_instance.test.id// hcso_apig_instance.test.id
  name        = "%[2]s"
}
`, common.TestBaseComputeResources(name), name)
}

func testAccAppAuth_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_application_authorization" "test" {
  depends_on = [hcso_apig_api_publishment.test]

  instance_id    = hcso_apig_instance.test.id
  application_id = hcso_apig_application.test.id
  env_id         = hcso_apig_environment.test.id
  api_ids        = slice(hcso_apig_api.test[*].id, 0, 2)
}
`, baseConfig)
}

func testAccAppAuth_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_application_authorization" "test" {
  depends_on = [hcso_apig_api_publishment.test]

  instance_id    = hcso_apig_instance.test.id
  application_id = hcso_apig_application.test.id
  env_id         = hcso_apig_environment.test.id
  api_ids        = slice(hcso_apig_api.test[*].id, 1, 3)
}
`, baseConfig)
}
