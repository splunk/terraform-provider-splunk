package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func splunkLookupDefinitions() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Lookup name",
			},
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the static lookup table file",
			},
			"acl": aclSchema(),
		},
		Read:   splunkLookupDefinitionsRead,
		Create: splunkLookupDefinitionsCreate,
		Delete: splunkLookupDefinitionsDelete,
		Update: splunkLookupDefinitionsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func splunkLookupDefinitionsCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	splunkLookupDefinitionObj := getSplunkLookupDefinitionConfig(d)
	aclObject := getResourceDataViewACL(d)
	name := d.Get("name").(string)

	err := (*provider.Client).CreateLookupDefinitionObject(aclObject.Owner, aclObject.App, splunkLookupDefinitionObj)
	if err != nil {
		return err
	}
	d.SetId(name)
	return splunkLookupDefinitionsRead(d, meta)
}

func splunkLookupDefinitionsRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getResourceDataViewACL(d)
	name := d.Id()
	readUser := aclObject.Owner

	resp, err := (*provider.Client).ReadLookupDefinitionObject(name, readUser, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getLookupDefintionByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("unable to find resource: %v", name)
	}

	if err = d.Set("name", entry.Name); err != nil {
		return err
	}

	if err = d.Set("filename", entry.Content.Filename); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func splunkLookupDefinitionsUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	splunkLookupDefinitionObj := getSplunkLookupDefinitionConfig(d)
	aclObject := getResourceDataViewACL(d)

	updateUser := "nobody"

	if aclObject.Sharing == "user" {
		// If we have a private dashboard we can only update it using the owner
		updateUser = aclObject.Owner
	}

	if err := (*provider.Client).UpdateLookupDefinitionObject(updateUser, aclObject.App, name, splunkLookupDefinitionObj); err != nil {
		return err
	}

	return nil
}

func splunkLookupDefinitionsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	aclObject := getResourceDataViewACL(d)
	if aclObject.Sharing != "user" {
		aclObject.Owner = "nobody"
	}
	resp, err := (*provider.Client).DeleteLookupDefinitionObject(aclObject.Owner, aclObject.App, name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.SplunkLookupDefinitionResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

func getLookupDefintionByName(name string, httpResponse *http.Response) (lookupDefinitionEntry *models.SplunkLookupDefinitionEntry, err error) {
	response := &models.SplunkLookupDefinitionResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:

		decoder := json.NewDecoder(httpResponse.Body)
		err := decoder.Decode(response)
		if err != nil {
			return nil, err
		}
		re := regexp.MustCompile(`(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return lookupDefinitionEntry, err
	}

	return lookupDefinitionEntry, nil
}

func getSplunkLookupDefinitionConfig(d *schema.ResourceData) (splunkLookupDefinitionObject *models.SplunkLookupDefinitionObject) {
	splunkLookupDefinitionObject = &models.SplunkLookupDefinitionObject{}
	splunkLookupDefinitionObject.Name = d.Get("name").(string)
	splunkLookupDefinitionObject.Filename = d.Get("filename").(string)
	return splunkLookupDefinitionObject
}
