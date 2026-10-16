package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
			State: splunkDashboardsImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultReplicationTimeout),
		},
	}
}

func splunkDashboardsImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	imported, err := importNamespacedResourceState(d, "data", "ui", "views")
	if err != nil {
		return nil, err
	}
	if !imported {
		if err := d.Set("name", d.Id()); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

// Functions
func splunkDashboardsCreate(d *schema.ResourceData, meta interface{}) error {
	remaining := startCreateDeadline(d)
	defer stopCreateDeadline(d)

	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	splunkDashboardsObj := getSplunkDashboardsConfig(d)
	aclObject := getResourceDataViewACL(d)

	// retry 409 as a dashboard deleted by replace may still exist on other cluster members
	resp, err := retryOnStatus(remaining(), http.StatusConflict, func() (*http.Response, error) {
		return (*provider.Client).CreateDashboardObject(aclObject.Owner, aclObject.App, splunkDashboardsObj)
	})
	closeBody(resp)
	if err != nil {
		return err
	}

	// A new dashboard is private to its owner until this call shares it, so nobody
	// returns 404. Later reads, updates, and deletes of a shared dashboard use nobody.
	err = retryUpdateAcl(remaining(), func() error {
		return (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "ui", "views")
	})
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

	resp, err := readKnowledgeObject(d, func() (*http.Response, error) {
		return (*provider.Client).ReadDashboardObject(name, sharedObjectOwner(aclObject), aclObject.App)
	})
	if err != nil {
		if hasStatus(resp, http.StatusNotFound) && !d.IsNewResource() {
			log.Printf("[WARN] Dashboard %s not found, removing from state", name)
			closeBody(resp)
			d.SetId("")
			return nil
		}
		closeBody(resp)
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

	// Updating a shared dashboard in the owner's namespace creates a private copy.
	updateUser := sharedObjectOwner(aclObject)

	if err := (*provider.Client).UpdateDashboardObject(updateUser, aclObject.App, name, splunkDashboardsObj); err != nil {
		return err
	}

	// changing acl needs owner or admin rights, so only send it when it changed
	if d.HasChange("acl") {
		if err := (*provider.Client).UpdateAcl(updateUser, aclObject.App, name, aclObject, "data", "ui", "views"); err != nil {
			return err
		}
	}

	return splunkDashboardsRead(d, meta)
}

func splunkDashboardsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	aclObject := getResourceDataViewACL(d)
	resp, err := retryOnStatus(replicationRefreshTimeout, http.StatusNotFound, func() (*http.Response, error) {
		return (*provider.Client).DeleteDashboardObject(sharedObjectOwner(aclObject), aclObject.App, name)
	})
	if hasStatus(resp, http.StatusNotFound) {
		closeBody(resp)
		return nil
	}
	if err != nil {
		closeBody(resp)
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
