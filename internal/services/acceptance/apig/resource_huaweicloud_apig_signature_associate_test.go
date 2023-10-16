package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/signs"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getSignatureAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opts := signs.ListBindOpts{
		InstanceId:  state.Primary.Attributes["instance_id"],
		SignatureId: state.Primary.Attributes["signature_id"],
		Limit:       500,
	}
	resp, err := signs.ListBind(c, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccSignatureAssociate_basic(t *testing.T) {
	var (
		apiDetails []signs.SignBindApiInfo

		name   = acceptance.RandomAccResourceName()
		rName1 = "hcso_apig_signature_associate.basic_bind"
		rName2 = "hcso_apig_signature_associate.hmac_bind"
		rName3 = "hcso_apig_signature_associate.aes_bind"

		rc1 = acceptance.InitResourceCheck(rName1, &apiDetails, getSignatureAssociateFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &apiDetails, getSignatureAssociateFunc)
		rc3 = acceptance.InitResourceCheck(rName3, &apiDetails, getSignatureAssociateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignatureAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName1, "instance_id",
						"hcso_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName1, "signature_id",
						"hcso_apig_signature.basic", "id"),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName2, "instance_id",
						"hcso_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName2, "signature_id",
						"hcso_apig_signature.hmac", "id"),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName3, "instance_id",
						"hcso_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName3, "signature_id",
						"hcso_apig_signature.aes", "id"),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				Config: testAccSignatureAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName2),
			},
			{
				ResourceName:      rName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName3),
			},
		},
	})
}

func testAccSignatureAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["signature_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<signature_id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["signature_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["signature_id"]), nil
	}
}

func testAccSignatureAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = hcso_vpc.test.id
  subnet_id             = hcso_vpc_subnet.test.id
  security_group_id     = hcso_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.hcso_availability_zones.test.names, 0, 1), null)
}

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
  instance_id             = hcso_apig_instance.test.id
  group_id                = hcso_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/{user_age}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Success response"
  failure_response        = "Failed response"
  description             = "Created by script"

  request_params {
    name     = "user_age"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }
  
  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  web {
    path             = "/getUserAge/{userAge}"
    vpc_channel_id   = hcso_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }

  web_policy {
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getUserAge/{userAge}"
    timeout          = 30000
    vpc_channel_id   = hcso_apig_vpc_channel.test.id

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }

    conditions {
      source     = "param"
      param_name = "user_age"
      type       = "Equal"
      value      = "28"
    }
  }
}

resource "hcso_apig_environment" "test" {
  count = 6

  name        = "%[2]s_${count.index}"
  instance_id = hcso_apig_instance.test.id
}

resource "hcso_apig_api_publishment" "test" {
  count = 6

  instance_id = hcso_apig_instance.test.id
  api_id      = hcso_apig_api.test.id
  env_id      = hcso_apig_environment.test[count.index].id
}

resource "hcso_apig_signature" "basic" {
  instance_id = hcso_apig_instance.test.id
  name        = "%[2]s_basic"
  type        = "basic"
}

resource "hcso_apig_signature" "hmac" {
  instance_id = hcso_apig_instance.test.id
  name        = "%[2]s_hmac"
  type        = "hmac"
}

resource "hcso_apig_signature" "aes" {
  instance_id = hcso_apig_instance.test.id
  name        = "%[2]s_aes"
  type        = "aes"
  algorithm   = "aes-128-cfb"
}
`, common.TestBaseComputeResources(name), name)
}

func testAccSignatureAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_signature_associate" "basic_bind" {
  instance_id  = hcso_apig_instance.test.id
  signature_id = hcso_apig_signature.basic.id
  publish_ids  = slice(hcso_apig_api_publishment.test[*].publish_id, 0, 2)
}

resource "hcso_apig_signature_associate" "hmac_bind" {
  instance_id  = hcso_apig_instance.test.id
  signature_id = hcso_apig_signature.hmac.id
  publish_ids  = slice(hcso_apig_api_publishment.test[*].publish_id, 2, 4)
}

resource "hcso_apig_signature_associate" "aes_bind" {
  instance_id  = hcso_apig_instance.test.id
  signature_id = hcso_apig_signature.aes.id
  publish_ids  = slice(hcso_apig_api_publishment.test[*].publish_id, 4, 6)
}
`, testAccSignatureAssociate_base(name))
}

func testAccSignatureAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_apig_signature_associate" "basic_bind" {
  instance_id  = hcso_apig_instance.test.id
  signature_id = hcso_apig_signature.basic.id
  publish_ids  = slice(hcso_apig_api_publishment.test[*].publish_id, 1, 3)
}

resource "hcso_apig_signature_associate" "hmac_bind" {
  instance_id  = hcso_apig_instance.test.id
  signature_id = hcso_apig_signature.hmac.id
  publish_ids  = slice(hcso_apig_api_publishment.test[*].publish_id, 3, 5)
}

resource "hcso_apig_signature_associate" "aes_bind" {
  instance_id  = hcso_apig_instance.test.id
  signature_id = hcso_apig_signature.aes.id
  publish_ids  = setunion(slice(hcso_apig_api_publishment.test[*].publish_id, 0, 1),
    slice(hcso_apig_api_publishment.test[*].publish_id, 5, 6))
}
`, testAccSignatureAssociate_base(name))
}
