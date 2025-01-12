package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/nealbrown/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func shIndexesManager() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datatype": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"event", "metric"}, false),
				Description:  `Valid values: (event | metric). Specifies the type of index.`,
			},
			"frozen_time_period_in_secs": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "94608000",
				Description: `Number of seconds after which indexed data rolls to frozen. Defaults to 188697600 (6 years).
				Freezing data means it is removed from the index. If you need to archive your data, refer to coldToFrozenDir and coldToFrozenScript parameter documentation.`,
			},
			"max_global_raw_data_size_mb": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "100",
				Description: `The maximum size of an index (in MB). If an index grows larger than the maximum size, the oldest data is frozen.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the index to create.`,
			},
			"acl": aclSchema(),
		},
		Read:   shIndexesManagerRead,
		Create: shIndexesManagerCreate,
		Delete: shIndexesManagerDelete,
		Update: shIndexesManagerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func shIndexesManagerCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	indexConfigObj := getShIndexesManagerConfig(d)
	aclObject := &models.ACLObject{
		Owner: "nobody",
		App:   "cloud_administration",
	}
	err := (*provider.Client).CreateShIndexesManagerObject(name, aclObject.Owner, aclObject.App, indexConfigObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return shIndexesManagerRead(d, meta)
}

func shIndexesManagerRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of indexes to get owner and app name for the specific index
	var err error
	var resp *http.Response
	var entry *models.ShIndexesManagerEntry
	// Retrying due to delay in backend to create resource
	for i := 0; i < 10; i++ {
		resp, err = (*provider.Client).ReadAllShIndexesManagerObject()
		if err != nil {
			continue
		}

		entry, err = getShIndexesManagerConfigByName(name, resp)
		if err != nil {
			continue
		}

		_ = resp.Body.Close()
		if entry == nil {
			err = fmt.Errorf("Unable to find resource: %s", name)
		} else {
			break
		}
	}

	if entry == nil {
		return err
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadShIndexesManagerObject(name, "nobody", "cloud_administration")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getShIndexesManagerConfigByName(name, resp)
	if err != nil {
		return err
	}

	if err = d.Set("datatype", entry.Content.Datatype); err != nil {
		return err
	}
	if err = d.Set("frozen_time_period_in_secs", entry.Content.FrozenTimePeriodInSecs); err != nil {
		return err
	}
	if err = d.Set("max_global_raw_data_size_mb", entry.Content.MaxGlobalRawDataSizeMB); err != nil {
		return err
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	return nil
}

func shIndexesManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	indexConfigObj := getShIndexesManagerConfig(d)
	aclObject := &models.ACLObject{
		Owner: "nobody",
		App:   "cloud_administration",
	}
	err := (*provider.Client).UpdateShIndexesManagerObject(d.Id(), aclObject.Owner, aclObject.App, indexConfigObj)
	if err != nil {
		return err
	}

	return shIndexesManagerRead(d, meta)
}

func shIndexesManagerDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := &models.ACLObject{
		Owner: "nobody",
		App:   "cloud_administration",
	}
	resp, err := (*provider.Client).DeleteShIndexesManagerObject(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.ShIndexesManagerResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getShIndexesManagerConfig(d *schema.ResourceData) (indexConfigObject *models.ShIndexesManagerObject) {
	indexConfigObject = &models.ShIndexesManagerObject{}
	indexConfigObject.Datatype = d.Get("datatype").(string)
	indexConfigObject.FrozenTimePeriodInSecs = d.Get("frozen_time_period_in_secs").(string)
	indexConfigObject.MaxGlobalRawDataSizeMB = d.Get("max_global_raw_data_size_mb").(string)
	return indexConfigObject
}

func getShIndexesManagerConfigByName(name string, httpResponse *http.Response) (indexEntry *models.ShIndexesManagerEntry, err error) {
	response := &models.ShIndexesManagerResponse{}
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
		return indexEntry, err
	}

	return indexEntry, nil
}
