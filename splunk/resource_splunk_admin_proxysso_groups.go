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

func adminProxyssoGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The Proxy SSO group name.",
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Required. List of internal roles assigned to group.",
			},
		},
		Read:   AdminProxyssoGroupsRead,
		Create: AdminProxyssoGroupsCreate,
		Delete: AdminProxyssoGroupsDelete,
		Update: AdminProxyssoGroupsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func AdminProxyssoGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	authenticationUserObj := getAdminProxyssoGroupsConfig(d)
	err := (*provider.Client).CreateAdminProxyssoGroups(name, authenticationUserObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return AdminProxyssoGroupsRead(d, meta)
}

func AdminProxyssoGroupsRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadAllAdminProxyssoGroups()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getAdminProxyssoGroupsByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadAdminProxyssoGroups(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getAdminProxyssoGroupsByName(name, resp)
	if err != nil {
		return err
	}

	// an empty entry (with no error) means the resource wasn't found
	// mark it as such so it can be re-created
	if entry == nil {
		d.SetId("")
		return nil
	}

	if err = d.Set("name", entry.Name); err != nil {
		return err
	}

	if err = d.Set("roles", entry.Content.Roles); err != nil {
		return err
	}

	return nil
}

func AdminProxyssoGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	authenticationUserObj := getAdminProxyssoGroupsConfig(d)
	err := (*provider.Client).UpdateAdminProxyssoGroups(d.Id(), authenticationUserObj)
	if err != nil {
		return err
	}

	return AdminProxyssoGroupsRead(d, meta)
}

func AdminProxyssoGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteAdminProxyssoGroups(d.Id())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.AdminProxyssoGroupsResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getAdminProxyssoGroupsConfig(d *schema.ResourceData) (authenticationUserObject *models.AdminProxyssoGroupsObject) {
	authenticationUserObject = &models.AdminProxyssoGroupsObject{}
	authenticationUserObject.Name = d.Get("name").(string)
	if val, ok := d.GetOk("roles"); ok {
		for _, v := range val.([]interface{}) {
			authenticationUserObject.Roles = append(authenticationUserObject.Roles, v.(string))
		}
	}
	return authenticationUserObject
}

func getAdminProxyssoGroupsByName(name string, httpResponse *http.Response) (AdminProxyssoGroupsEntry *models.AdminProxyssoGroupsEntry, err error) {
	response := &models.AdminProxyssoGroupsResponse{}
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
		return AdminProxyssoGroupsEntry, err
	}

	return AdminProxyssoGroupsEntry, nil
}
