package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/go-splunk-client/pkg/entry"
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
	c := provider.ExternalClient

	samlGroup, err := getAdminSAMLGroupsConfig(d)
	if err != nil {
		return err
	}

	if err := c.Create(samlGroup); err != nil {
		return err
	}

	return adminSAMLGroupsRead(d, meta)
}

func adminSAMLGroupsRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	c := provider.ExternalClient

	samlGroup, err := getAdminSAMLGroupsConfig(d)
	if err != nil {
		return err
	}

	if err := c.Read(&samlGroup); err != nil {
		if clientError, ok := err.(client.Error); ok {
			if clientError.Code == client.ErrorNotFound {
				d.SetId("")
				return nil
			}
		}

		return err
	}

	id, err := samlGroup.ID.URL()
	if err != nil {
		return err
	}
	d.SetId(id)

	if err := d.Set("name", samlGroup.ID.Title); err != nil {
		return err
	}

	if err := d.Set("roles", samlGroup.Content.Roles); err != nil {
		return err
	}

	return nil
}

func adminSAMLGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	c := provider.ExternalClient

	samlGroup, err := getAdminSAMLGroupsConfig(d)
	if err != nil {
		return err
	}

	if err := c.Update(samlGroup); err != nil {
		return err
	}

	return adminSAMLGroupsRead(d, meta)
}

func adminSAMLGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	c := provider.ExternalClient

	samlGroup, err := getAdminSAMLGroupsConfig(d)
	if err != nil {
		return err
	}

	if err := c.Delete(samlGroup); err != nil {
		return err
	}

	return nil
}

func getAdminSAMLGroupsConfig(d *schema.ResourceData) (entry.SAMLGroup, error) {
	var samlGroup entry.SAMLGroup

	if d.Id() == "" {
		samlGroup.ID = client.ID{
			Title: d.Get("name").(string),
		}
	} else {
		id, err := client.ParseID(d.Id())
		if err != nil {
			return entry.SAMLGroup{}, err
		}

		samlGroup.ID = id
	}

	// undefined roles indicates an empty list of roles
	roles := []string{}
	for _, roleI := range d.Get("roles").([]interface{}) {
		role := roleI.(string)
		roles = append(roles, role)
	}
	samlGroup.Content.Roles = roles

	return samlGroup, nil
}
