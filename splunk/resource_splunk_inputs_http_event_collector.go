package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/splunk/terraform-provider-splunk/client/models"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func inputsHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Token name (inputs.conf key)",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Token value for sending data to collector/event endpoint.",
			},
			"index": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Index to store generated events",
			},
			"indexes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of indexes allowed for events with this token.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default host value for events with this token",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default source for events with this token",
			},
			"sourcetype": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Default sourcetype for events with this token",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Input disabled indicator: 0 = Input Not disabled, 1 = Input disabled",
			},
			"use_ack": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Indexer acknowledgement for this token: 0 = disabled, 1 = enabled",
			},
			"acl": aclSchema(),
		},
		Read:   httpEventCollectorInputRead,
		Create: httpEventCollectorInputCreate,
		Delete: httpEventCollectorInputDelete,
		Update: httpEventCollectorInputUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func httpEventCollectorInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	httpInputConfigObj := getHttpEventCollectorConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.App = "splunk_httpinput"
		aclObject.Owner = "nobody"
		aclObject.Sharing = "app"
	}
	err := (*provider.Client).CreateHttpEventCollectorObject(name, aclObject.Owner, aclObject.App, httpInputConfigObj)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "http")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return httpEventCollectorInputRead(d, meta)
}

func httpEventCollectorInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadAllHttpEventCollectorObject()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getHECConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadHttpEventCollectorObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getHECConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("host", entry.Content.Host); err != nil {
		return err
	}

	if err = d.Set("token", entry.Content.Token); err != nil {
		return err
	}

	if err = d.Set("index", entry.Content.Index); err != nil {
		return err
	}

	if err = d.Set("indexes", entry.Content.Indexes); err != nil {
		return err
	}

	if err = d.Set("source", entry.Content.Source); err != nil {
		return err
	}

	if err = d.Set("sourcetype", entry.Content.SourceType); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("use_ack", entry.Content.UseACK); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func httpEventCollectorInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	httpInputConfigObj := getHttpEventCollectorConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateHttpEventCollectorObject(d.Id(), aclObject.Owner, aclObject.App, httpInputConfigObj)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "inputs", "http")
	if err != nil {
		return err
	}

	return httpEventCollectorInputRead(d, meta)
}

func httpEventCollectorInputDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteHttpEventCollectorObject(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.HECResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getHttpEventCollectorConfig(d *schema.ResourceData) (httpInputConfigObject *models.HttpEventCollectorObject) {
	httpInputConfigObject = &models.HttpEventCollectorObject{}
	httpInputConfigObject.Host = d.Get("host").(string)
	httpInputConfigObject.Index = d.Get("index").(string)
	httpInputConfigObject.Token = d.Get("token").(string)
	httpInputConfigObject.Indexes = d.Get("indexes").([]interface{})
	httpInputConfigObject.Source = d.Get("source").(string)
	httpInputConfigObject.SourceType = d.Get("sourcetype").(string)
	httpInputConfigObject.UseACK = d.Get("use_ack").(int)
	httpInputConfigObject.Disabled = d.Get("disabled").(bool)
	return httpInputConfigObject
}

func getHECConfigByName(name string, httpResponse *http.Response) (hecEntry *models.HECEntry, err error) {
	response := &models.HECResponse{}
	//body, err := ioutil.ReadAll(httpResponse.Body)
	//fmt.Println(body)
	switch httpResponse.StatusCode {
	case 200, 201:

		decoder := json.NewDecoder(httpResponse.Body)
		err := decoder.Decode(response)
		if err != nil {
			return nil, err
		}
		re := regexp.MustCompile(`http://(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return hecEntry, err
	}

	return hecEntry, nil
}
