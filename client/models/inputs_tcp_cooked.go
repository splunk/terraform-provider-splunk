package models

// Scripted Input Response Schema
type InputsTCPCookedResponse struct {
	Entry    []InputsTCPCookedEntry `json:"entry"`
	Messages []ErrorMessage         `json:"messages"`
}

type InputsTCPCookedEntry struct {
	Name    string                `json:"name"`
	ACL     ACLObject             `json:"acl"`
	Content InputsTCPCookedObject `json:"content"`
}

type InputsTCPCookedObject struct {
	Host           string `json:"host,omitempty" url:"host,omitempty"`
	Index          string `json:"index,omitempty" url:"index,omitempty"`
	Disabled       bool   `json:"disabled,omitempty" url:"disabled"`
	ConnectionHost string `json:"connection_host,omitempty" url:"connection_host,omitempty"`
	RestrictToHost string `json:"restrictToHost,omitempty" url:"restrictToHost,omitempty"`
}
