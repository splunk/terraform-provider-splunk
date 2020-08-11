package models

// Scripted Input Response Schema
type InputsSplunkTCPTokenResponse struct {
	Entry    []InputsSplunkTCPTokenEntry `json:"entry"`
	Messages []ErrorMessage              `json:"messages"`
}

type InputsSplunkTCPTokenEntry struct {
	Name    string                     `json:"name"`
	ACL     ACLObject                  `json:"acl"`
	Content InputsSplunkTCPTokenObject `json:"content"`
}

type InputsSplunkTCPTokenObject struct {
	Name  string `json:"name,omitempty" url:"name,omitempty"`
	Token string `json:"token,omitempty" url:"token,omitempty"`
}
