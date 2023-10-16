package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func TestAccRdsInstance_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))
	newPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_basic(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.52"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.52"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				Config: testAccRdsInstance_update(name, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_updated"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.62"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.62"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8636"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", newPwd),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
					"availability_zone",
				},
			},
		},
	})
}

func TestAccRdsInstance_without_password(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_without_password(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
				),
			},
			{
				Config: testAccRdsInstance_without_password_update(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
					"availability_zone",
				},
			},
		},
	})
}

func TestAccRdsInstance_withEpsId(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_epsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccRdsInstance_ha(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_ha(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "async"),
				),
			},
		},
	})
}

func TestAccRdsInstance_mysql(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))
	newPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_mysql_step1(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				Config: testAccRdsInstance_mysql_step2(updateName, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.67"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.67"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3308"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", newPwd),
				),
			},
			{
				Config: testAccRdsInstance_mysql_step3(updateName, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "0"),
				),
			},
		},
	})
}

func TestAccRdsInstance_sqlserver(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_sqlserver(name, pwd, "192.168.0.56"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "collation", "Chinese_PRC_CI_AS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.56"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.56"),
				),
			},
			{
				Config: testAccRdsInstance_sqlserver(name, pwd, "192.168.0.66"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "collation", "Chinese_PRC_CI_AS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.66"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.66"),
				),
			},
		},
	})
}

func TestAccRdsInstance_prePaid(t *testing.T) {
	var (
		instance instances.RdsInstanceResponse

		resourceType = "hcso_rds_instance"
		resourceName = "hcso_rds_instance.test"
		name         = acceptance.RandomAccResourceName()
		password     = fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
			acctest.RandIntRange(10, 99))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_prePaid(name, password, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccRdsInstance_prePaid(name, password, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccRdsInstance_withParameters(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_parameters(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "div_precision_increment"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "12"),
				),
			},
			{
				Config: testAccRdsInstance_newParameters(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "connect_timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "14"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_mysql(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test_backup"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))
	newPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_mysql(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				Config: testAccRdsInstance_restore_mysql_update(name, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3308"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", newPwd),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_sqlserver(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test_backup"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))
	newPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_sqlserver(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				Config: testAccRdsInstance_restore_sqlserver_update(name, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8636"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", newPwd),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_pg(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "hcso_rds_instance"
	resourceName := "hcso_rds_instance.test_backup"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))
	newPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_pg(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8732"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				Config: testAccRdsInstance_restore_pg_update(name, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.hcso_rds_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8733"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", newPwd),
				),
			},
		},
	})
}

func testAccCheckRdsInstanceDestroy(rsType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != rsType {
				continue
			}

			id := rs.Primary.ID
			instance, err := rds.GetRdsInstanceByID(client, id)
			if err != nil {
				return err
			}
			if instance.Id != "" {
				return fmt.Errorf("%s (%s) still exists", rsType, id)
			}
		}
		return nil
	}
}

func testAccCheckRdsInstanceExists(name string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		found, err := rds.GetRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("error checking %s exist, err=%s", name, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", name)
		}

		instance = found
		return nil
	}
}

func testAccRdsInstance_basic(name, password string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name              = "%s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.hcso_availability_zones.test.names[0]]
  security_group_id = hcso_networking_secgroup.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  vpc_id            = hcso_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.52"

  db {
    password = "%s"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, common.TestBaseNetwork(name), name, password)
}

// name, volume.size, backup_strategy, flavor, tags and password will be updated
func testAccRdsInstance_update(name, password string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name              = "%s-update"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.hcso_availability_zones.test.names[0]]
  security_group_id = hcso_networking_secgroup.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  vpc_id            = hcso_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.62"

  db {
    password = "%s"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8636
  }
  volume {
    type = "CLOUDSSD"
    size = 100
  }
  backup_strategy {
    start_time = "09:00-10:00"
    keep_days  = 2
  }

  tags = {
    key1 = "value"
    foo  = "bar_updated"
  }
}
`, common.TestBaseNetwork(name), name, password)
}

func testAccRdsInstance_without_password(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name              = "%s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.hcso_availability_zones.test.names[0]]
  security_group_id = hcso_networking_secgroup.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  vpc_id            = hcso_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "12"
    port    = 8635
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccRdsInstance_without_password_update(name, password string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name              = "%s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.hcso_availability_zones.test.names[0]]
  security_group_id = hcso_networking_secgroup.test.id
  subnet_id         = hcso_vpc_subnet.test.id
  vpc_id            = hcso_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "%s"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, common.TestBaseNetwork(name), name, password)
}

func testAccRdsInstance_epsId(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name                  = "%s"
  flavor                = "rds.pg.n1.large.2"
  availability_zone     = [data.hcso_availability_zones.test.names[0]]
  security_group_id     = hcso_networking_secgroup.test.id
  subnet_id             = hcso_vpc_subnet.test.id
  vpc_id                = hcso_vpc.test.id
  enterprise_project_id = "%s"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
`, common.TestBaseNetwork(name), name, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRdsInstance_ha(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.n1.large.2.ha"
  security_group_id   = hcso_networking_secgroup.test.id
  subnet_id           = hcso_vpc_subnet.test.id
  vpc_id              = hcso_vpc.test.id
  time_zone           = "UTC+08:00"
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "async"
  availability_zone   = [
    data.hcso_availability_zones.test.names[0],
    data.hcso_availability_zones.test.names[1],
  ]

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, common.TestBaseNetwork(name), name)
}

// if the instance flavor has been changed, then a temp instance will be kept for 12 hours,
// the binding relationship between instance and security group or subnet cannot be unbound
// when deleting the instance in this period time, so we cannot create a new vpc, subnet and
// security group in the test case, otherwise, they cannot be deleted when destroy the resource
func testAccRdsInstance_mysql_step1(name, pwd string) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test" {
  name = "default"
}

data "hcso_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "hcso_rds_instance" "test" {
  name              = "%[1]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = true  

  db {
    password = "%[2]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
    period     = 1
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 400
    trigger_threshold = 15
  }
}
`, name, pwd)
}

func testAccRdsInstance_mysql_step2(name, pwd string) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test" {
  name = "default"
}

data "hcso_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "hcso_rds_instance" "test" {
  name              = "%[1]s"
  flavor            = data.hcso_rds_flavors.test.flavors[1].name
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = false
  fixed_ip          = "192.168.0.67"

  db {
    password = "%[2]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  backup_strategy {
    start_time = "18:15-19:15"
    keep_days  = 5
    period     = 3
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 500
    trigger_threshold = 20
  }
}
`, name, pwd)
}

func testAccRdsInstance_mysql_step3(name, pwd string) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test" {
  name = "default"
}

data "hcso_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "hcso_rds_instance" "test" {
  name              = "%[1]s"
  flavor            = data.hcso_rds_flavors.test.flavors[1].name
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = false
  fixed_ip          = "192.168.0.67"

  db {
    password = "%[2]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, name, pwd)
}

func testAccRdsInstance_sqlserver(name, pwd, fixedIp string) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test" {
  name = "default"
}

data "hcso_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2017_EE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "hcso_rds_instance" "test" {
  name              = "%s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"
  fixed_ip          = "%s"

  availability_zone = [
    data.hcso_availability_zones.test.names[0],
  ]

  db {
    password = "%s"
    type     = "SQLServer"
    version  = "2017_EE"
    port     = 8635
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, name, fixedIp, pwd)
}

func testAccRdsInstance_prePaid(name, pwd string, isAutoRenew bool) string {
	return fmt.Sprintf(`
data "hcso_availability_zones" "test" {}

data "hcso_vpc" "test" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test" {
  name = "default"
}

data "hcso_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "hcso_rds_instance" "test" {
  vpc_id            = data.hcso_vpc.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  security_group_id = data.hcso_networking_secgroup.test.id
  
  availability_zone = [
    data.hcso_availability_zones.test.names[0],
  ]

  name      = "%[1]s"
  flavor    = data.hcso_rds_flavors.test.flavors[0].name
  collation = "Chinese_PRC_CI_AS"

  db {
    password = "%[2]s"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8638
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]v"
}
`, name, pwd, isAutoRenew)
}

func testAccRdsInstance_parameters(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.mysql.sld4.large.ha"
  security_group_id   = hcso_networking_secgroup.test.id
  subnet_id           = hcso_vpc_subnet.test.id
  vpc_id              = hcso_vpc.test.id
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.hcso_availability_zones.test.names[0],
    data.hcso_availability_zones.test.names[3],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "LOCALSSD"
    size = 40
  }

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccRdsInstance_newParameters(name string) string {
	return fmt.Sprintf(`
%s

data "hcso_availability_zones" "test" {}

resource "hcso_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.mysql.sld4.large.ha"
  security_group_id   = hcso_networking_secgroup.test.id
  subnet_id           = hcso_vpc_subnet.test.id
  vpc_id              = hcso_vpc.test.id
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.hcso_availability_zones.test.names[0],
    data.hcso_availability_zones.test.names[3],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "LOCALSSD"
    size = 40
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccRdsInstance_restore_mysql(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_vpc" "test_update" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test_update" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test_update" {
  name = "default"
}

resource "hcso_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  security_group_id = data.hcso_networking_secgroup.test_update.id
  subnet_id         = data.hcso_vpc_subnet.test_update.id
  vpc_id            = data.hcso_vpc.test_update.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = true  

  restore {
    instance_id = hcso_rds_backup.test.instance_id
    backup_id   = hcso_rds_backup.test.id
  }

  db {
    password = "%[3]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
    period     = 1
  }

  volume {
    type              = "CLOUDSSD"
    size              = 50
    limit_size        = 400
    trigger_threshold = 15
  }
}
`, testBackup_mysql_basic(name), name, pwd)
}

func testAccRdsInstance_restore_mysql_update(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_vpc" "test_update" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test_update" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test_update" {
  name = "default"
}

resource "hcso_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[1].name
  security_group_id = data.hcso_networking_secgroup.test_update.id
  subnet_id         = data.hcso_vpc_subnet.test_update.id
  vpc_id            = data.hcso_vpc.test_update.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = false

  restore {
    instance_id = hcso_rds_backup.test.instance_id
    backup_id   = hcso_rds_backup.test.id
  }

  db {
    password = "%[3]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  backup_strategy {
    start_time = "18:15-19:15"
    keep_days  = 5
    period     = 3
  }

  volume {
    type              = "CLOUDSSD"
    size              = 60
    limit_size        = 500
    trigger_threshold = 20
  }
}
`, testBackup_mysql_basic(name), name, pwd)
}

func testAccRdsInstance_restore_sqlserver(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = hcso_rds_backup.test.instance_id
    backup_id   = hcso_rds_backup.test.id
  }

  db {
    password = "%[3]s"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8635
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testBackup_sqlserver_basic(name), name, pwd)
}

func testAccRdsInstance_restore_sqlserver_update(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[1].name
  security_group_id = data.hcso_networking_secgroup.test.id
  subnet_id         = data.hcso_vpc_subnet.test.id
  vpc_id            = data.hcso_vpc.test.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = hcso_rds_backup.test.instance_id
    backup_id   = hcso_rds_backup.test.id
  }

  db {
    password = "%[3]s"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8636
  }

  volume {
    type = "CLOUDSSD"
    size = 60
  }
}
`, testBackup_sqlserver_basic(name), name, pwd)
}

func testAccRdsInstance_restore_pg(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_vpc" "test_update" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test_update" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test_update" {
  name = "default"
}

resource "hcso_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[0].name
  security_group_id = data.hcso_networking_secgroup.test_update.id
  subnet_id         = data.hcso_vpc_subnet.test_update.id
  vpc_id            = data.hcso_vpc.test_update.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = hcso_rds_backup.test.instance_id
    backup_id   = hcso_rds_backup.test.id
  }

  db {
    password = "%[3]s"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8732
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testBackup_pg_basic(name), name, pwd)
}

func testAccRdsInstance_restore_pg_update(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "hcso_vpc" "test_update" {
  name = "vpc-default"
}

data "hcso_vpc_subnet" "test_update" {
  name = "subnet-default"
}

data "hcso_networking_secgroup" "test_update" {
  name = "default"
}

resource "hcso_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.hcso_rds_flavors.test.flavors[1].name
  security_group_id = data.hcso_networking_secgroup.test_update.id
  subnet_id         = data.hcso_vpc_subnet.test_update.id
  vpc_id            = data.hcso_vpc.test_update.id
  availability_zone = slice(sort(data.hcso_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = hcso_rds_backup.test.instance_id
    backup_id   = hcso_rds_backup.test.id
  }

  db {
    password = "%[3]s"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8733
  }

  volume {
    type = "CLOUDSSD"
    size = 60
  }
}
`, testBackup_pg_basic(name), name, pwd)
}
