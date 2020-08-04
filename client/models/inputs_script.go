package models

// Scripted Input Response Schema
type ScriptedInputResponse struct {
	Entry    []ScriptedInputEntry `json:"entry"`
	Messages []ErrorMessage       `json:"messages"`
}

type ScriptedInputEntry struct {
	Name    string             `json:"name"`
	ACL     ACLObject          `json:"acl"`
	Content InputsScriptObject `json:"content"`
}

type InputsScriptObject struct {
	Host         string `json:"host,omitempty" url:"host,omitempty"`
	Index        string `json:"index,omitempty" url:"index,omitempty"`
	Source       string `json:"source,omitempty" url:"source,omitempty"`
	SourceType   string `json:"sourcetype,omitempty" url:"sourcetype,omitempty"`
	RenameSource string `json:"rename-source,omitempty" url:"rename-source,omitempty"`
	PassAuth     string `json:"passAuth,omitempty" url:"passAuth,omitempty"`
	Disabled     bool   `json:"disabled,omitempty" url:"disabled,omitempty"`
	Interval     int    `json:"interval,omitempty" url:"interval,omitempty"`
}
