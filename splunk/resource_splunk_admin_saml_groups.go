package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jaware-splunk/terraform-provider-splunk/client/models"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func adminSAMLGroups() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The external group name.",
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
		Read:   adminSAMLGroupsRead,
		Create: adminSAMLGroupsCreate,
		Delete: adminSAMLGroupsDelete,
		Update: adminSAMLGroupsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func adminSAMLGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	adminSAMLGroupsObj := getAdminSAMLGroupsConfig(d)
	err := (*provider.Client).CreateAdminSAMLGroups(name, adminSAMLGroupsObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return adminSAMLGroupsRead(d, meta)
}

func adminSAMLGroupsRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()

	// Read the SAML group
	resp, err := (*provider.Client).ReadAdminSAMLGroups(name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getAdminSAMLGroupsByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("name", entry.Name); err != nil {
		return err
	}

	if err = d.Set("roles", entry.Content.Roles); err != nil {
		return err
	}

	return nil
}

func adminSAMLGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	adminSAMLGroupsObj := getAdminSAMLGroupsConfig(d)
	err := (*provider.Client).UpdateAdminSAMLGroups(d.Id(), adminSAMLGroupsObj)
	if err != nil {
		return err
	}

	return adminSAMLGroupsRead(d, meta)
}

func adminSAMLGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteAdminSAMLGroups(d.Id())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.AdminSAMLGroupsResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getAdminSAMLGroupsConfig(d *schema.ResourceData) (adminSAMLGroupsObject *models.AdminSAMLGroupsObject) {
	adminSAMLGroupsObject = &models.AdminSAMLGroupsObject{}
	adminSAMLGroupsObject.Name = d.Get("name").(string)
	if val, ok := d.GetOk("roles"); ok {
		for _, v := range val.([]interface{}) {
			adminSAMLGroupsObject.Roles = append(adminSAMLGroupsObject.Roles, v.(string))
		}
	}
	return adminSAMLGroupsObject
}

func getAdminSAMLGroupsByName(name string, httpResponse *http.Response) (AdminSAMLGroupsEntry *models.AdminSAMLGroupsEntry, err error) {
	response := &models.AdminSAMLGroupsResponse{}
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
		return AdminSAMLGroupsEntry, err
	}

	return AdminSAMLGroupsEntry, nil
}
