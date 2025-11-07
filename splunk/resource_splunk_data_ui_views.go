package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/avast/retry-go/v4"

	"github.com/rsrdesarrollo/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func splunkDashboards() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Dashboard Name.",
			},
			"eai_data": {
				Type:        schema.TypeString,
				Required:    true,
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
	aclObject := getResourceDataViewACL(d)

	err := (*provider.Client).CreateDashboardObject(aclObject.Owner, aclObject.App, splunkDashboardsObj)
	if err != nil {
		return err
	}

	// add retry as sometimes dashboard object is not yet propagated and acl endpoint return 404
	err = retry.Do(
		func() error {
			err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "ui", "views")
			if err != nil {
				return err
			}
			return nil
		}, retry.Attempts(10), retry.OnRetry(func(n uint, err error) {
			log.Printf("#%d: %s. Retrying...\n", n, err)
		}), retry.DelayType(retry.BackOffDelay),
	)

	if err != nil {
		return err
	}

	d.SetId(name)
	return splunkDashboardsRead(d, meta)
}

func splunkDashboardsRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()

	aclObject := getResourceDataViewACL(d)

	readUser := "nobody"

	if aclObject.Sharing == "user" {
		// If we have a private dashboard we can only query it using the owner
		readUser = aclObject.Owner
	}

	resp, err := (*provider.Client).ReadDashboardObject(name, readUser, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getDashboardByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("unable to find resource: %v", name)
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
	aclObject := getResourceDataViewACL(d)

	updateUser := "nobody"

	if aclObject.Sharing == "user" {
		// If we have a private dashboard we can only update it using the owner
		updateUser = aclObject.Owner
	}

	if err := (*provider.Client).UpdateDashboardObject(updateUser, aclObject.App, name, splunkDashboardsObj); err != nil {
		return err
	}

	if err := (*provider.Client).UpdateAcl(updateUser, aclObject.App, name, aclObject, "data", "ui", "views"); err != nil {
		return err
	}

	return splunkDashboardsRead(d, meta)
}

func splunkDashboardsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	aclObject := getResourceDataViewACL(d)
	if aclObject.Sharing != "user" {
		aclObject.Owner = "nobody"
	}
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

// getResourceDataViewACL implements psuedo-defaults for the acl field for view resources.
func getResourceDataViewACL(d *schema.ResourceData) *models.ACLObject {
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.App = "search"
		aclObject.Owner = "admin"
		aclObject.Sharing = "user"
	}

	return aclObject
}
