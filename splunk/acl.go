package splunk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-splunk/client/models"
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
					Description: "The app context for the resource. Required for updating saved search ACL properties. Allowed values are:" +
						"The name of an app and system",
				},
				"owner": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					Description: "User name of resource owner. Defaults to the resource creator. Required for updating any knowledge object ACL properties." +
						"nobody = All users may access the resource, but write access to the resource might be restricted.",
				},
				"sharing": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					Description: "Indicates how the resource is shared. Required for updating any knowledge object ACL properties." +
						"app: Shared within a specific app" +
						"global: (Default) Shared globally to all apps." +
						"user: Private to a user",
				},
				"read": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Properties that indicate resource read permissions.",
				},
				"write": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Properties that indicate write permissions of the resource.",
				},
				"can_change_perms": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Indicates if the active user can change permissions for this object. Defaults to true.",
				},
				"can_share_app": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Indicates if the active user can change sharing to app level. Defaults to true.",
				},
				"can_share_global": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Indicates if the active user can change sharing to system level. Defaults to true.",
				},
				"can_share_user": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Indicates if the active user can change sharing to user level. Defaults to true.",
				},
				"can_write": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Indicates if the active user can edit this object. Defaults to true.",
				},
				"removable": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "Indicates whether an admin or user with sufficient permissions can delete the entity.",
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
		"app":              acl.App,
		"owner":            acl.Owner,
		"sharing":          acl.Sharing,
		"read":             acl.Perms.Read,
		"write":            acl.Perms.Write,
		"removable":        acl.Removable,
		"can_write":        acl.CanWrite,
		"can_share_app":    acl.CanShareApp,
		"can_share_user":   acl.CanShareUser,
		"can_share_global": acl.CanShareGlobal,
		"can_change_perms": acl.CanChangePerms,
	}

	return []interface{}{m}
}
