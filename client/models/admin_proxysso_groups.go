package models

// Scripted Input Response Schema
type AdminProxyssoGroupsResponse struct {
	Entry    []AdminProxyssoGroupsEntry `json:"entry"`
	Messages []ErrorMessage            `json:"messages"`
}

type AdminProxyssoGroupsEntry struct {
	Name    string                   `json:"name"`
	ACL     ACLObject                `json:"acl"`
	Content AdminProxyssoGroupsObject `json:"content"`
}

type AdminProxyssoGroupsObject struct {
	Roles []string `json:"roles,omitempty" url:"roles,omitempty"`
	Name  string   `json:"name,omitempty" url:"name,omitempty"`
}
