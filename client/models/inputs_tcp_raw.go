package models

// Scripted Input Response Schema
type InputsTCPRawResponse struct {
	Entry    []InputsTCPRawEntry `json:"entry"`
	Messages []ErrorMessage      `json:"messages"`
}

type InputsTCPRawEntry struct {
	Name    string             `json:"name"`
	ACL     ACLObject          `json:"acl"`
	Content InputsTCPRawObject `json:"content"`
}

type InputsTCPRawObject struct {
	Host           string `json:"host,omitempty" url:"host,omitempty"`
	Index          string `json:"index,omitempty" url:"index,omitempty"`
	Source         string `json:"source,omitempty" url:"source,omitempty"`
	SourceType     string `json:"sourcetype,omitempty" url:"sourcetype,omitempty"`
	Disabled       bool   `json:"disabled,omitempty" url:"disabled"`
	ConnectionHost string `json:"connection_host,omitempty" url:"connection_host,omitempty"`
	Queue          string `json:"queue,omitempty" url:"queue,omitempty"`
	RestrictToHost string `json:"restrictToHost,omitempty" url:"restrictToHost,omitempty"`
	RawTcpDoneTime int    `json:"rawTcpDoneTimeout,omitempty" url:"rawTcpDoneTimeout,omitempty"`
}
