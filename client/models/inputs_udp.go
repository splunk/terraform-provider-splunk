package models

// Scripted Input Response Schema
type InputsUDPResponse struct {
	Entry    []InputsUDPEntry `json:"entry"`
	Messages []ErrorMessage   `json:"messages"`
}

type InputsUDPEntry struct {
	Name    string          `json:"name"`
	ACL     ACLObject       `json:"acl"`
	Content InputsUDPObject `json:"content"`
}

type InputsUDPObject struct {
	Host                 string `json:"host,omitempty" url:"host,omitempty"`
	Index                string `json:"index,omitempty" url:"index,omitempty"`
	Source               string `json:"source,omitempty" url:"source,omitempty"`
	SourceType           string `json:"sourcetype,omitempty" url:"sourcetype,omitempty"`
	Disabled             bool   `json:"disabled,omitempty" url:"disabled"`
	ConnectionHost       string `json:"connection_host,omitempty" url:"connection_host,omitempty"`
	Queue                string `json:"queue,omitempty" url:"queue,omitempty"`
	RestrictToHost       string `json:"restrictToHost,omitempty" url:"restrictToHost,omitempty"`
	NoAppendingTimestamp bool   `json:"no_appending_timestamp,omitempty" url:"no_appending_timestamp"`
	NoPriorityStripping  bool   `json:"no_priority_stripping,omitempty" url:"no_priority_stripping"`
}
