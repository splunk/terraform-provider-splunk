package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-splunk/client/models"
)

func aclSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"app": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"owner": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"sharing": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"read": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"write": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func getACLConfig(r []interface{}) (acl *models.ACLObject) {
	acl = &models.ACLObject{}
	for _, v := range r {
		a := v.(map[string]interface{})

		if a["app"] != "" {
			acl.App = a["app"].(string)
		} else {
			acl.App = "search"
		}

		if a["owner"] != "" {
			acl.Owner = a["owner"].(string)
		} else {
			acl.Owner = "nobody"
		}

		if a["sharing"] != "" {
			acl.Sharing = a["sharing"].(string)
		} else {
			acl.Sharing = "app"
		}

		for _, v := range a["read"].([]interface{}) {
			acl.Perms.Read = append(acl.Perms.Read, v.(string))
		}

		for _, w := range a["write"].([]interface{}) {
			acl.Perms.Write = append(acl.Perms.Write, w.(string))
		}
	}

	return acl
}

func flattenACL(acl *models.ACLObject) []interface{} {
	if acl == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"app":     acl.App,
		"owner":   acl.Owner,
		"sharing": acl.Sharing,
		"read":    acl.Perms.Read,
		"write":   acl.Perms.Write,
	}

	return []interface{}{m}
}
