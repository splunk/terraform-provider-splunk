package models

// Scripted Input Response Schema
type AdminSAMLGroupsResponse struct {
	Entry    []AdminSAMLGroupsEntry `json:"entry"`
	Messages []ErrorMessage         `json:"messages"`
}

type AdminSAMLGroupsEntry struct {
	Name    string                `json:"name"`
	ACL     ACLObject             `json:"acl"`
	Content AdminSAMLGroupsObject `json:"content"`
}

type AdminSAMLGroupsObject struct {
	Name  string   `json:"name,omitempty" url:"name,omitempty"`
	Roles []string `json:"roles,omitempty" url:"roles,omitempty"`
}
