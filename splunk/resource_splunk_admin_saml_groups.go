package splunk

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/go-splunk-client/pkg/entry"
	"github.com/splunk/terraform-provider-splunk/client/models"
	"github.com/splunk/terraform-provider-splunk/internal/sync"

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
			"use_client": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateClientValueFunc(true),
				Description: "Set to explicitly specify which client to use for this resource. Leave unset to use the provider's default. Permitted non-empty values are legacy and external. " +
					"The legacy client is being replaced by a standalone Splunk client with improved error and drift handling. The legacy client will be deprecated in a future version.",
			},
		},
		Read:   readFunc(samlGroupSync, adminSAMLGroupsRead),
		Create: createFunc(samlGroupSync, adminSAMLGroupsCreate),
		Delete: deleteFunc(samlGroupSync, adminSAMLGroupsDelete),
		Update: updateFunc(samlGroupSync, adminSAMLGroupsUpdate),
		Importer: &schema.ResourceImporter{
			State: importStateSAMLGroups,
		},
	}
}

// samlGroupSync returns a SyncGetter that manages entry.SAMLGroup.
func samlGroupSync() sync.SyncGetter {
	var group entry.SAMLGroup

	samlGroupSync := sync.ComposeSync(
		sync.NewClientID(&group.ID, "name"),
		sync.NewDirectListField("roles", &group.Content.Roles),
	)

	return sync.NewIndirectObject(&group, samlGroupSync)
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

	// name defaults to the stored Id, unless it's a parseable client.ID
	name := d.Id()
	if id, err := client.ParseID(name); err == nil {
		// if d.Id is a parseable client.ID, this was most recently using the non-legacy client
		// restore the internal workings back to the legacy client format
		name = id.Title
		d.SetId(name)
	}

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
	err := (*provider.Client).UpdateAdminSAMLGroups(d.Get("name").(string), adminSAMLGroupsObj)
	if err != nil {
		return err
	}

	return adminSAMLGroupsRead(d, meta)
}

func adminSAMLGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteAdminSAMLGroups(d.Get("name").(string))
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

// importStateSAMLGroups calls schema.ImportStatePassthrough after setting use_client based on the ID's
// ability to be parsed into a client.ID.
func importStateSAMLGroups(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if _, err := client.ParseID(d.Id()); err == nil {
		if err := d.Set("use_client", useClientExternal); err != nil {
			return nil, err
		}
	} else {
		if err := d.Set("use_client", useClientLegacy); err != nil {
			return nil, err
		}
	}

	return schema.ImportStatePassthrough(d, m)
}
