package internal

import (
	"context"
	"fmt"

	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcso/internal/hcso_config"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aad"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/antiddos"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apigateway"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/as"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/bms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbh"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ces"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cloudtable"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cmdb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cnad"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codearts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cph"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cpts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/css"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dbss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ddm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dis"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ecs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eg"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/er"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ga"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/gaussdb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ges"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/identitycenter"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/meeting"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/nat"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/obs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/oms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ram"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/scm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sdrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/swr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/tms"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ucs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vod"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpcep"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

const (
	defaultCloud       string = "myhuaweicloud.com"
	defaultEuropeCloud string = "myhuaweicloud.eu"
)

// Provider returns a schema.Provider for HuaweiCloud.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["region"],
				InputDefault: "cn-north-1",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_REGION_NAME",
					"OS_REGION_NAME",
				}, nil),
			},

			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["access_key"],
				RequiredWith: []string{"secret_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_ACCESS_KEY",
					"OS_ACCESS_KEY",
				}, nil),
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["secret_key"],
				RequiredWith: []string{"access_key"},
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_SECRET_KEY",
					"OS_SECRET_KEY",
				}, nil),
			},

			"security_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["security_token"],
				RequiredWith: []string{"access_key"},
				DefaultFunc:  schema.EnvDefaultFunc("HCSO_SECURITY_TOKEN", nil),
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_DOMAIN_ID",
					"OS_DOMAIN_ID",
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
				}, ""),
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
				}, ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_USER_NAME",
					"OS_USERNAME",
				}, ""),
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_USER_ID",
					"OS_USER_ID",
				}, ""),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: descriptions["password"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_USER_PASSWORD",
					"OS_PASSWORD",
				}, ""),
			},

			"assume_role": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_agency_name"],
							DefaultFunc: schema.EnvDefaultFunc("HCSO_ASSUME_ROLE_AGENCY_NAME", nil),
						},
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: descriptions["assume_role_domain_name"],
							DefaultFunc: schema.EnvDefaultFunc("HCSO_ASSUME_ROLE_DOMAIN_NAME", nil),
						},
					},
				},
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_PROJECT_ID",
					"OS_PROJECT_ID",
				}, nil),
			},

			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["project_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_PROJECT_NAME",
					"OS_PROJECT_NAME",
				}, nil),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_id"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_ID", ""),
			},

			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_name"],
				DefaultFunc: schema.EnvDefaultFunc("OS_TENANT_NAME", ""),
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["token"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_AUTH_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["insecure"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_INSECURE",
					"OS_INSECURE",
				}, false),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"agency_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_NAME", nil),
				Description:  descriptions["agency_name"],
				RequiredWith: []string{"agency_domain_name"},
			},

			"agency_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OS_AGENCY_DOMAIN_NAME", nil),
				Description:  descriptions["agency_domain_name"],
				RequiredWith: []string{"agency_name"},
			},

			"delegated_project": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELEGATED_PROJECT", ""),
				Description: descriptions["delegated_project"],
			},

			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["auth_url"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HCSO_AUTH_URL",
					"OS_AUTH_URL",
				}, nil),
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["cloud"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_CLOUD", ""),
			},

			"endpoints": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoints"],
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"regional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["regional"],
			},

			"shared_config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_config_file"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_SHARED_CONFIG_FILE", ""),
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_PROFILE", ""),
			},

			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["enterprise_project_id"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_ENTERPRISE_PROJECT_ID", ""),
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["max_retries"],
				DefaultFunc: schema.EnvDefaultFunc("HCSO_MAX_RETRIES", 5),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"hcso_apig_environments": apig.DataSourceEnvironments(),

			"hcso_as_configurations": as.DataSourceASConfigurations(),
			"hcso_as_groups":         as.DataSourceASGroups(),

			"hcso_account":            DataSourceAccount(),
			"hcso_availability_zones": DataSourceAvailabilityZones(),

			"hcso_bms_flavors": bms.DataSourceBmsFlavors(),

			"hcso_cbr_backup": cbr.DataSourceBackup(),
			"hcso_cbr_vaults": cbr.DataSourceVaults(),

			"hcso_cbh_instances": cbh.DataSourceCbhInstances(),

			"hcso_cce_addon_template": cce.DataSourceAddonTemplate(),
			"hcso_cce_cluster":        cce.DataSourceCCEClusterV3(),
			"hcso_cce_clusters":       cce.DataSourceCCEClusters(),
			"hcso_cce_node":           cce.DataSourceNode(),
			"hcso_cce_nodes":          cce.DataSourceNodes(),
			"hcso_cce_node_pool":      cce.DataSourceCCENodePoolV3(),
			"hcso_cci_namespaces":     cci.DataSourceCciNamespaces(),

			"hcso_cdm_flavors": DataSourceCdmFlavorV1(),

			"hcso_cdn_domain_statistics": cdn.DataSourceStatistics(),

			"hcso_cfw_firewalls": cfw.DataSourceFirewalls(),

			"hcso_cnad_advanced_instances":         cnad.DataSourceInstances(),
			"hcso_cnad_advanced_available_objects": cnad.DataSourceAvailableProtectedObjects(),
			"hcso_cnad_advanced_protected_objects": cnad.DataSourceProtectedObjects(),

			"hcso_compute_flavors":      ecs.DataSourceEcsFlavors(),
			"hcso_compute_instance":     ecs.DataSourceComputeInstance(),
			"hcso_compute_instances":    ecs.DataSourceComputeInstances(),
			"hcso_compute_servergroups": ecs.DataSourceComputeServerGroups(),

			"hcso_cdm_clusters": cdm.DataSourceCdmClusters(),

			"hcso_cph_server_flavors": cph.DataSourceServerFlavors(),
			"hcso_cph_phone_flavors":  cph.DataSourcePhoneFlavors(),
			"hcso_cph_phone_images":   cph.DataSourcePhoneImages(),

			"hcso_csms_secret_version": dew.DataSourceDewCsmsSecret(),
			"hcso_css_flavors":         css.DataSourceCssFlavors(),

			"hcso_dcs_flavors":        dcs.DataSourceDcsFlavorsV2(),
			"hcso_dcs_maintainwindow": dcs.DataSourceDcsMaintainWindow(),
			"hcso_dcs_instances":      dcs.DataSourceDcsInstance(),

			"hcso_dds_flavors":   dds.DataSourceDDSFlavorV3(),
			"hcso_dds_instances": dds.DataSourceDdsInstance(),

			"hcso_dms_kafka_flavors":   dms.DataSourceKafkaFlavors(),
			"hcso_dms_kafka_instances": dms.DataSourceDmsKafkaInstances(),
			"hcso_dms_product":         dms.DataSourceDmsProduct(),
			"hcso_dms_maintainwindow":  dms.DataSourceDmsMaintainWindow(),

			"hcso_dms_rabbitmq_flavors": dms.DataSourceRabbitMQFlavors(),

			"hcso_dms_rocketmq_broker":    dms.DataSourceDmsRocketMQBroker(),
			"hcso_dms_rocketmq_instances": dms.DataSourceDmsRocketMQInstances(),

			"hcso_dns_zones":      dns.DataSourceZones(),
			"hcso_dns_recordsets": dns.DataSourceRecordsets(),

			"hcso_eg_custom_event_sources": eg.DataSourceCustomEventSources(),

			"hcso_enterprise_project": eps.DataSourceEnterpriseProject(),

			"hcso_er_attachments":  er.DataSourceAttachments(),
			"hcso_er_instances":    er.DataSourceInstances(),
			"hcso_er_route_tables": er.DataSourceRouteTables(),

			"hcso_evs_volumes":      evs.DataSourceEvsVolumesV2(),
			"hcso_fgs_dependencies": fgs.DataSourceFunctionGraphDependencies(),

			"hcso_gaussdb_cassandra_dedicated_resource": gaussdb.DataSourceGeminiDBDehResource(),
			"hcso_gaussdb_cassandra_flavors":            gaussdb.DataSourceCassandraFlavors(),
			"hcso_gaussdb_nosql_flavors":                gaussdb.DataSourceGaussDBNoSQLFlavors(),
			"hcso_gaussdb_cassandra_instance":           gaussdb.DataSourceGeminiDBInstance(),
			"hcso_gaussdb_cassandra_instances":          gaussdb.DataSourceGeminiDBInstances(),
			"hcso_gaussdb_opengauss_instance":           gaussdb.DataSourceOpenGaussInstance(),
			"hcso_gaussdb_opengauss_instances":          gaussdb.DataSourceOpenGaussInstances(),
			"hcso_gaussdb_mysql_configuration":          gaussdb.DataSourceGaussdbMysqlConfigurations(),
			"hcso_gaussdb_mysql_dedicated_resource":     gaussdb.DataSourceGaussDBMysqlDehResource(),
			"hcso_gaussdb_mysql_flavors":                gaussdb.DataSourceGaussdbMysqlFlavors(),
			"hcso_gaussdb_mysql_instance":               gaussdb.DataSourceGaussDBMysqlInstance(),
			"hcso_gaussdb_mysql_instances":              gaussdb.DataSourceGaussDBMysqlInstances(),
			"hcso_gaussdb_redis_instance":               gaussdb.DataSourceGaussRedisInstance(),

			"hcso_identity_permissions": iam.DataSourceIdentityPermissions(),
			"hcso_identity_role":        iam.DataSourceIdentityRole(),
			"hcso_identity_custom_role": iam.DataSourceIdentityCustomRole(),
			"hcso_identity_group":       iam.DataSourceIdentityGroup(),
			"hcso_identity_projects":    iam.DataSourceIdentityProjects(),
			"hcso_identity_users":       iam.DataSourceIdentityUsers(),

			"hcso_identitycenter_instance": identitycenter.DataSourceIdentityCenter(),
			"hcso_identitycenter_groups":   identitycenter.DataSourceIdentityCenterGroups(),
			"hcso_identitycenter_users":    identitycenter.DataSourceIdentityCenterUsers(),

			"hcso_iec_bandwidths":     dataSourceIECBandWidths(),
			"hcso_iec_eips":           dataSourceIECNetworkEips(),
			"hcso_iec_flavors":        dataSourceIecFlavors(),
			"hcso_iec_images":         dataSourceIecImages(),
			"hcso_iec_keypair":        dataSourceIECKeypair(),
			"hcso_iec_network_acl":    dataSourceIECNetworkACL(),
			"hcso_iec_port":           DataSourceIECPort(),
			"hcso_iec_security_group": dataSourceIECSecurityGroup(),
			"hcso_iec_server":         dataSourceIECServer(),
			"hcso_iec_sites":          dataSourceIecSites(),
			"hcso_iec_vpc":            DataSourceIECVpc(),
			"hcso_iec_vpc_subnets":    DataSourceIECVpcSubnets(),

			"hcso_images_image":  ims.DataSourceImagesImageV2(),
			"hcso_images_images": ims.DataSourceImagesImages(),

			"hcso_kms_key":      dew.DataSourceKmsKey(),
			"hcso_kms_data_key": dew.DataSourceKmsDataKeyV1(),
			"hcso_kps_keypairs": dew.DataSourceKeypairs(),

			"hcso_lb_listeners":    lb.DataSourceListeners(),
			"hcso_lb_loadbalancer": lb.DataSourceELBV2Loadbalancer(),
			"hcso_lb_certificate":  lb.DataSourceLBCertificateV2(),
			"hcso_lb_pools":        lb.DataSourcePools(),

			"hcso_elb_certificate": elb.DataSourceELBCertificateV3(),
			"hcso_elb_flavors":     elb.DataSourceElbFlavorsV3(),
			"hcso_elb_pools":       elb.DataSourcePools(),

			"hcso_nat_gateway": nat.DataSourcePublicGateway(),

			"hcso_networking_port":      vpc.DataSourceNetworkingPortV2(),
			"hcso_networking_secgroups": vpc.DataSourceNetworkingSecGroups(),

			"hcso_modelarts_datasets":         modelarts.DataSourceDatasets(),
			"hcso_modelarts_dataset_versions": modelarts.DataSourceDatasetVerions(),
			"hcso_modelarts_notebook_images":  modelarts.DataSourceNotebookImages(),
			"hcso_modelarts_notebook_flavors": modelarts.DataSourceNotebookFlavors(),
			"hcso_modelarts_service_flavors":  modelarts.DataSourceServiceFlavors(),
			"hcso_modelarts_models":           modelarts.DataSourceModels(),
			"hcso_modelarts_model_templates":  modelarts.DataSourceModelTemplates(),
			"hcso_modelarts_workspaces":       modelarts.DataSourceWorkspaces(),
			"hcso_modelarts_services":         modelarts.DataSourceServices(),
			"hcso_modelarts_resource_flavors": modelarts.DataSourceResourceFlavors(),

			"hcso_obs_buckets":       obs.DataSourceObsBuckets(),
			"hcso_obs_bucket_object": obs.DataSourceObsBucketObject(),

			"hcso_ram_resource_permissions": ram.DataSourceRAMPermissions(),

			"hcso_rds_flavors":         rds.DataSourceRdsFlavor(),
			"hcso_rds_engine_versions": rds.DataSourceRdsEngineVersionsV3(),
			"hcso_rds_instances":       rds.DataSourceRdsInstances(),
			"hcso_rds_backups":         rds.DataSourceBackup(),
			"hcso_rds_storage_types":   rds.DataSourceStoragetype(),

			"hcso_rms_policy_definitions": rms.DataSourcePolicyDefinitions(),

			"hcso_sdrs_domain": sdrs.DataSourceSDRSDomain(),

			"hcso_servicestage_component_runtimes": servicestage.DataSourceComponentRuntimes(),

			"hcso_smn_topics": smn.DataSourceTopics(),

			"hcso_sms_source_servers": sms.DataSourceServers(),

			"hcso_scm_certificates": scm.DataSourceCertificates(),

			"hcso_sfs_file_system": sfs.DataSourceSFSFileSystemV2(),
			"hcso_sfs_turbos":      sfs.DataSourceTurbos(),

			"hcso_vpc_bandwidth": eip.DataSourceBandWidth(),
			"hcso_vpc_eip":       eip.DataSourceVpcEip(),
			"hcso_vpc_eips":      eip.DataSourceVpcEips(),

			"hcso_vpc":                    vpc.DataSourceVpcV1(),
			"hcso_vpcs":                   vpc.DataSourceVpcs(),
			"hcso_vpc_ids":                vpc.DataSourceVpcIdsV1(),
			"hcso_vpc_peering_connection": vpc.DataSourceVpcPeeringConnectionV2(),
			"hcso_vpc_route_table":        vpc.DataSourceVPCRouteTable(),
			"hcso_vpc_subnet":             vpc.DataSourceVpcSubnetV1(),
			"hcso_vpc_subnets":            vpc.DataSourceVpcSubnets(),
			"hcso_vpc_subnet_ids":         vpc.DataSourceVpcSubnetIdsV1(),

			"hcso_vpcep_public_services": vpcep.DataSourceVPCEPPublicServices(),

			"hcso_waf_certificate":         waf.DataSourceWafCertificateV1(),
			"hcso_waf_policies":            waf.DataSourceWafPoliciesV1(),
			"hcso_waf_dedicated_instances": waf.DataSourceWafDedicatedInstancesV1(),
			"hcso_waf_reference_tables":    waf.DataSourceWafReferenceTablesV1(),
			"hcso_waf_instance_groups":     waf.DataSourceWafInstanceGroups(),
			"hcso_dws_flavors":             dws.DataSourceDwsFlavors(),

			// Legacy
			"hcso_images_image_v2":        ims.DataSourceImagesImageV2(),
			"hcso_networking_port_v2":     vpc.DataSourceNetworkingPortV2(),
			"hcso_networking_secgroup_v2": vpc.DataSourceNetworkingSecGroup(),

			"hcso_kms_key_v1":      dew.DataSourceKmsKey(),
			"hcso_kms_data_key_v1": dew.DataSourceKmsDataKeyV1(),

			"hcso_rds_flavors_v3":     rds.DataSourceRdsFlavor(),
			"hcso_sfs_file_system_v2": sfs.DataSourceSFSFileSystemV2(),

			"hcso_vpc_v1":                    vpc.DataSourceVpcV1(),
			"hcso_vpc_ids_v1":                vpc.DataSourceVpcIdsV1(),
			"hcso_vpc_peering_connection_v2": vpc.DataSourceVpcPeeringConnectionV2(),
			"hcso_vpc_subnet_v1":             vpc.DataSourceVpcSubnetV1(),
			"hcso_vpc_subnet_ids_v1":         vpc.DataSourceVpcSubnetIdsV1(),

			"hcso_cce_cluster_v3": cce.DataSourceCCEClusterV3(),
			"hcso_cce_node_v3":    cce.DataSourceNode(),

			"hcso_dms_product_v1":        dms.DataSourceDmsProduct(),
			"hcso_dms_maintainwindow_v1": dms.DataSourceDmsMaintainWindow(),

			"hcso_dcs_maintainwindow_v1": dcs.DataSourceDcsMaintainWindow(),

			"hcso_dds_flavors_v3":   dds.DataSourceDDSFlavorV3(),
			"hcso_identity_role_v3": iam.DataSourceIdentityRole(),
			"hcso_cdm_flavors_v1":   DataSourceCdmFlavorV1(),

			"hcso_ddm_engines":        ddm.DataSourceDdmEngines(),
			"hcso_ddm_flavors":        ddm.DataSourceDdmFlavors(),
			"hcso_ddm_instance_nodes": ddm.DataSourceDdmInstanceNodes(),
			"hcso_ddm_instances":      ddm.DataSourceDdmInstances(),
			"hcso_ddm_schemas":        ddm.DataSourceDdmSchemas(),
			"hcso_ddm_accounts":       ddm.DataSourceDdmAccounts(),

			"hcso_organizations_organization":         organizations.DataSourceOrganization(),
			"hcso_organizations_organizational_units": organizations.DataSourceOrganizationalUnits(),
			"hcso_organizations_accounts":             organizations.DataSourceAccounts(),

			// Deprecated ongoing (without DeprecationMessage), used by other providers
			"hcso_vpc_route":        vpc.DataSourceVpcRouteV2(),
			"hcso_vpc_route_ids":    vpc.DataSourceVpcRouteIdsV2(),
			"hcso_vpc_route_v2":     vpc.DataSourceVpcRouteV2(),
			"hcso_vpc_route_ids_v2": vpc.DataSourceVpcRouteIdsV2(),

			// Deprecated
			"hcso_antiddos":                      deprecated.DataSourceAntiDdosV1(),
			"hcso_antiddos_v1":                   deprecated.DataSourceAntiDdosV1(),
			"hcso_compute_availability_zones_v2": deprecated.DataSourceComputeAvailabilityZonesV2(),
			"hcso_csbs_backup":                   deprecated.DataSourceCSBSBackupV1(),
			"hcso_csbs_backup_policy":            deprecated.DataSourceCSBSBackupPolicyV1(),
			"hcso_csbs_backup_policy_v1":         deprecated.DataSourceCSBSBackupPolicyV1(),
			"hcso_csbs_backup_v1":                deprecated.DataSourceCSBSBackupV1(),
			"hcso_networking_network_v2":         deprecated.DataSourceNetworkingNetworkV2(),
			"hcso_networking_subnet_v2":          deprecated.DataSourceNetworkingSubnetV2(),
			"hcso_cts_tracker":                   deprecated.DataSourceCTSTrackerV1(),
			"hcso_dcs_az":                        deprecated.DataSourceDcsAZV1(),
			"hcso_dcs_az_v1":                     deprecated.DataSourceDcsAZV1(),
			"hcso_dcs_product":                   deprecated.DataSourceDcsProductV1(),
			"hcso_dcs_product_v1":                deprecated.DataSourceDcsProductV1(),
			"hcso_dms_az":                        deprecated.DataSourceDmsAZ(),
			"hcso_dms_az_v1":                     deprecated.DataSourceDmsAZ(),
			"hcso_vbs_backup_policy":             deprecated.DataSourceVBSBackupPolicyV2(),
			"hcso_vbs_backup":                    deprecated.DataSourceVBSBackupV2(),
			"hcso_vbs_backup_policy_v2":          deprecated.DataSourceVBSBackupPolicyV2(),
			"hcso_vbs_backup_v2":                 deprecated.DataSourceVBSBackupV2(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"hcso_aad_forward_rule": aad.ResourceForwardRule(),

			"hcso_antiddos_basic": antiddos.ResourceCloudNativeAntiDdos(),

			"hcso_aom_alarm_rule":             aom.ResourceAlarmRule(),
			"hcso_aom_event_alarm_rule":       aom.ResourceEventAlarmRule(),
			"hcso_aom_service_discovery_rule": aom.ResourceServiceDiscoveryRule(),
			"hcso_aom_alarm_action_rule":      aom.ResourceAlarmActionRule(),
			"hcso_aom_alarm_silence_rule":     aom.ResourceAlarmSilenceRule(),

			"hcso_rfs_stack": rfs.ResourceStack(),

			"hcso_api_gateway_api":         ResourceAPIGatewayAPI(),
			"hcso_api_gateway_environment": apigateway.ResourceEnvironment(),
			"hcso_api_gateway_group":       ResourceAPIGatewayGroup(),

			"hcso_apig_acl_policy":                  apig.ResourceAclPolicy(),
			"hcso_apig_acl_policy_associate":        apig.ResourceAclPolicyAssociate(),
			"hcso_apig_api":                         apig.ResourceApigAPIV2(),
			"hcso_apig_api_publishment":             apig.ResourceApigApiPublishment(),
			"hcso_apig_appcode":                     apig.ResourceAppcode(),
			"hcso_apig_application":                 apig.ResourceApigApplicationV2(),
			"hcso_apig_application_authorization":   apig.ResourceAppAuth(),
			"hcso_apig_certificate":                 apig.ResourceCertificate(),
			"hcso_apig_channel":                     apig.ResourceChannel(),
			"hcso_apig_custom_authorizer":           apig.ResourceApigCustomAuthorizerV2(),
			"hcso_apig_environment":                 apig.ResourceApigEnvironmentV2(),
			"hcso_apig_group":                       apig.ResourceApigGroupV2(),
			"hcso_apig_instance_routes":             apig.ResourceInstanceRoutes(),
			"hcso_apig_instance":                    apig.ResourceApigInstanceV2(),
			"hcso_apig_plugin_associate":            apig.ResourcePluginAssociate(),
			"hcso_apig_plugin":                      apig.ResourcePlugin(),
			"hcso_apig_response":                    apig.ResourceApigResponseV2(),
			"hcso_apig_signature_associate":         apig.ResourceSignatureAssociate(),
			"hcso_apig_signature":                   apig.ResourceSignature(),
			"hcso_apig_throttling_policy_associate": apig.ResourceThrottlingPolicyAssociate(),
			"hcso_apig_throttling_policy":           apig.ResourceApigThrottlingPolicyV2(),

			"hcso_as_configuration":    as.ResourceASConfiguration(),
			"hcso_as_group":            as.ResourceASGroup(),
			"hcso_as_lifecycle_hook":   as.ResourceASLifecycleHook(),
			"hcso_as_instance_attach":  as.ResourceASInstanceAttach(),
			"hcso_as_notification":     as.ResourceAsNotification(),
			"hcso_as_policy":           as.ResourceASPolicy(),
			"hcso_as_bandwidth_policy": as.ResourceASBandWidthPolicy(),

			"hcso_bms_instance": bms.ResourceBmsInstance(),
			"hcso_bcs_instance": resourceBCSInstanceV2(),

			"hcso_cbr_policy": cbr.ResourcePolicy(),
			"hcso_cbr_vault":  cbr.ResourceVault(),

			"hcso_cbh_instance": cbh.ResourceCBHInstance(),

			"hcso_cc_connection":             cc.ResourceCloudConnection(),
			"hcso_cc_network_instance":       cc.ResourceNetworkInstance(),
			"hcso_cc_bandwidth_package":      cc.ResourceBandwidthPackage(),
			"hcso_cc_inter_region_bandwidth": cc.ResourceInterRegionBandwidth(),

			"hcso_cce_cluster":     cce.ResourceCluster(),
			"hcso_cce_node":        cce.ResourceNode(),
			"hcso_cce_node_attach": cce.ResourceNodeAttach(),
			"hcso_cce_addon":       cce.ResourceAddon(),
			"hcso_cce_node_pool":   cce.ResourceNodePool(),
			"hcso_cce_namespace":   cce.ResourceCCENamespaceV1(),
			"hcso_cce_pvc":         cce.ResourceCcePersistentVolumeClaimsV1(),
			"hcso_cce_partition":   cce.ResourcePartition(),

			"hcso_cts_tracker":      cts.ResourceCTSTracker(),
			"hcso_cts_data_tracker": cts.ResourceCTSDataTracker(),
			"hcso_cts_notification": cts.ResourceCTSNotification(),
			"hcso_cci_namespace":    cci.ResourceCciNamespace(),
			"hcso_cci_network":      cci.ResourceCciNetworkV1(),
			"hcso_cci_pvc":          ResourceCCIPersistentVolumeClaimV1(),

			"hcso_cdm_cluster": cdm.ResourceCdmCluster(),
			"hcso_cdm_job":     cdm.ResourceCdmJob(),
			"hcso_cdm_link":    cdm.ResourceCdmLink(),

			"hcso_cdn_domain":         resourceCdnDomainV1(),
			"hcso_ces_alarmrule":      ces.ResourceAlarmRule(),
			"hcso_ces_resource_group": ces.ResourceResourceGroup(),
			"hcso_ces_alarm_template": ces.ResourceCesAlarmTemplate(),

			"hcso_cfw_address_group":        cfw.ResourceAddressGroup(),
			"hcso_cfw_address_group_member": cfw.ResourceAddressGroupMember(),
			"hcso_cfw_black_white_list":     cfw.ResourceBlackWhiteList(),
			"hcso_cfw_eip_protection":       cfw.ResourceEipProtection(),
			"hcso_cfw_protection_rule":      cfw.ResourceProtectionRule(),
			"hcso_cfw_service_group":        cfw.ResourceServiceGroup(),
			"hcso_cfw_service_group_member": cfw.ResourceServiceGroupMember(),

			"hcso_cloudtable_cluster": cloudtable.ResourceCloudTableCluster(),

			"hcso_cnad_advanced_black_white_list": cnad.ResourceBlackWhiteList(),
			"hcso_cnad_advanced_policy":           cnad.ResourceCNADAdvancedPolicy(),
			"hcso_cnad_advanced_policy_associate": cnad.ResourcePolicyAssociate(),
			"hcso_cnad_advanced_protected_object": cnad.ResourceProtectedObject(),

			"hcso_compute_instance":         ecs.ResourceComputeInstance(),
			"hcso_compute_interface_attach": ecs.ResourceComputeInterfaceAttach(),
			"hcso_compute_keypair":          ResourceComputeKeypairV2(),
			"hcso_compute_servergroup":      ecs.ResourceComputeServerGroup(),
			"hcso_compute_eip_associate":    ecs.ResourceComputeEIPAssociate(),
			"hcso_compute_volume_attach":    ecs.ResourceComputeVolumeAttach(),

			"hcso_cph_server": cph.ResourceCphServer(),

			"hcso_cse_microservice":          cse.ResourceMicroservice(),
			"hcso_cse_microservice_engine":   cse.ResourceMicroserviceEngine(),
			"hcso_cse_microservice_instance": cse.ResourceMicroserviceInstance(),

			"hcso_csms_secret": dew.ResourceCsmsSecret(),

			"hcso_css_cluster":   css.ResourceCssCluster(),
			"hcso_css_snapshot":  css.ResourceCssSnapshot(),
			"hcso_css_thesaurus": css.ResourceCssthesaurus(),

			"hcso_dbss_instance": dbss.ResourceInstance(),

			"hcso_dc_virtual_gateway":   dc.ResourceVirtualGateway(),
			"hcso_dc_virtual_interface": dc.ResourceVirtualInterface(),

			"hcso_dcs_instance": dcs.ResourceDcsInstance(),
			"hcso_dcs_backup":   dcs.ResourceDcsBackup(),

			"hcso_dds_database_role":      dds.ResourceDatabaseRole(),
			"hcso_dds_database_user":      dds.ResourceDatabaseUser(),
			"hcso_dds_instance":           dds.ResourceDdsInstanceV3(),
			"hcso_dds_backup":             dds.ResourceDdsBackup(),
			"hcso_dds_parameter_template": dds.ResourceDdsParameterTemplate(),
			"hcso_dds_audit_log_policy":   dds.ResourceDdsAuditLogPolicy(),

			"hcso_ddm_instance": ddm.ResourceDdmInstance(),
			"hcso_ddm_schema":   ddm.ResourceDdmSchema(),
			"hcso_ddm_account":  ddm.ResourceDdmAccount(),

			"hcso_dis_stream": dis.ResourceDisStream(),

			"hcso_dli_database":              dli.ResourceDliSqlDatabaseV1(),
			"hcso_dli_package":               dli.ResourceDliPackageV2(),
			"hcso_dli_queue":                 dli.ResourceDliQueue(),
			"hcso_dli_spark_job":             dli.ResourceDliSparkJobV2(),
			"hcso_dli_sql_job":               dli.ResourceSqlJob(),
			"hcso_dli_table":                 dli.ResourceDliTable(),
			"hcso_dli_flinksql_job":          dli.ResourceFlinkSqlJob(),
			"hcso_dli_flinkjar_job":          dli.ResourceFlinkJarJob(),
			"hcso_dli_permission":            dli.ResourceDliPermission(),
			"hcso_dli_datasource_connection": dli.ResourceDatasourceConnection(),
			"hcso_dli_datasource_auth":       dli.ResourceDatasourceAuth(),
			"hcso_dli_template_sql":          dli.ResourceSQLTemplate(),
			"hcso_dli_template_flink":        dli.ResourceFlinkTemplate(),
			"hcso_dli_global_variable":       dli.ResourceGlobalVariable(),
			"hcso_dli_template_spark":        dli.ResourceSparkTemplate(),
			"hcso_dli_agency":                dli.ResourceDliAgency(),

			"hcso_dms_kafka_user":        dms.ResourceDmsKafkaUser(),
			"hcso_dms_kafka_permissions": dms.ResourceDmsKafkaPermissions(),
			"hcso_dms_kafka_instance":    dms.ResourceDmsKafkaInstance(),
			"hcso_dms_kafka_topic":       dms.ResourceDmsKafkaTopic(),
			"hcso_dms_rabbitmq_instance": dms.ResourceDmsRabbitmqInstance(),

			"hcso_dms_rocketmq_instance":       dms.ResourceDmsRocketMQInstance(),
			"hcso_dms_rocketmq_consumer_group": dms.ResourceDmsRocketMQConsumerGroup(),
			"hcso_dms_rocketmq_topic":          dms.ResourceDmsRocketMQTopic(),
			"hcso_dms_rocketmq_user":           dms.ResourceDmsRocketMQUser(),

			"hcso_dns_custom_line": dns.ResourceDNSCustomLine(),
			"hcso_dns_ptrrecord":   dns.ResourceDNSPtrRecord(),
			"hcso_dns_recordset":   dns.ResourceDNSRecordset(),
			"hcso_dns_zone":        dns.ResourceDNSZone(),

			"hcso_drs_job": drs.ResourceDrsJob(),

			"hcso_dws_cluster":            dws.ResourceDwsCluster(),
			"hcso_dws_event_subscription": dws.ResourceDwsEventSubs(),
			"hcso_dws_alarm_subscription": dws.ResourceDwsAlarmSubs(),
			"hcso_dws_snapshot":           dws.ResourceDwsSnapshot(),
			"hcso_dws_snapshot_policy":    dws.ResourceDwsSnapshotPolicy(),
			"hcso_dws_ext_data_source":    dws.ResourceDwsExtDataSource(),

			"hcso_eg_custom_event_source": eg.ResourceCustomEventSource(),
			"hcso_eg_endpoint":            eg.ResourceEndpoint(),

			"hcso_elb_certificate":     elb.ResourceCertificateV3(),
			"hcso_elb_l7policy":        elb.ResourceL7PolicyV3(),
			"hcso_elb_l7rule":          elb.ResourceL7RuleV3(),
			"hcso_elb_listener":        elb.ResourceListenerV3(),
			"hcso_elb_loadbalancer":    elb.ResourceLoadBalancerV3(),
			"hcso_elb_monitor":         elb.ResourceMonitorV3(),
			"hcso_elb_ipgroup":         elb.ResourceIpGroupV3(),
			"hcso_elb_pool":            elb.ResourcePoolV3(),
			"hcso_elb_member":          elb.ResourceMemberV3(),
			"hcso_elb_logtank":         elb.ResourceLogTank(),
			"hcso_elb_security_policy": elb.ResourceSecurityPolicy(),

			"hcso_enterprise_project": eps.ResourceEnterpriseProject(),

			"hcso_er_association":    er.ResourceAssociation(),
			"hcso_er_instance":       er.ResourceInstance(),
			"hcso_er_propagation":    er.ResourcePropagation(),
			"hcso_er_route_table":    er.ResourceRouteTable(),
			"hcso_er_static_route":   er.ResourceStaticRoute(),
			"hcso_er_vpc_attachment": er.ResourceVpcAttachment(),

			"hcso_evs_snapshot": ResourceEvsSnapshotV2(),
			"hcso_evs_volume":   evs.ResourceEvsVolume(),

			"hcso_fgs_async_invoke_configuration": fgs.ResourceAsyncInvokeConfiguration(),
			"hcso_fgs_dependency":                 fgs.ResourceFgsDependency(),
			"hcso_fgs_function":                   fgs.ResourceFgsFunctionV2(),
			"hcso_fgs_trigger":                    fgs.ResourceFunctionGraphTrigger(),

			"hcso_ga_accelerator":    ga.ResourceAccelerator(),
			"hcso_ga_listener":       ga.ResourceListener(),
			"hcso_ga_endpoint_group": ga.ResourceEndpointGroup(),
			"hcso_ga_endpoint":       ga.ResourceEndpoint(),
			"hcso_ga_health_check":   ga.ResourceHealthCheck(),

			"hcso_gaussdb_cassandra_instance": gaussdb.ResourceGeminiDBInstanceV3(),

			"hcso_gaussdb_mysql_instance":           gaussdb.ResourceGaussDBInstance(),
			"hcso_gaussdb_mysql_proxy":              gaussdb.ResourceGaussDBProxy(),
			"hcso_gaussdb_mysql_database":           gaussdb.ResourceGaussDBDatabase(),
			"hcso_gaussdb_mysql_account":            gaussdb.ResourceGaussDBAccount(),
			"hcso_gaussdb_mysql_account_privilege":  gaussdb.ResourceGaussDBAccountPrivilege(),
			"hcso_gaussdb_mysql_sql_control_rule":   gaussdb.ResourceGaussDBSqlControlRule(),
			"hcso_gaussdb_mysql_parameter_template": gaussdb.ResourceGaussDBMysqlTemplate(),

			"hcso_gaussdb_opengauss_instance": gaussdb.ResourceOpenGaussInstance(),

			"hcso_gaussdb_redis_instance":      gaussdb.ResourceGaussRedisInstanceV3(),
			"hcso_gaussdb_redis_eip_associate": gaussdb.ResourceGaussRedisEipAssociate(),

			"hcso_gaussdb_influx_instance": gaussdb.ResourceGaussDBInfluxInstanceV3(),
			"hcso_gaussdb_mongo_instance":  gaussdb.ResourceGaussDBMongoInstanceV3(),

			"hcso_ges_graph":    ges.ResourceGesGraph(),
			"hcso_ges_metadata": ges.ResourceGesMetadata(),
			"hcso_ges_backup":   ges.ResourceGesBackup(),

			"hcso_hss_host_group": hss.ResourceHostGroup(),

			"hcso_identity_access_key":            iam.ResourceIdentityKey(),
			"hcso_identity_acl":                   iam.ResourceIdentityACL(),
			"hcso_identity_agency":                iam.ResourceIAMAgencyV3(),
			"hcso_identity_group":                 iam.ResourceIdentityGroup(),
			"hcso_identity_group_membership":      iam.ResourceIdentityGroupMembership(),
			"hcso_identity_group_role_assignment": iam.ResourceIdentityGroupRoleAssignment(),
			"hcso_identity_project":               iam.ResourceIdentityProject(),
			"hcso_identity_role":                  iam.ResourceIdentityRole(),
			"hcso_identity_role_assignment":       iam.ResourceIdentityGroupRoleAssignment(),
			"hcso_identity_user":                  iam.ResourceIdentityUser(),
			"hcso_identity_user_role_assignment":  iam.ResourceIdentityUserRoleAssignment(),
			"hcso_identity_provider":              iam.ResourceIdentityProvider(),
			"hcso_identity_password_policy":       iam.ResourceIdentityPasswordPolicy(),
			"hcso_identity_protection_policy":     iam.ResourceIdentityProtectionPolicy(),

			"hcso_identitycenter_user":                     identitycenter.ResourceIdentityCenterUser(),
			"hcso_identitycenter_group":                    identitycenter.ResourceIdentityCenterGroup(),
			"hcso_identitycenter_group_membership":         identitycenter.ResourceGroupMembership(),
			"hcso_identitycenter_permission_set":           identitycenter.ResourcePermissionSet(),
			"hcso_identitycenter_system_policy_attachment": identitycenter.ResourceSystemPolicyAttachment(),
			"hcso_identitycenter_account_assignment":       identitycenter.ResourceIdentityCenterAccountAssignment(),

			"hcso_iec_eip":                 resourceIecNetworkEip(),
			"hcso_iec_keypair":             resourceIecKeypair(),
			"hcso_iec_network_acl":         resourceIecNetworkACL(),
			"hcso_iec_network_acl_rule":    resourceIecNetworkACLRule(),
			"hcso_iec_security_group":      resourceIecSecurityGroup(),
			"hcso_iec_security_group_rule": resourceIecSecurityGroupRule(),
			"hcso_iec_server":              resourceIecServer(),
			"hcso_iec_vip":                 resourceIecVipV1(),
			"hcso_iec_vpc":                 ResourceIecVpc(),
			"hcso_iec_vpc_subnet":          resourceIecSubnet(),

			"hcso_images_image":                ims.ResourceImsImage(),
			"hcso_images_image_copy":           ims.ResourceImsImageCopy(),
			"hcso_images_image_share":          ims.ResourceImsImageShare(),
			"hcso_images_image_share_accepter": ims.ResourceImsImageShareAccepter(),

			"hcso_iotda_space":               iotda.ResourceSpace(),
			"hcso_iotda_product":             iotda.ResourceProduct(),
			"hcso_iotda_device":              iotda.ResourceDevice(),
			"hcso_iotda_device_group":        iotda.ResourceDeviceGroup(),
			"hcso_iotda_dataforwarding_rule": iotda.ResourceDataForwardingRule(),
			"hcso_iotda_amqp":                iotda.ResourceAmqp(),
			"hcso_iotda_device_certificate":  iotda.ResourceDeviceCertificate(),
			"hcso_iotda_device_linkage_rule": iotda.ResourceDeviceLinkageRule(),

			"hcso_kms_key":     dew.ResourceKmsKey(),
			"hcso_kps_keypair": dew.ResourceKeypair(),
			"hcso_kms_grant":   dew.ResourceKmsGrant(),

			"hcso_lb_certificate":  lb.ResourceCertificateV2(),
			"hcso_lb_l7policy":     lb.ResourceL7PolicyV2(),
			"hcso_lb_l7rule":       lb.ResourceL7RuleV2(),
			"hcso_lb_loadbalancer": lb.ResourceLoadBalancer(),
			"hcso_lb_listener":     lb.ResourceListener(),
			"hcso_lb_member":       lb.ResourceMemberV2(),
			"hcso_lb_monitor":      lb.ResourceMonitorV2(),
			"hcso_lb_pool":         lb.ResourcePoolV2(),
			"hcso_lb_whitelist":    lb.ResourceWhitelistV2(),

			"hcso_live_domain":               live.ResourceDomain(),
			"hcso_live_recording":            live.ResourceRecording(),
			"hcso_live_record_callback":      live.ResourceRecordCallback(),
			"hcso_live_transcoding":          live.ResourceTranscoding(),
			"hcso_live_snapshot":             live.ResourceLiveSnapshot(),
			"hcso_live_bucket_authorization": live.ResourceLiveBucketAuthorization(),

			"hcso_lts_group":      lts.ResourceLTSGroup(),
			"hcso_lts_stream":     lts.ResourceLTSStream(),
			"hcso_lts_host_group": lts.ResourceHostGroup(),
			"hcso_lts_transfer":   lts.ResourceLtsTransfer(),

			"hcso_mapreduce_cluster":         mrs.ResourceMRSClusterV2(),
			"hcso_mapreduce_job":             mrs.ResourceMRSJobV2(),
			"hcso_mapreduce_data_connection": mrs.ResourceDataConnection(),

			"hcso_meeting_admin_assignment": meeting.ResourceAdminAssignment(),
			"hcso_meeting_conference":       meeting.ResourceConference(),
			"hcso_meeting_user":             meeting.ResourceUser(),

			"hcso_modelarts_dataset":                modelarts.ResourceDataset(),
			"hcso_modelarts_dataset_version":        modelarts.ResourceDatasetVersion(),
			"hcso_modelarts_notebook":               modelarts.ResourceNotebook(),
			"hcso_modelarts_notebook_mount_storage": modelarts.ResourceNotebookMountStorage(),
			"hcso_modelarts_model":                  modelarts.ResourceModelartsModel(),
			"hcso_modelarts_service":                modelarts.ResourceModelartsService(),
			"hcso_modelarts_workspace":              modelarts.ResourceModelartsWorkspace(),
			"hcso_modelarts_authorization":          modelarts.ResourceModelArtsAuthorization(),
			"hcso_modelarts_network":                modelarts.ResourceModelartsNetwork(),
			"hcso_modelarts_resource_pool":          modelarts.ResourceModelartsResourcePool(),

			"hcso_dataarts_studio_instance": dataarts.ResourceStudioInstance(),

			"hcso_mpc_transcoding_template":       mpc.ResourceTranscodingTemplate(),
			"hcso_mpc_transcoding_template_group": mpc.ResourceTranscodingTemplateGroup(),

			"hcso_mrs_cluster": ResourceMRSClusterV1(),
			"hcso_mrs_job":     ResourceMRSJobV1(),

			"hcso_nat_dnat_rule": nat.ResourcePublicDnatRule(),
			"hcso_nat_gateway":   nat.ResourcePublicGateway(),
			"hcso_nat_snat_rule": nat.ResourcePublicSnatRule(),

			"hcso_nat_private_dnat_rule":  nat.ResourcePrivateDnatRule(),
			"hcso_nat_private_gateway":    nat.ResourcePrivateGateway(),
			"hcso_nat_private_snat_rule":  nat.ResourcePrivateSnatRule(),
			"hcso_nat_private_transit_ip": nat.ResourcePrivateTransitIp(),

			"hcso_network_acl":              ResourceNetworkACL(),
			"hcso_network_acl_rule":         ResourceNetworkACLRule(),
			"hcso_networking_secgroup":      vpc.ResourceNetworkingSecGroup(),
			"hcso_networking_secgroup_rule": vpc.ResourceNetworkingSecGroupRule(),
			"hcso_networking_vip":           vpc.ResourceNetworkingVip(),
			"hcso_networking_vip_associate": vpc.ResourceNetworkingVIPAssociateV2(),

			"hcso_obs_bucket":             obs.ResourceObsBucket(),
			"hcso_obs_bucket_acl":         obs.ResourceOBSBucketAcl(),
			"hcso_obs_bucket_object":      obs.ResourceObsBucketObject(),
			"hcso_obs_bucket_object_acl":  obs.ResourceOBSBucketObjectAcl(),
			"hcso_obs_bucket_policy":      obs.ResourceObsBucketPolicy(),
			"hcso_obs_bucket_replication": obs.ResourceObsBucketReplication(),

			"hcso_oms_migration_task":       oms.ResourceMigrationTask(),
			"hcso_oms_migration_task_group": oms.ResourceMigrationTaskGroup(),

			"hcso_ram_resource_share": ram.ResourceRAMShare(),

			"hcso_rds_mysql_account":                rds.ResourceMysqlAccount(),
			"hcso_rds_mysql_database":               rds.ResourceMysqlDatabase(),
			"hcso_rds_mysql_database_privilege":     rds.ResourceMysqlDatabasePrivilege(),
			"hcso_rds_instance":                     rds.ResourceRdsInstance(),
			"hcso_rds_parametergroup":               rds.ResourceRdsConfiguration(),
			"hcso_rds_read_replica_instance":        rds.ResourceRdsReadReplicaInstance(),
			"hcso_rds_backup":                       rds.ResourceBackup(),
			"hcso_rds_cross_region_backup_strategy": rds.ResourceBackupStrategy(),
			"hcso_rds_sql_audit":                    rds.ResourceSQLAudit(),

			"hcso_rms_policy_assignment":                  rms.ResourcePolicyAssignment(),
			"hcso_rms_resource_aggregator":                rms.ResourceAggregator(),
			"hcso_rms_resource_aggregation_authorization": rms.ResourceAggregationAuthorization(),
			"hcso_rms_resource_recorder":                  rms.ResourceRecorder(),

			"hcso_sdrs_drill":              sdrs.ResourceDrill(),
			"hcso_sdrs_replication_pair":   sdrs.ResourceReplicationPair(),
			"hcso_sdrs_protection_group":   sdrs.ResourceProtectionGroup(),
			"hcso_sdrs_protected_instance": sdrs.ResourceProtectedInstance(),
			"hcso_sdrs_replication_attach": sdrs.ResourceReplicationAttach(),

			"hcso_secmaster_incident": secmaster.ResourceIncident(),

			"hcso_servicestage_application":                 servicestage.ResourceApplication(),
			"hcso_servicestage_component_instance":          servicestage.ResourceComponentInstance(),
			"hcso_servicestage_component":                   servicestage.ResourceComponent(),
			"hcso_servicestage_environment":                 servicestage.ResourceEnvironment(),
			"hcso_servicestage_repo_token_authorization":    servicestage.ResourceRepoTokenAuth(),
			"hcso_servicestage_repo_password_authorization": servicestage.ResourceRepoPwdAuth(),

			"hcso_sfs_access_rule": sfs.ResourceSFSAccessRuleV2(),
			"hcso_sfs_file_system": sfs.ResourceSFSFileSystemV2(),
			"hcso_sfs_turbo":       sfs.ResourceSFSTurbo(),

			"hcso_smn_topic":            smn.ResourceTopic(),
			"hcso_smn_subscription":     smn.ResourceSubscription(),
			"hcso_smn_message_template": smn.ResourceSmnMessageTemplate(),

			"hcso_sms_server_template": sms.ResourceServerTemplate(),
			"hcso_sms_task":            sms.ResourceMigrateTask(),

			"hcso_swr_organization":             swr.ResourceSWROrganization(),
			"hcso_swr_organization_permissions": swr.ResourceSWROrganizationPermissions(),
			"hcso_swr_repository":               swr.ResourceSWRRepository(),
			"hcso_swr_repository_sharing":       swr.ResourceSWRRepositorySharing(),
			"hcso_swr_image_permissions":        swr.ResourceSwrImagePermissions(),
			"hcso_swr_image_trigger":            swr.ResourceSwrImageTrigger(),
			"hcso_swr_image_retention_policy":   swr.ResourceSwrImageRetentionPolicy(),
			"hcso_swr_image_auto_sync":          swr.ResourceSwrImageAutoSync(),

			"hcso_tms_tags": tms.ResourceTmsTag(),

			"hcso_ucs_fleet":   ucs.ResourceFleet(),
			"hcso_ucs_cluster": ucs.ResourceCluster(),
			"hcso_ucs_policy":  ucs.ResourcePolicy(),

			"hcso_vod_media_asset":                vod.ResourceMediaAsset(),
			"hcso_vod_media_category":             vod.ResourceMediaCategory(),
			"hcso_vod_transcoding_template_group": vod.ResourceTranscodingTemplateGroup(),
			"hcso_vod_watermark_template":         vod.ResourceWatermarkTemplate(),

			"hcso_vpc_bandwidth":           eip.ResourceVpcBandWidthV2(),
			"hcso_vpc_bandwidth_associate": eip.ResourceBandWidthAssociate(),
			"hcso_vpc_eip":                 eip.ResourceVpcEIPV1(),
			"hcso_vpc_eip_associate":       eip.ResourceEIPAssociate(),

			"hcso_vpc_peering_connection":          vpc.ResourceVpcPeeringConnectionV2(),
			"hcso_vpc_peering_connection_accepter": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"hcso_vpc_route_table":                 vpc.ResourceVPCRouteTable(),
			"hcso_vpc_route":                       vpc.ResourceVPCRouteTableRoute(),
			"hcso_vpc":                             vpc.ResourceVirtualPrivateCloudV1(),
			"hcso_vpc_subnet":                      vpc.ResourceVpcSubnetV1(),
			"hcso_vpc_address_group":               vpc.ResourceVpcAddressGroup(),
			"hcso_vpc_flow_log":                    vpc.ResourceVpcFlowLog(),

			"hcso_vpcep_approval": vpcep.ResourceVPCEndpointApproval(),
			"hcso_vpcep_endpoint": vpcep.ResourceVPCEndpoint(),
			"hcso_vpcep_service":  vpcep.ResourceVPCEndpointService(),

			"hcso_vpn_gateway":                 vpn.ResourceGateway(),
			"hcso_vpn_customer_gateway":        vpn.ResourceCustomerGateway(),
			"hcso_vpn_connection":              vpn.ResourceConnection(),
			"hcso_vpn_connection_health_check": vpn.ResourceConnectionHealthCheck(),

			"hcso_scm_certificate": scm.ResourceScmCertificate(),

			"hcso_waf_address_group":                       waf.ResourceWafAddressGroup(),
			"hcso_waf_certificate":                         waf.ResourceWafCertificateV1(),
			"hcso_waf_cloud_instance":                      waf.ResourceCloudInstance(),
			"hcso_waf_domain":                              waf.ResourceWafDomainV1(),
			"hcso_waf_policy":                              waf.ResourceWafPolicyV1(),
			"hcso_waf_rule_anti_crawler":                   waf.ResourceRuleAntiCrawler(),
			"hcso_waf_rule_blacklist":                      waf.ResourceWafRuleBlackListV1(),
			"hcso_waf_rule_cc_protection":                  waf.ResourceRuleCCProtection(),
			"hcso_waf_rule_data_masking":                   waf.ResourceWafRuleDataMaskingV1(),
			"hcso_waf_rule_global_protection_whitelist":    waf.ResourceRuleGlobalProtectionWhitelist(),
			"hcso_waf_rule_known_attack_source":            waf.ResourceRuleKnownAttack(),
			"hcso_waf_rule_precise_protection":             waf.ResourceRulePreciseProtection(),
			"hcso_waf_rule_web_tamper_protection":          waf.ResourceWafRuleWebTamperProtectionV1(),
			"hcso_waf_rule_geolocation_access_control":     waf.ResourceRuleGeolocation(),
			"hcso_waf_rule_information_leakage_prevention": waf.ResourceRuleLeakagePrevention(),
			"hcso_waf_dedicated_instance":                  waf.ResourceWafDedicatedInstance(),
			"hcso_waf_dedicated_domain":                    waf.ResourceWafDedicatedDomainV1(),
			"hcso_waf_instance_group":                      waf.ResourceWafInstanceGroup(),
			"hcso_waf_instance_group_associate":            waf.ResourceWafInstGroupAssociate(),
			"hcso_waf_reference_table":                     waf.ResourceWafReferenceTableV1(),

			"hcso_workspace_desktop": workspace.ResourceDesktop(),
			"hcso_workspace_service": workspace.ResourceService(),
			"hcso_workspace_user":    workspace.ResourceUser(),

			"hcso_cpts_project": cpts.ResourceProject(),
			"hcso_cpts_task":    cpts.ResourceTask(),

			// CodeArts
			"hcso_codearts_project":            codearts.ResourceProject(),
			"hcso_codearts_repository":         codearts.ResourceRepository(),
			"hcso_codearts_deploy_application": codearts.ResourceDeployApplication(),
			"hcso_codearts_deploy_group":       codearts.ResourceDeployGroup(),
			"hcso_codearts_deploy_host":        codearts.ResourceDeployHost(),

			"hcso_dsc_instance":  dsc.ResourceDscInstance(),
			"hcso_dsc_asset_obs": dsc.ResourceAssetObs(),

			// internal only
			"hcso_apm_aksk":                apm.ResourceApmAkSk(),
			"hcso_aom_alarm_policy":        aom.ResourceAlarmPolicy(),
			"hcso_aom_prometheus_instance": aom.ResourcePrometheusInstance(),

			"hcso_aom_application":                 cmdb.ResourceAomApplication(),
			"hcso_aom_component":                   cmdb.ResourceAomComponent(),
			"hcso_aom_cmdb_resource_relationships": cmdb.ResourceCiRelationships(),
			"hcso_aom_environment":                 cmdb.ResourceAomEnvironment(),

			"hcso_lts_access_rule":     lts.ResourceAomMappingRule(),
			"hcso_lts_dashboard":       lts.ResourceLtsDashboard(),
			"hcso_elb_log":             lts.ResourceLtsElb(),
			"hcso_lts_struct_template": lts.ResourceLtsStruct(),

			// Legacy
			"hcso_networking_eip_associate": eip.ResourceEIPAssociate(),

			"hcso_projectman_project": codearts.ResourceProject(),
			"hcso_codehub_repository": codearts.ResourceRepository(),

			"hcso_compute_instance_v2":             ecs.ResourceComputeInstance(),
			"hcso_compute_interface_attach_v2":     ecs.ResourceComputeInterfaceAttach(),
			"hcso_compute_keypair_v2":              ResourceComputeKeypairV2(),
			"hcso_compute_servergroup_v2":          ecs.ResourceComputeServerGroup(),
			"hcso_compute_volume_attach_v2":        ecs.ResourceComputeVolumeAttach(),
			"hcso_compute_floatingip_associate_v2": ecs.ResourceComputeEIPAssociate(),

			"hcso_dns_ptrrecord_v2": dns.ResourceDNSPtrRecord(),
			"hcso_dns_recordset_v2": dns.ResourceDNSRecordSetV2(),
			"hcso_dns_zone_v2":      dns.ResourceDNSZone(),

			"hcso_dcs_instance_v1": dcs.ResourceDcsInstance(),
			"hcso_dds_instance_v3": dds.ResourceDdsInstanceV3(),

			"hcso_kms_key_v1": dew.ResourceKmsKey(),

			"hcso_lb_certificate_v2":  lb.ResourceCertificateV2(),
			"hcso_lb_loadbalancer_v2": lb.ResourceLoadBalancer(),
			"hcso_lb_listener_v2":     lb.ResourceListener(),
			"hcso_lb_pool_v2":         lb.ResourcePoolV2(),
			"hcso_lb_member_v2":       lb.ResourceMemberV2(),
			"hcso_lb_monitor_v2":      lb.ResourceMonitorV2(),
			"hcso_lb_l7policy_v2":     lb.ResourceL7PolicyV2(),
			"hcso_lb_l7rule_v2":       lb.ResourceL7RuleV2(),
			"hcso_lb_whitelist_v2":    lb.ResourceWhitelistV2(),

			"hcso_mrs_cluster_v1": ResourceMRSClusterV1(),
			"hcso_mrs_job_v1":     ResourceMRSJobV1(),

			"hcso_networking_secgroup_v2":      vpc.ResourceNetworkingSecGroup(),
			"hcso_networking_secgroup_rule_v2": vpc.ResourceNetworkingSecGroupRule(),

			"hcso_smn_topic_v2":        smn.ResourceTopic(),
			"hcso_smn_subscription_v2": smn.ResourceSubscription(),

			"hcso_rds_account":            rds.ResourceMysqlAccount(),
			"hcso_rds_database":           rds.ResourceMysqlDatabase(),
			"hcso_rds_database_privilege": rds.ResourceMysqlDatabasePrivilege(),
			"hcso_rds_instance_v3":        rds.ResourceRdsInstance(),
			"hcso_rds_parametergroup_v3":  rds.ResourceRdsConfiguration(),

			"hcso_rf_stack": rfs.ResourceStack(),

			"hcso_nat_dnat_rule_v2": nat.ResourcePublicDnatRule(),
			"hcso_nat_gateway_v2":   nat.ResourcePublicGateway(),
			"hcso_nat_snat_rule_v2": nat.ResourcePublicSnatRule(),

			"hcso_sfs_access_rule_v2": sfs.ResourceSFSAccessRuleV2(),
			"hcso_sfs_file_system_v2": sfs.ResourceSFSFileSystemV2(),

			"hcso_iam_agency":    iam.ResourceIAMAgencyV3(),
			"hcso_iam_agency_v3": iam.ResourceIAMAgencyV3(),

			"hcso_vpc_bandwidth_v2":                   eip.ResourceVpcBandWidthV2(),
			"hcso_vpc_eip_v1":                         eip.ResourceVpcEIPV1(),
			"hcso_vpc_peering_connection_v2":          vpc.ResourceVpcPeeringConnectionV2(),
			"hcso_vpc_peering_connection_accepter_v2": vpc.ResourceVpcPeeringConnectionAccepterV2(),
			"hcso_vpc_v1":                             vpc.ResourceVirtualPrivateCloudV1(),
			"hcso_vpc_subnet_v1":                      vpc.ResourceVpcSubnetV1(),

			"hcso_cce_cluster_v3": cce.ResourceCCEClusterV3(),
			"hcso_cce_node_v3":    cce.ResourceNode(),

			"hcso_as_configuration_v1": as.ResourceASConfiguration(),
			"hcso_as_group_v1":         as.ResourceASGroup(),
			"hcso_as_policy_v1":        as.ResourceASPolicy(),

			"hcso_identity_project_v3":          iam.ResourceIdentityProject(),
			"hcso_identity_role_assignment_v3":  iam.ResourceIdentityGroupRoleAssignment(),
			"hcso_identity_user_v3":             iam.ResourceIdentityUser(),
			"hcso_identity_group_v3":            iam.ResourceIdentityGroup(),
			"hcso_identity_group_membership_v3": iam.ResourceIdentityGroupMembership(),
			"hcso_identity_provider_conversion": iam.ResourceIAMProviderConversion(),

			"hcso_cdm_cluster_v1": cdm.ResourceCdmCluster(),
			"hcso_css_cluster_v1": css.ResourceCssCluster(),
			"hcso_dis_stream_v2":  dis.ResourceDisStream(),

			"hcso_organizations_organization":            organizations.ResourceOrganization(),
			"hcso_organizations_organizational_unit":     organizations.ResourceOrganizationalUnit(),
			"hcso_organizations_account":                 organizations.ResourceAccount(),
			"hcso_organizations_account_associate":       organizations.ResourceAccountAssociate(),
			"hcso_organizations_account_invite":          organizations.ResourceAccountInvite(),
			"hcso_organizations_account_invite_accepter": organizations.ResourceAccountInviteAccepter(),
			"hcso_organizations_trusted_service":         organizations.ResourceTrustedService(),

			"hcso_dli_queue_v1":                dli.ResourceDliQueue(),
			"hcso_networking_vip_v2":           vpc.ResourceNetworkingVip(),
			"hcso_networking_vip_associate_v2": vpc.ResourceNetworkingVIPAssociateV2(),
			"hcso_fgs_function_v2":             fgs.ResourceFgsFunctionV2(),
			"hcso_cdn_domain_v1":               resourceCdnDomainV1(),

			// Deprecated
			"hcso_apig_vpc_channel":               deprecated.ResourceApigVpcChannelV2(),
			"hcso_blockstorage_volume_v2":         deprecated.ResourceBlockStorageVolumeV2(),
			"hcso_csbs_backup":                    deprecated.ResourceCSBSBackupV1(),
			"hcso_csbs_backup_policy":             deprecated.ResourceCSBSBackupPolicyV1(),
			"hcso_csbs_backup_policy_v1":          deprecated.ResourceCSBSBackupPolicyV1(),
			"hcso_csbs_backup_v1":                 deprecated.ResourceCSBSBackupV1(),
			"hcso_networking_network_v2":          deprecated.ResourceNetworkingNetworkV2(),
			"hcso_networking_subnet_v2":           deprecated.ResourceNetworkingSubnetV2(),
			"hcso_networking_floatingip_v2":       deprecated.ResourceNetworkingFloatingIPV2(),
			"hcso_networking_router_v2":           deprecated.ResourceNetworkingRouterV2(),
			"hcso_networking_router_interface_v2": deprecated.ResourceNetworkingRouterInterfaceV2(),
			"hcso_networking_router_route_v2":     deprecated.ResourceNetworkingRouterRouteV2(),
			"hcso_networking_port":                deprecated.ResourceNetworkingPortV2(),
			"hcso_networking_port_v2":             deprecated.ResourceNetworkingPortV2(),
			"hcso_vpc_route_v2":                   deprecated.ResourceVPCRouteV2(),
			"hcso_ecs_instance_v1":                deprecated.ResourceEcsInstanceV1(),
			"hcso_compute_secgroup_v2":            deprecated.ResourceComputeSecGroupV2(),
			"hcso_compute_floatingip_v2":          deprecated.ResourceComputeFloatingIPV2(),
			"hcso_oms_task":                       deprecated.ResourceMaasTaskV1(),

			"hcso_fw_firewall_group_v2": deprecated.ResourceFWFirewallGroupV2(),
			"hcso_fw_policy_v2":         deprecated.ResourceFWPolicyV2(),
			"hcso_fw_rule_v2":           deprecated.ResourceFWRuleV2(),

			"hcso_images_image_v2": deprecated.ResourceImagesImageV2(),

			"hcso_dms_instance":    deprecated.ResourceDmsInstancesV1(),
			"hcso_dms_instance_v1": deprecated.ResourceDmsInstancesV1(),
			"hcso_dms_group":       deprecated.ResourceDmsGroups(),
			"hcso_dms_group_v1":    deprecated.ResourceDmsGroups(),
			"hcso_dms_queue":       deprecated.ResourceDmsQueues(),
			"hcso_dms_queue_v1":    deprecated.ResourceDmsQueues(),

			"hcso_cs_cluster":            deprecated.ResourceCsClusterV1(),
			"hcso_cs_cluster_v1":         deprecated.ResourceCsClusterV1(),
			"hcso_cs_route":              deprecated.ResourceCsRouteV1(),
			"hcso_cs_route_v1":           deprecated.ResourceCsRouteV1(),
			"hcso_cs_peering_connect":    deprecated.ResourceCsPeeringConnectV1(),
			"hcso_cs_peering_connect_v1": deprecated.ResourceCsPeeringConnectV1(),

			"hcso_vbs_backup":           deprecated.ResourceVBSBackupV2(),
			"hcso_vbs_backup_policy":    deprecated.ResourceVBSBackupPolicyV2(),
			"hcso_vbs_backup_policy_v2": deprecated.ResourceVBSBackupPolicyV2(),
			"hcso_vbs_backup_v2":        deprecated.ResourceVBSBackupV2(),

			"hcso_vpnaas_ipsec_policy_v2":    deprecated.ResourceVpnIPSecPolicyV2(),
			"hcso_vpnaas_service_v2":         deprecated.ResourceVpnServiceV2(),
			"hcso_vpnaas_ike_policy_v2":      deprecated.ResourceVpnIKEPolicyV2(),
			"hcso_vpnaas_endpoint_group_v2":  deprecated.ResourceVpnEndpointGroupV2(),
			"hcso_vpnaas_site_connection_v2": deprecated.ResourceVpnSiteConnectionV2(),

			"hcso_vpnaas_endpoint_group":  deprecated.ResourceVpnEndpointGroupV2(),
			"hcso_vpnaas_ike_policy":      deprecated.ResourceVpnIKEPolicyV2(),
			"hcso_vpnaas_ipsec_policy":    deprecated.ResourceVpnIPSecPolicyV2(),
			"hcso_vpnaas_service":         deprecated.ResourceVpnServiceV2(),
			"hcso_vpnaas_site_connection": deprecated.ResourceVpnSiteConnectionV2(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11 cc
			terraformVersion = "0.11+compatible"
		}

		return configureProvider(ctx, d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"region": "The HuaweiCloud region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"project_id": "The ID of the project to login with.",

		"project_name": "The name of the project to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"domain_id": "The ID of the Domain to scope to.",

		"domain_name": "The name of the Domain to scope to.",

		"access_key":     "The access key of the HuaweiCloud to use.",
		"secret_key":     "The secret key of the HuaweiCloud to use.",
		"security_token": "The security token to authenticate with a temporary security credential.",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"agency_name": "The name of agency",

		"agency_domain_name": "The name of domain who created the agency (Identity v3).",

		"delegated_project": "The name of delegated project (Identity v3).",

		"assume_role_agency_name": "The name of agency for assume role.",

		"assume_role_domain_name": "The name of domain for assume role.",

		"cloud": "The endpoint of cloud provider, defaults to myhuaweicloud.com",

		"endpoints": "The custom endpoints used to override the default endpoint URL.",

		"regional": "Whether the service endpoints are regional",

		"shared_config_file": "The path to the shared config file. If not set, the default is ~/.hcloud/config.json.",

		"profile": "The profile name as set in the shared config file.",

		"max_retries": "How many times HTTP connection should be retried until giving up.",

		"enterprise_project_id": "enterprise project id",
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData, terraformVersion string) (interface{},
	diag.Diagnostics) {
	var tenantName, tenantID, delegatedProject, identityEndpoint string
	region := d.Get("region").(string)
	isRegional := d.Get("regional").(bool)
	cloud := getCloudDomain(d.Get("cloud").(string), region)

	// project_name is prior to tenant_name
	// if neither of them was set, use region as the default project
	if v, ok := d.GetOk("project_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else if v, ok := d.GetOk("tenant_name"); ok && v.(string) != "" {
		tenantName = v.(string)
	} else {
		tenantName = region
	}

	// project_id is prior to tenant_id
	if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		tenantID = v.(string)
	} else {
		tenantID = d.Get("tenant_id").(string)
	}

	// Use region as delegated_project if it's not set
	if v, ok := d.GetOk("delegated_project"); ok && v.(string) != "" {
		delegatedProject = v.(string)
	} else {
		delegatedProject = region
	}

	// use auth_url as identityEndpoint if specified
	if v, ok := d.GetOk("auth_url"); ok {
		identityEndpoint = v.(string)
	} else {
		// use cloud as basis for identityEndpoint
		identityEndpoint = fmt.Sprintf("https://iam.%s.%s/v3", region, cloud)
	}

	hcsoConfig := hcso_config.HCSOConfig{
		Config: config.Config{
			AccessKey:           d.Get("access_key").(string),
			SecretKey:           d.Get("secret_key").(string),
			CACertFile:          d.Get("cacert_file").(string),
			ClientCertFile:      d.Get("cert").(string),
			ClientKeyFile:       d.Get("key").(string),
			DomainID:            d.Get("domain_id").(string),
			DomainName:          d.Get("domain_name").(string),
			IdentityEndpoint:    identityEndpoint,
			Insecure:            d.Get("insecure").(bool),
			Password:            d.Get("password").(string),
			Token:               d.Get("token").(string),
			SecurityToken:       d.Get("security_token").(string),
			Region:              region,
			TenantID:            tenantID,
			TenantName:          tenantName,
			Username:            d.Get("user_name").(string),
			UserID:              d.Get("user_id").(string),
			AgencyName:          d.Get("agency_name").(string),
			AgencyDomainName:    d.Get("agency_domain_name").(string),
			DelegatedProject:    delegatedProject,
			Cloud:               cloud,
			RegionClient:        isRegional,
			MaxRetries:          d.Get("max_retries").(int),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
			SharedConfigFile:    d.Get("shared_config_file").(string),
			Profile:             d.Get("profile").(string),
			TerraformVersion:    terraformVersion,
			RegionProjectIDMap:  make(map[string]string),
			RPLock:              new(sync.Mutex),
			SecurityKeyLock:     new(sync.Mutex),
		},
	}

	hcsoConfig.Metadata = &hcsoConfig.Config

	// get assume role
	assumeRoleList := d.Get("assume_role").([]interface{})
	if len(assumeRoleList) == 0 {
		// without assume_role block in provider
		delegatedAgencyName := os.Getenv("HCSO_ASSUME_ROLE_AGENCY_NAME")
		delegatedDomianName := os.Getenv("HCSO_ASSUME_ROLE_DOMAIN_NAME")
		if delegatedAgencyName != "" && delegatedDomianName != "" {
			hcsoConfig.AssumeRoleAgency = delegatedAgencyName
			hcsoConfig.AssumeRoleDomain = delegatedDomianName
		}
	} else {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		hcsoConfig.AssumeRoleAgency = assumeRole["agency_name"].(string)
		hcsoConfig.AssumeRoleDomain = assumeRole["domain_name"].(string)
	}

	// get custom endpoints
	endpoints, err := flattenProviderEndpoints(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	hcsoConfig.Endpoints = endpoints

	if err := hcsoConfig.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &hcsoConfig.Config, nil
}

func flattenProviderEndpoints(d *schema.ResourceData) (map[string]string, error) {
	endpoints := d.Get("endpoints").(map[string]interface{})
	epMap := make(map[string]string)

	for key, val := range endpoints {
		endpoint := strings.TrimSpace(val.(string))
		// check empty string
		if endpoint == "" {
			return nil, fmt.Errorf("the value of customer endpoint %s must be specified", key)
		}

		// add prefix "https://" and suffix "/"
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = fmt.Sprintf("%s/", endpoint)
		}
		epMap[key] = endpoint
	}

	// unify the endpoint which has multiple versions
	for key := range endpoints {
		ep, ok := epMap[key]
		if !ok {
			continue
		}

		multiKeys := config.GetServiceDerivedCatalogKeys(key)
		for _, k := range multiKeys {
			epMap[k] = ep
		}
	}

	log.Printf("[DEBUG] customer endpoints: %+v", epMap)
	return epMap, nil
}

func getCloudDomain(cloud, region string) string {
	// first, use the specified value
	if cloud != "" {
		return cloud
	}

	// then check whether the region(eu-west-1xx) is located in Europe
	if strings.HasPrefix(region, "eu-west-1") {
		return defaultEuropeCloud
	}
	return defaultCloud
}
