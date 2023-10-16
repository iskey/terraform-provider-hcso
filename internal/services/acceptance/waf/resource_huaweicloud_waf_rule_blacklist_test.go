package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/whiteblackip_rules"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccWafRuleBlackList_basic(t *testing.T) {
	var rule rules.WhiteBlackIP
	randName := acceptance.RandomAccResourceName()
	rName1 := "hcso_waf_rule_blacklist.rule_1"
	rName2 := "hcso_waf_rule_blacklist.rule_2"
	rName3 := "hcso_waf_rule_blacklist.rule_3"
	addressGroupResourceName := "hcso_waf_address_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafRuleBlackListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleBlackList_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(rName1, &rule),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(rName1, "action", "0"),

					testAccCheckWafRuleBlackListExists(rName2, &rule),
					resource.TestCheckResourceAttr(rName2, "ip_address", "192.165.0.0/24"),
					resource.TestCheckResourceAttr(rName2, "action", "1"),

					testAccCheckWafRuleBlackListExists(rName3, &rule),
					resource.TestCheckResourceAttr(rName3, "ip_address", "192.160.0.0/24"),
					resource.TestCheckResourceAttr(rName3, "action", "0"),
					resource.TestCheckResourceAttr(rName3, "name", randName),
					resource.TestCheckResourceAttr(rName3, "description", "test description"),
				),
			},
			{
				Config: testAccWafRuleBlackList_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(rName1, &rule),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.125"),
					resource.TestCheckResourceAttr(rName1, "action", "2"),

					testAccCheckWafRuleBlackListExists(rName2, &rule),
					resource.TestCheckResourceAttr(rName2, "ip_address", "192.150.0.0/24"),
					resource.TestCheckResourceAttr(rName2, "action", "0"),

					testAccCheckWafRuleBlackListExists(rName3, &rule),
					resource.TestCheckResourceAttrPair(rName3, "address_group_id", addressGroupResourceName, "id"),
					resource.TestCheckResourceAttr(rName3, "action", "2"),
					resource.TestCheckResourceAttr(rName3, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(rName3, "description", ""),
					resource.TestCheckResourceAttrSet(rName3, "address_group_name"),
					resource.TestCheckResourceAttrSet(rName3, "address_group_size"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName1),
			},
		},
	})
}

func TestAccWafRuleBlackList_withEpsID(t *testing.T) {
	var rule rules.WhiteBlackIP
	randName := acceptance.RandomAccResourceName()
	rName := "hcso_waf_rule_blacklist.rule"
	addressGroupResourceName := "hcso_waf_address_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafRuleBlackListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleBlackList_basic_withEpsID(randName, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(rName, &rule),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "ip_address", "192.160.0.0/24"),
					resource.TestCheckResourceAttr(rName, "action", "0"),
					resource.TestCheckResourceAttr(rName, "name", randName),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
				),
			},
			{
				Config: testAccWafRuleBlackList_update_withEpsID(randName, acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(rName, &rule),
					resource.TestCheckResourceAttrPair(rName, "address_group_id", addressGroupResourceName, "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HCSO_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "action", "2"),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "address_group_name"),
					resource.TestCheckResourceAttrSet(rName, "address_group_size"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName),
			},
		},
	})
}

func testAccCheckWafRuleBlackListDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	wafClient, err := config.WafV1Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hcso_waf_rule_blacklist" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.GetWithEpsId(wafClient, policyID, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"]).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleBlackListExists(n string, rule *rules.WhiteBlackIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		wafClient, err := config.WafV1Client(acceptance.HCSO_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
		}

		policyID := rs.Primary.Attributes["policy_id"]
		found, err := rules.GetWithEpsId(wafClient, policyID, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"]).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF black list rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafRuleBlackList_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "hcso_waf_rule_blacklist" "rule_1" {
  policy_id  = hcso_waf_policy.policy_1.id
  ip_address = "192.168.0.0/24"
}

resource "hcso_waf_rule_blacklist" "rule_2" {
  policy_id  = hcso_waf_policy.policy_1.id
  ip_address = "192.165.0.0/24"
  action     = 1
}

resource "hcso_waf_rule_blacklist" "rule_3" {
  policy_id   = hcso_waf_policy.policy_1.id
  ip_address  = "192.160.0.0/24"
  name        = "%s"
  description = "test description"
}
`, testAccWafPolicyV1_basic(name), name)
}

func testAccWafRuleBlackList_update(name string) string {
	return fmt.Sprintf(`
%s

resource "hcso_waf_address_group" "test" {
  name         = "%s"
  description  = "example_description"
  ip_addresses = ["192.168.1.0/24"]

  depends_on   = [hcso_waf_dedicated_instance.instance_1]
}

resource "hcso_waf_rule_blacklist" "rule_1" {
  policy_id  = hcso_waf_policy.policy_1.id
  ip_address = "192.168.0.125"
  action     = 2
}

resource "hcso_waf_rule_blacklist" "rule_2" {
  policy_id  = hcso_waf_policy.policy_1.id
  ip_address = "192.150.0.0/24"
  action     = 0
}

resource "hcso_waf_rule_blacklist" "rule_3" {
  policy_id        = hcso_waf_policy.policy_1.id
  address_group_id = hcso_waf_address_group.test.id
  action           = 2
  name             = "%s_update"
  description      = ""
}
`, testAccWafPolicyV1_basic(name), name, name)
}

func testAccWafRuleBlackList_basic_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "hcso_waf_rule_blacklist" "rule" {
  policy_id             = hcso_waf_policy.policy_1.id
  ip_address            = "192.160.0.0/24"
  name                  = "%s"
  description           = "test description"
  enterprise_project_id = "%s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), name, epsID)
}

func testAccWafRuleBlackList_update_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_waf_address_group" "test" {
  name                  = "%[2]s"
  description           = "example_description"
  ip_addresses          = ["192.168.1.0/24"]
  enterprise_project_id = "%[3]s"

  depends_on   = [hcso_waf_dedicated_instance.instance_1]
}

resource "hcso_waf_rule_blacklist" "rule" {
  policy_id             = hcso_waf_policy.policy_1.id
  address_group_id      = hcso_waf_address_group.test.id
  action                = 2
  name                  = "%[2]s_update"
  description           = ""
  enterprise_project_id = "%[3]s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), name, epsID)
}
