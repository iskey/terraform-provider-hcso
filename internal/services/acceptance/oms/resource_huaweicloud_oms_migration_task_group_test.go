package oms

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	oms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"

	"github.com/huaweicloud/terraform-provider-hcso/internal/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getMigrationTaskGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.HcOmsV2Client(acceptance.HCSO_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating OMS client: %s", err)
	}

	return c.ShowTaskGroup(&oms.ShowTaskGroupRequest{GroupId: state.Primary.ID})
}

func TestAccMigrationTaskGroup_prefix(t *testing.T) {
	var taskGroup oms.ShowTaskGroupResponse
	rName := acceptance.RandomAccResourceNameWithDash()
	sourceBucketName := rName + "-source"
	destBucketName := rName + "-dest"
	resourceName := "hcso_oms_migration_task_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&taskGroup,
		getMigrationTaskGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOmsInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMigrationTaskGroup_prefix(testMigrationTask_base(sourceBucketName, destBucketName)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "source_object.0.data_source", "HuaweiCloud"),
					resource.TestCheckResourceAttrPair(resourceName, "source_object.0.bucket",
						"hcso_obs_bucket.source", "bucket"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_object.0.bucket",
						"hcso_obs_bucket.dest", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "action", "stop"),
					resource.TestCheckResourceAttr(resourceName, "type", "PREFIX"),
					resource.TestCheckResourceAttr(resourceName, "enable_kms", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "test task group"),
					resource.TestCheckResourceAttr(resourceName, "migrate_since", "2023-01-02 15:04:05"),
					resource.TestCheckResourceAttr(resourceName, "object_overwrite_mode", "CRC64_COMPARISON_OVERWRITE"),
					resource.TestCheckResourceAttr(resourceName, "consistency_check", "crc64"),
					resource.TestCheckResourceAttr(resourceName, "enable_requester_pays", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_failed_object_recording", "true"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.max_bandwidth", "1"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.start", "15:00"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.end", "16:00"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.max_bandwidth", "2"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.start", "16:00"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.end", "17:00"),
				),
			},
			{
				Config: testMigrationTaskGroup_prefix_update(testMigrationTask_base(sourceBucketName, destBucketName)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "action", "start"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.max_bandwidth", "2"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.start", "15:00"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.0.end", "16:00"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.max_bandwidth", "3"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.start", "16:00"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_policy.1.end", "17:00"),
				),
			},
		},
	})
}

func TestAccMigrationTaskGroup_list(t *testing.T) {
	var taskGroup oms.ShowTaskGroupResponse
	rName := acceptance.RandomAccResourceNameWithDash()
	sourceBucketName := rName + "-source"
	destBucketName := rName + "-dest"
	listFileBucketName := rName + "-list"
	resourceName := "hcso_oms_migration_task_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&taskGroup,
		getMigrationTaskGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOmsInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMigrationTaskGroup_list(testMigrationTask_base(sourceBucketName, destBucketName), listFileBucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "LIST"),
					resource.TestCheckResourceAttr(resourceName, "description", "test task group"),
					resource.TestCheckResourceAttr(resourceName, "action", "stop"),
				),
			},
		},
	})
}

func TestAccMigrationTaskGroup_urlList(t *testing.T) {
	var taskGroup oms.ShowTaskGroupResponse
	rName := acceptance.RandomAccResourceNameWithDash()
	destBucketName := rName + "-dest"
	listFileBucketName := rName + "-list"
	resourceName := "hcso_oms_migration_task_group.test"

	tmpFile, err := os.Create("temp.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tmpFile.Close(); err != nil {
			t.Fatal(err)
		}

		err := os.Remove("temp.txt")
		if err != nil {
			t.Fatal(err)
		}
	}()

	url := fmt.Sprintf("https://%s.obs.%s.myhuaweicloud.com/folder/temp.txt", listFileBucketName, acceptance.HCSO_REGION_NAME)
	data := fmt.Sprintf("%s	%s", url, "folder/temp.txt")
	err = os.WriteFile(tmpFile.Name(), []byte(data), 0600)
	if err != nil {
		t.Fatal(err)
	}

	rc := acceptance.InitResourceCheck(
		resourceName,
		&taskGroup,
		getMigrationTaskGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOmsInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMigrationTaskGroup_urlList(listFileBucketName, destBucketName, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "URL_LIST"),
					resource.TestCheckResourceAttr(resourceName, "description", "test task group"),
					resource.TestCheckResourceAttr(resourceName, "action", "stop"),
				),
			},
		},
	})
}

func testMigrationTaskGroup_prefix(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_oms_migration_task_group" "test" {
  source_object {
    data_source = "HuaweiCloud"
    region      = "%[2]s"
    bucket      = hcso_obs_bucket.source.bucket
    access_key  = "%[3]s"
    secret_key  = "%[4]s"
    object      = [""]
  }

  destination_object {
    region     = "%[2]s"
    bucket     = hcso_obs_bucket.dest.bucket
    access_key = "%[3]s"
    secret_key = "%[4]s"
  }

  action                         = "stop"
  type                           = "PREFIX"
  enable_kms                     = true
  description                    = "test task group"
  migrate_since                  = "2023-01-02 15:04:05"
  object_overwrite_mode          = "CRC64_COMPARISON_OVERWRITE"
  consistency_check              = "crc64"
  enable_requester_pays          = true
  enable_failed_object_recording = true

  bandwidth_policy {
    max_bandwidth = 1
    start         = "15:00"
    end           = "16:00"
  }

  bandwidth_policy {
    max_bandwidth = 2
    start         = "16:00"
    end           = "17:00"
  }
}
`, baseConfig, acceptance.HCSO_REGION_NAME, acceptance.HCSO_ACCESS_KEY, acceptance.HCSO_SECRET_KEY)
}

func testMigrationTaskGroup_prefix_update(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_oms_migration_task_group" "test" {
  source_object {
    data_source = "HuaweiCloud"
    region      = "%[2]s"
    bucket      = hcso_obs_bucket.source.bucket
    access_key  = "%[3]s"
    secret_key  = "%[4]s"
    object      = [""]
  }

  destination_object {
    region     = "%[2]s"
    bucket     = hcso_obs_bucket.dest.bucket
    access_key = "%[3]s"
    secret_key = "%[4]s"
  }

  action                         = "start"
  type                           = "PREFIX"
  enable_kms                     = true
  description                    = "test task group"
  migrate_since                  = "2023-01-02 15:04:05"
  object_overwrite_mode          = "CRC64_COMPARISON_OVERWRITE"
  consistency_check              = "crc64"
  enable_requester_pays          = true
  enable_failed_object_recording = true

  bandwidth_policy {
    max_bandwidth = 2
    start         = "15:00"
    end           = "16:00"
  }

  bandwidth_policy {
    max_bandwidth = 3
    start         = "16:00"
    end           = "17:00"
  }
}
`, baseConfig, acceptance.HCSO_REGION_NAME, acceptance.HCSO_ACCESS_KEY, acceptance.HCSO_SECRET_KEY)
}

func testMigrationTaskGroup_list(baseConfig, listFileBucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "hcso_obs_bucket" "list_file_bucket" {
  bucket = "%[2]s"
  acl    = "private"
}

resource "hcso_obs_bucket_object" "list_file_object" {
  bucket        = hcso_obs_bucket.list_file_bucket.bucket
  key           = "test_folder/list_file.txt"
  content       = "test.txt"
}

resource "hcso_oms_migration_task_group" "test" {
  source_object {
    data_source      = "HuaweiCloud"
    region           = "%[3]s"
    bucket           = hcso_obs_bucket.source.bucket
    access_key       = "%[4]s"
    secret_key       = "%[5]s"
    list_file_bucket = hcso_obs_bucket.list_file_bucket.bucket
    list_file_key    = "test_folder/"
  }

  destination_object {
    region     = "%[3]s"
    bucket     = hcso_obs_bucket.dest.bucket
    access_key = "%[4]s"
    secret_key = "%[5]s"
  }

  action      = "stop"
  type        = "LIST"  
  description = "test task group"
}
`, baseConfig, listFileBucketName, acceptance.HCSO_REGION_NAME, acceptance.HCSO_ACCESS_KEY, acceptance.HCSO_SECRET_KEY)
}

func testMigrationTaskGroup_urlList(listFileBucketName, destBucketName, source string) string {
	return fmt.Sprintf(`
resource "hcso_obs_bucket" "list_file_bucket" {
  bucket = "%[1]s"
  acl    = "private"
}

resource "hcso_obs_bucket" "dest" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}

resource "hcso_obs_bucket_object" "list_file_object" {
  bucket        = hcso_obs_bucket.list_file_bucket.bucket
  key           = "folder/temp.txt"
  source        = "%[3]s"
  content_type  = "binary/octet-stream"
}

resource "hcso_oms_migration_task_group" "test" {
  source_object {
    data_source      = "URLSource"
    list_file_bucket = hcso_obs_bucket.list_file_bucket.bucket
    list_file_key    = "folder/"
  }

  destination_object {
    region     = "%[4]s"
    bucket     = hcso_obs_bucket.dest.bucket
    access_key = "%[5]s"
    secret_key = "%[6]s"
  }

  action      = "stop"
  type        = "URL_LIST"  
  description = "test task group"
}
`, listFileBucketName, destBucketName, source, acceptance.HCSO_REGION_NAME, acceptance.HCSO_ACCESS_KEY, acceptance.HCSO_SECRET_KEY)
}
