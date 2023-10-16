---
subcategory: "Enterprise Project Management Service (EPS)"
---

# hcso_enterprise_project

Use this data source to get an enterprise project from HuaweiCloud

## Example Usage

```hcl
data "hcso_enterprise_project" "test" {
  name = "test"
}
```

## Resources Supported Currently

<!-- markdownlint-disable MD033 -->
Service Name | Resource Name | Sub Resource Name
---- | --- | ---
AS  | hcso_as_group |
BCS | hcso_bcs_instance |
BMS | hcso_bms_instance |
CBR | hcso_cbr_vault |
CCE | hcso_cce_cluster | hcso_cce_node<br>hcso_cce_node_pool<br>hcso_cce_addon
CDM | hcso_cdm_cluster |
CDN | hcso_cdn_domain |
CES | hcso_ces_alarmrule |
DCS | hcso_dcs_instance |
DDS | hcso_dds_instance |
DMS | hcso_dms_kafka_instance<br>hcso_dms_rabbitmq_instance |
DNS | hcso_dns_ptrrecord<br>hcso_dns_zone |
ECS | hcso_compute_instance |
EIP | hcso_vpc_eip<br>hcso_vpc_bandwidth |
ELB | hcso_lb_loadbalancer |
Dedicated ELB | hcso_elb_certificate<br>hcso_elb_ipgroup<br>hcso_elb_loadbalancer |
EVS | hcso_evs_volume |
FGS | hcso_fgs_function |
GaussDB | hcso_gaussdb_cassandra_instance<br>hcso_gaussdb_mysql_instance<br>hcso_gaussdb_opengauss_instance |
IMS | hcso_images_image |
KMS | hcso_kms_key |
NAT | hcso_nat_gateway | hcso_nat_snat_rule<br>hcso_nat_dnat_rule
OBS | hcso_obs_bucket | hcso_obs_bucket_object<br>hcso_obs_bucket_policy
RDS | hcso_rds_instance<br>hcso_rds_read_replica_instance |
SFS | hcso_sfs_file_system<br>hcso_sfs_turbo | hcso_sfs_access_rule
SMN | hcso_smn_topic |
VPC | hcso_vpc<br>hcso_networking_secgroup | hcso_vpc_subnet<br>hcso_vpc_route<br>hcso_networking_secgroup_rule
<!-- markdownlint-enable MD033 -->

## Argument Reference

* `name` - (Optional, String) Specifies the enterprise project name. Fuzzy search is supported.

* `id` - (Optional, String) Specifies the ID of an enterprise project. The value 0 indicates enterprise project default.

* `status` - (Optional, Int) Specifies the status of an enterprise project.
    + 1 indicates Enabled.
    + 2 indicates Disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Provides supplementary information about the enterprise project.

* `created_at` - Specifies the time (UTC) when the enterprise project was created. Example: 2018-05-18T06:49:06Z

* `updated_at` - Specifies the time (UTC) when the enterprise project was modified. Example: 2018-05-28T02:21:36Z
