package splunk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

var federatedIndexNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

func validateFederatedIndexName(value interface{}, key string) ([]string, []error) {
	name, ok := value.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%s must be a string", key)}
	}
	if len(name) > 2048 {
		return nil, []error{fmt.Errorf("%s must not exceed 2048 characters", key)}
	}
	if strings.Contains(name, "kvstore") {
		return nil, []error{fmt.Errorf("%s must not contain kvstore", key)}
	}
	if !federatedIndexNamePattern.MatchString(name) {
		return nil, []error{fmt.Errorf("%s must begin with a lowercase letter or number and contain only lowercase letters, numbers, underscores, or hyphens; omit the federated: prefix", key)}
	}
	return nil, nil
}

// federatedIndexes manages local indexes that map to remote provider datasets.
func federatedIndexes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateFederatedIndexName,
				Description:  "Unique local federated index name without the federated: prefix, for example remote-main.",
			},
			"federated_provider": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Name of the federated provider containing the remote dataset.",
			},
			"dataset_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "index",
				ValidateFunc: validation.StringInSlice([]string{"index", "savedsearch", "lastjob", "datamodel"}, false),
				Description:  "Remote dataset type: index, savedsearch, lastjob, or datamodel.",
			},
			"dataset_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Remote dataset name without its type prefix, for example main.",
			},
		},
		Create: federatedIndexCreate,
		Read:   federatedIndexRead,
		Update: federatedIndexUpdate,
		Delete: federatedIndexDelete,
		Importer: &schema.ResourceImporter{
			State: federatedIndexImport,
		},
	}
}

func federatedIndexCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	if err := (*provider.Client).CreateFederatedIndex(name, federatedIndexConfig(d)); err != nil {
		return err
	}

	// Splunk addresses federated indexes by their prefixed name in REST paths.
	d.SetId("federated:" + name)
	return federatedIndexRead(d, meta)
}

func federatedIndexRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).ReadFederatedIndex(d.Id())
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	entry, err := federatedIndexEntry(resp)
	if err != nil {
		return err
	}
	if entry == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", strings.TrimPrefix(entry.Name, "federated:")); err != nil {
		return err
	}
	if err := d.Set("federated_provider", entry.Content.FederatedProvider); err != nil {
		return err
	}
	// Split only the type prefix, preserving colons within remote object names.
	dataset := strings.SplitN(entry.Content.FederatedDataset, ":", 2)
	if len(dataset) != 2 || dataset[1] == "" {
		return fmt.Errorf("federated index %q returned an invalid dataset mapping", d.Id())
	}
	if err := d.Set("dataset_type", dataset[0]); err != nil {
		return err
	}
	return d.Set("dataset_name", dataset[1])
}

// federatedIndexImport accepts either a short name or Splunk's full index ID.
func federatedIndexImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	name := strings.TrimPrefix(d.Id(), "federated:")
	if _, errs := validateFederatedIndexName(name, "id"); len(errs) > 0 {
		return nil, errs[0]
	}
	d.SetId("federated:" + name)
	return []*schema.ResourceData{d}, nil
}

func federatedIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	if err := (*provider.Client).UpdateFederatedIndex(d.Id(), federatedIndexConfig(d)); err != nil {
		return err
	}
	return federatedIndexRead(d, meta)
}

func federatedIndexDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteFederatedIndex(d.Id())
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil && (resp == nil || resp.StatusCode != http.StatusNotFound) {
		return err
	}
	d.SetId("")
	return nil
}

// federatedIndexConfig combines the separate Terraform fields for Splunk's API.
func federatedIndexConfig(d *schema.ResourceData) *models.FederatedIndexObject {
	return &models.FederatedIndexObject{
		FederatedProvider: d.Get("federated_provider").(string),
		FederatedDataset:  d.Get("dataset_type").(string) + ":" + d.Get("dataset_name").(string),
	}
}

func federatedIndexEntry(resp *http.Response) (*models.FederatedIndexEntry, error) {
	response := &models.FederatedIndexResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	if len(response.Entry) == 0 {
		return nil, nil
	}
	if len(response.Entry) != 1 {
		return nil, fmt.Errorf("expected one federated index response entry, got %d", len(response.Entry))
	}
	return &response.Entry[0], nil
}
