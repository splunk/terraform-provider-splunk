package splunk

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

func adminProxyssoGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The Proxy SSO group name. Changing this forces a new resource.",
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Required. Set of internal Splunk roles assigned to the group. Order does not matter.",
			},
		},
		Read:   adminProxyssoGroupsRead,
		Create: adminProxyssoGroupsCreate,
		Delete: adminProxyssoGroupsDelete,
		Update: adminProxyssoGroupsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func adminProxyssoGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	adminProxyssoGroupsObj := getAdminProxyssoGroupsConfig(d)
	err := (*provider.Client).CreateAdminProxyssoGroups(name, adminProxyssoGroupsObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return adminProxyssoGroupsRead(d, meta)
}

func adminProxyssoGroupsRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()

	resp, err := (*provider.Client).ReadAdminProxyssoGroups(name)
	if err != nil {
		// Splunk returns 400 "Unable to find a role mapping" for a missing group.
		// The client surfaces that as an error after reading the body, so refresh
		// has to recognize it here and drop state.
		if adminProxyssoGroupMissing(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	defer resp.Body.Close()

	entry, err := getAdminProxyssoGroupsByName(name, resp)
	if err != nil {
		return err
	}
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

func adminProxyssoGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	adminProxyssoGroupsObj := getAdminProxyssoGroupsConfig(d)
	err := (*provider.Client).UpdateAdminProxyssoGroups(d.Id(), adminProxyssoGroupsObj)
	if err != nil {
		return err
	}

	return adminProxyssoGroupsRead(d, meta)
}

func adminProxyssoGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteAdminProxyssoGroups(d.Id())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return nil
	}
	return errors.New(resp.Status)
}

func getAdminProxyssoGroupsConfig(d *schema.ResourceData) *models.AdminProxyssoGroupsObject {
	obj := &models.AdminProxyssoGroupsObject{
		Name: d.Get("name").(string),
	}
	for _, v := range d.Get("roles").(*schema.Set).List() {
		obj.Roles = append(obj.Roles, v.(string))
	}
	return obj
}

func getAdminProxyssoGroupsByName(name string, httpResponse *http.Response) (*models.AdminProxyssoGroupsEntry, error) {
	response := &models.AdminProxyssoGroupsResponse{}
	decodeErr := json.NewDecoder(httpResponse.Body).Decode(response)

	switch httpResponse.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if decodeErr != nil {
			return nil, decodeErr
		}
		for i := range response.Entry {
			if response.Entry[i].Name == name {
				return &response.Entry[i], nil
			}
		}
		return nil, nil
	case http.StatusBadRequest:
		if adminProxyssoResponseMissing(response) {
			return nil, nil
		}
		return nil, errors.New(firstAdminProxyssoMessage(response))
	case http.StatusNotFound:
		return nil, nil
	default:
		if decodeErr != nil {
			return nil, decodeErr
		}
		return nil, errors.New(firstAdminProxyssoMessage(response))
	}
}

// adminProxyssoGroupMissing reports Splunk's "group does not exist" errors.
// client.DoRequest returns those as an error string that still contains the response body.
func adminProxyssoGroupMissing(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.Contains(msg, "Unable to find a role mapping") {
		return true
	}
	return strings.HasPrefix(msg, "404 ")
}

func adminProxyssoResponseMissing(response *models.AdminProxyssoGroupsResponse) bool {
	if len(response.Messages) == 0 {
		return false
	}
	return strings.Contains(response.Messages[0].Text, "Unable to find a role mapping")
}

func firstAdminProxyssoMessage(response *models.AdminProxyssoGroupsResponse) string {
	if len(response.Messages) == 0 {
		return "unexpected response from admin/ProxySSO-groups"
	}
	return response.Messages[0].Text
}
