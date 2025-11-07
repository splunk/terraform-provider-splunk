package splunk

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rsrdesarrollo/terraform-provider-splunk/client/models"
)

func genericAcl() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path to the object whose ACL is being managed.",
			},
			"acl": aclSchema(),
		},
		Create: genericAclCreate,
		Read:   genericAclRead,
		// Update does the same thing as Create, because the resource being managed has to already exist
		Update: genericAclCreate,
		Delete: genericAclDelete,
		Importer: &schema.ResourceImporter{
			State: genericAclImportState,
		},
	}
}

func genericAclCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	path := d.Get("path").(string)
	resources, name, ok := (*provider.Client).ResourcesAndNameForPath(path)
	if !ok {
		return fmt.Errorf("unable to parse path %s into resource and name parts", path)
	}

	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.App = "search"
		aclObject.Owner = "nobody"
	}

	err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, resources...)
	if err != nil {
		return err
	}

	d.SetId(path)

	return genericAclRead(d, meta)
}

func genericAclRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	path := d.Id()

	resources, name, ok := (*provider.Client).ResourcesAndNameForPath(path)
	if !ok {
		return fmt.Errorf("unable to parse path %s into resource and name parts", path)
	}

	r := d.Get("acl")
	aclObject := getACLConfig(r.([]interface{}))

	resp, err := (*provider.Client).GetAcl(aclObject.Owner, aclObject.App, name, resources...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	aclResponse := &models.ACLResponse{}
	if err := json.NewDecoder(resp.Body).Decode(aclResponse); err != nil {
		return err
	}

	if len(aclResponse.Entry) != 1 {
		return fmt.Errorf("ACLResponse has %d entries, expected exactly 1", len(aclResponse.Entry))
	}

	err = d.Set("acl", flattenACL(&aclResponse.Entry[0].Content))
	if err != nil {
		return err
	}

	return nil
}

func genericAclDelete(d *schema.ResourceData, meta interface{}) error {
	// Delete doesn't actually do anything, because an ACL can't be deleted.
	return nil
}

func genericAclImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	idParts := strings.Split(id, ":")
	if len(idParts) != 3 {
		return nil, fmt.Errorf("import ID must be owner:app:path")
	}

	owner := idParts[0]
	app := idParts[1]
	d.SetId(idParts[2])

	aclObject := &models.ACLObject{
		Owner: owner,
		App:   app,
	}

	if err := d.Set("acl", flattenACL(aclObject)); err != nil {
		return nil, err
	}

	if err := d.Set("path", d.Id()); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
