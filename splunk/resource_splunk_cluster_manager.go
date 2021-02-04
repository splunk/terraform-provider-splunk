package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func clusterManagerConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"available_sites": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Sets the various sites that are recognized for this master. Valid values include site1 to site64.`,
			},
			"cluster_label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Label for this cluster.`,
			},
			"cxn_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Low-level timeout, in seconds, for establishing connection between cluster nodes`,
			},
			"heartbeat_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Only valid for the master node in a cluster configuration. Time, in seconds, before a master considers a peer down. Once a peer is down, the master initiates steps to replicate buckets from the dead peer to its live peers. Defaults to 60 seconds`,
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"master", "slave", "searchhead", "disabled"}, true),
				Description:  `Required. Valid values: (master | slave | searchhead | disabled) Defaults to disabled. Sets operational mode for this cluster node. Only one master may exist per cluster.`,
			},
			"multisite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Enable or disable the multisite feature for this cluster.`,
			},
			"replication_factor": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Only valid for nodes configured as a master. Determines how many copies of raw data are created in the cluster. This could be less than the number of cluster peers. Must be greater than 0 and greater than or equal to the search factor. Defaults to 3.`,
			},
			"restart_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Only valid for nodes configured as a master. The amount of time, in seconds, the master waits for a peer to come back when the peer is restarted (to avoid the overhead of trying to fix the buckets that were on the peer). Defaults to 600 seconds.`,
			},
			"search_factor": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Only valid for nodes configured as a master. Determines how many searchable copies of each bucket to maintain. Must be less than or equal to replication_factor and greater than 0. Defaults to 2.`,
			},
			"use_batch_mask_changes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies if the master should process bucket mask changes in batch or individually one by one. Defaults to true. Set to false when there are 6.1 peers in the cluster for backwards compatibility.`,
			},
			"acl": aclSchema(),
		},
		Read:   clusterManagerRead,
		Create: clusterManagerCreate,
		Delete: clusterManagerDelete,
		Update: clusterManagerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func clusterManagerCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("cluster_label").(string)
	clusterManagerObj := getClusterManagerConfig(d)
	err := (*provider.Client).CreateClusterManager(name, clusterManagerObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return clusterManagerRead(d, meta)
}

func clusterManagerRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	//
	resp, err := (*provider.Client).ReadClusterManager()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getClusterManagerConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	if err = d.Set("available_sites", entry.Content.AvailableSites); err != nil {
		return err
	}

	if err = d.Set("cluster_label", entry.Content.ClusterLabel); err != nil {
		return err
	}

	if err = d.Set("cxn_timeout", entry.Content.ConnectionTimeout); err != nil {
		return err
	}

	if err = d.Set("heartbeat_timeout", entry.Content.HeartbeatTimeout); err != nil {
		return err
	}
	if err = d.Set("mode", entry.Content.Mode); err != nil {
		return err
	}

	if err = d.Set("multisite", entry.Content.Multisite); err != nil {
		return err
	}

	if err = d.Set("replication_factor", entry.Content.ReplicationFactor); err != nil {
		return err
	}

	if err = d.Set("restart_timeout", entry.Content.RestartTimeout); err != nil {
		return err
	}

	if err = d.Set("search_factor", entry.Content.SearchFactor); err != nil {
		return err
	}

	if err = d.Set("use_batch_mask_changes", entry.Content.UseBatchMaskChanges); err != nil {
		return err
	}

	return nil
}

func clusterManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	clusterManagerObject := getClusterManagerConfig(d)
	err := (*provider.Client).UpdateClusterManager(d.Id(), clusterManagerObject)
	if err != nil {
		return err
	}

	return clusterManagerRead(d, meta)
}

func clusterManagerDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteClusterManager(d.Id())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.ClusterManagerResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getClusterManagerConfig(d *schema.ResourceData) (clusterManagerObject *models.ClusterManagerObject) {
	clusterManagerObject = &models.ClusterManagerObject{}
	clusterManagerObject.AvailableSites = d.Get("available_sites").(string)
	clusterManagerObject.ClusterLabel = d.Get("cluster_label").(string)
	clusterManagerObject.ConnectionTimeout = d.Get("cxn_timeout").(int)
	clusterManagerObject.HeartbeatTimeout = d.Get("heartbeat_timeout").(int)
	clusterManagerObject.Mode = d.Get("mode").(string)
	clusterManagerObject.Multisite = d.Get("multisite").(bool)
	clusterManagerObject.ReplicationFactor = d.Get("replication_factor").(int)
	clusterManagerObject.RestartTimeout = d.Get("restart_timeout").(int)
	clusterManagerObject.SearchFactor = d.Get("search_factor").(int)
	clusterManagerObject.UseBatchMaskChanges = d.Get("use_batch_mask_changes").(bool)
	return clusterManagerObject
}

func getClusterManagerConfigByName(name string, httpResponse *http.Response) (ClusterManagerEntry *models.ClusterManagerEntry, err error) {
	response := &models.ClusterManagerResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		re := regexp.MustCompile(`(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return ClusterManagerEntry, err
	}

	return ClusterManagerEntry, nil
}
