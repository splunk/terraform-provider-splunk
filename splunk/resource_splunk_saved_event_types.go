package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

func savedEventTypes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Saved event type name",
			},
			"search": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search terms for this event type.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description of this event type",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If True, disables the event type",
			},
			"color": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Color for this event type. The supported colors are: none, et_blue, et_green, et_magenta, et_orange, et_purple, et_red, et_sky, et_teal, et_yellow.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the order in which matching event types are displayed for an event. 1 is the highest, and 10 is the lowest.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "[Deprecated] Use tags.conf.spec file to assign tags to groups of events with related field values. ",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"acl": aclSchema(),
		},
		Create: savedEventTypesCreate,
		Read:   savedEventTypesRead,
		Update: savedEventTypesUpdate,
		Delete: savedEventTypesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func savedEventTypesCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	savedEventTypeConfig := getSavedEventTypeConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.App = "search"
		aclObject.Owner = "nobody"
	}
	err := (*provider.Client).CreateSavedEventTypes(name, aclObject.Owner, aclObject.App, savedEventTypeConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "saved", "eventtypes")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return savedEventTypesRead(d, meta)
}

func savedEventTypesRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)

	name := d.Id()
	// We first get list of searches to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadAllSavedEventTypes()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getSavedEventTypeConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the configuration with proper owner and app
	resp, err = (*provider.Client).ReadSavedEventTypes(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getSavedEventTypeConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}
	if err = d.Set("description", entry.Content.Description); err != nil {
		return err
	}
	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}
	if err = d.Set("color", entry.Content.Color); err != nil {
		return err
	}
	if err = d.Set("priority", entry.Content.Priority); err != nil {
		return err
	}
	if err = d.Set("search", entry.Content.Search); err != nil {
		return err
	}

	if err = d.Set("tags", entry.Content.Tags); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func savedEventTypesUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	savedEventTypesConfig := getSavedEventTypeConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))

	// Update will create a new resource with private `user` permissions if resource had shared permissions set
	var owner string
	if aclObject.Sharing != "user" {
		owner = "nobody"
	} else {
		owner = aclObject.Owner
	}

	err := (*provider.Client).UpdateSavedEventTypes(d.Id(), owner, aclObject.App, savedEventTypesConfig)
	if err != nil {
		return err
	}

	// Update ACL
	err = (*provider.Client).UpdateAcl(owner, aclObject.App, d.Id(), aclObject, "saved", "eventtypes")
	if err != nil {
		return err
	}

	return savedEventTypesRead(d, meta)
}

func savedEventTypesDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteSavedEventTypes(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.InputsUDPResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

func getSavedEventTypeConfig(d *schema.ResourceData) (savedEventTypeObj *models.SavedEventTypeObject) {
	savedEventTypeObj = &models.SavedEventTypeObject{
		Description: d.Get("description").(string),
		Disabled:    d.Get("disabled").(bool),
		Color:       d.Get("color").(string),
		Priority:    d.Get("priority").(int),
		Search:      d.Get("search").(string),
	}

	if val, ok := d.GetOk("tags"); ok {
		for _, v := range val.([]interface{}) {
			savedEventTypeObj.Tags = append(savedEventTypeObj.Tags, v.(string))
		}
	}

	return savedEventTypeObj
}

func getSavedEventTypeConfigByName(name string, httpResponse *http.Response) (savedEventTypesEntry *models.SavedEventTypesEntry, err error) {
	response := &models.SavedEventTypesResponse{}
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
		return savedEventTypesEntry, err
	}

	return savedEventTypesEntry, nil
}
