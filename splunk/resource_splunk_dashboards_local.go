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

func splunkDashboards() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "App context for the dashboard.",
			},
			"eai_data": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Dashboard XML definition.",
			},
			"acl": aclSchema(),
		},
		Read:   splunkDashboardsRead,
		Create: splunkDashboardsCreate,
		Delete: splunkDashboardsDelete,
		Update: splunkDashboardsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func splunkDashboardsCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	splunkDashboardsObj := getSplunkDashboardsConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.App = "search"
		aclObject.Owner = "nobody"
		aclObject.Sharing = "app"
	}
	err := (*provider.Client).CreateDashboardObject(aclObject.Owner, aclObject.App, splunkDashboardsObj)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "ui", "views")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return splunkDashboardsRead(d, meta)
}

func splunkDashboardsRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadAllDashboardObject()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getDashboardByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadDashboardObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getDashboardByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("name", entry.Name); err != nil {
		return err
	}

	if err = d.Set("eai_data", entry.Content.EAIData); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func splunkDashboardsUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	splunkDashboardsObj := getSplunkDashboardsConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateDashboardObject(aclObject.Owner, aclObject.App, name, splunkDashboardsObj)
	if err != nil {
		return err
	}

	//ACL update
	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "ui", "views")
		if err != nil {
			return err
		}
	}

	return splunkDashboardsRead(d, meta)
}

func splunkDashboardsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// name := d.Get("name").(string)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	// splunkDashboardsObj := getSplunkDashboardsConfig(d)
	resp, err := (*provider.Client).DeleteDashboardObject(aclObject.Owner, aclObject.App, name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.DashboardResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getSplunkDashboardsConfig(d *schema.ResourceData) (splunkDashboardsObject *models.SplunkDashboardsObject) {
	splunkDashboardsObject = &models.SplunkDashboardsObject{}
	splunkDashboardsObject.Name = d.Get("name").(string)
	splunkDashboardsObject.EAIData = d.Get("eai_data").(string)
	return splunkDashboardsObject
}

func getDashboardByName(name string, httpResponse *http.Response) (dashboardEntry *models.DashboardEntry, err error) {
	response := &models.DashboardResponse{}
	//body, err := ioutil.ReadAll(httpResponse.Body)
	//fmt.Println(body)
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
		return dashboardEntry, err
	}

	return dashboardEntry, nil
}
