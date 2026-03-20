package splunk

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/splunk/terraform-provider-splunk/client/models"
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
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice([]string{"user", "app", "global"}, false),
					ForceNew:     true,
					Description: "Indicates how the resource is shared. Required for updating any knowledge object ACL properties." +
						"app: Shared within a specific app" +
						"global: (Default) Shared globally to all apps." +
						"user: Private to a user",
				},
				"read": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Properties that indicate resource read permissions.",
				},
				"write": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Properties that indicate resource write permissions.",
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

func aclValidator(diff *schema.ResourceDiff, v interface{}) error {
	acl := diff.Get("acl").([]interface{})
	if len(acl) == 0 {
		return nil
	}
	// Assert that each item is a map[string]interface{}
	aclMap, ok := acl[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Value cannot be mapped to map!")
	}

	// Check if sharing has changed to "user"
	if diff.HasChange("acl.0.sharing") {
		_, new := diff.GetChange("acl.0.sharing")
		if new == "user" {
			// Check if `read` or `write` attributes are explicitly set in the configuration, ignoring persisted state
			if diff.HasChange("acl.0.read") {
				if aclMap["read"] != nil && len(aclMap["read"].([]interface{})) > 0 {
					return fmt.Errorf("`acl.read` cannot be set when `acl.sharing` is `user`")
				}
			}
			if diff.HasChange("acl.0.write") {
				if aclMap["write"] != nil && len(aclMap["write"].([]interface{})) > 0 {
					return fmt.Errorf("`acl.write` cannot be set when `acl.sharing` is `user`")
				}
			}
		}
	}
	return nil
}

// defaultACLConfigForGenericResource matches splunk_generic_acl create when no acl block is set.
func defaultACLConfigForGenericResource() *models.ACLObject {
	return &models.ACLObject{
		App:     "search",
		Owner:   "nobody",
		Sharing: "app",
	}
}

func aclStringFromMapOptional(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func aclStringListFromInterface(raw interface{}) []string {
	if raw == nil {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, v := range list {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func getACLConfig(r []interface{}) (acl *models.ACLObject) {
	acl = &models.ACLObject{}
	for _, v := range r {
		a, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if s := aclStringFromMapOptional(a, "app"); s != "" {
			acl.App = s
		} else {
			acl.App = "search"
		}
		if s := aclStringFromMapOptional(a, "owner"); s != "" {
			acl.Owner = s
		} else {
			acl.Owner = "nobody"
		}
		if s := aclStringFromMapOptional(a, "sharing"); s != "" {
			acl.Sharing = s
		} else {
			acl.Sharing = "app"
		}
		for _, x := range aclStringListFromInterface(a["read"]) {
			acl.Perms.Read = append(acl.Perms.Read, x)
		}
		for _, x := range aclStringListFromInterface(a["write"]) {
			acl.Perms.Write = append(acl.Perms.Write, x)
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
