package models

type SplunkLookupDefinitionResponse struct {
	Entry    []SplunkLookupDefinitionEntry `json:"entry"`
	Messages []ErrorMessage                `json:"messages"`
}

type SplunkLookupDefinitionEntry struct {
	Name    string                       `json:"name"`
	ACL     ACLObject                    `json:"acl"`
	Content SplunkLookupDefinitionObject `json:"content"`
}

type SplunkLookupDefinitionObject struct {
	Name     string `json:"name,omitempty" url:"name,omitempty"`
	Filename string `json:"filename,omitempty" url:"filename,omitempty"`
}
