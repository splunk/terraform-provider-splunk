package models

// AdminProxyssoGroupsResponse is the JSON body from admin/ProxySSO-groups.
type AdminProxyssoGroupsResponse struct {
	Entry    []AdminProxyssoGroupsEntry `json:"entry"`
	Messages []ErrorMessage             `json:"messages"`
}

// AdminProxyssoGroupsEntry is one Proxy SSO group mapping.
type AdminProxyssoGroupsEntry struct {
	Name    string                    `json:"name"`
	ACL     ACLObject                 `json:"acl"`
	Content AdminProxyssoGroupsObject `json:"content"`
}

// AdminProxyssoGroupsObject is the group name and the Splunk roles it maps to.
type AdminProxyssoGroupsObject struct {
	Name  string   `json:"name,omitempty" url:"name,omitempty"`
	Roles []string `json:"roles,omitempty" url:"roles,omitempty"`
}
