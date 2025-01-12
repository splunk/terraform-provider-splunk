package splunk

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/nealbrown/terraform-provider-splunk/client/models"

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
	_ = json.NewDecoder(httpResponse.Body).Decode(response)

	switch httpResponse.StatusCode {
	case 200, 201:
		re := regexp.MustCompile(`(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	case 400:
		// Splunk returns a 400 when a SAML group mapping is not found
		// try to catch that here
		re := regexp.MustCompile("Unable to find a role mapping for")
		if re.MatchString(response.Messages[0].Text) {
			return nil, nil
		}

		// but if the error didn't match, don't assume the 400 status was just a missing resource
		err := errors.New(response.Messages[0].Text)
		return nil, err
	}

	return nil, errors.New(response.Messages[0].Text)
}
