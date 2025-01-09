package models

// ServerClassResponse represents the response structure for server classes
type ServerClassResponse struct {
	Entry    []ServerClassEntry `json:"entry"`
	Messages []ErrorMessage     `json:"messages"`
}

// ServerClassEntry represents an individual server class entry
type ServerClassEntry struct {
	Name    string            `json:"name"`
	ACL     ACLObject         `json:"acl"`
	Content ServerClassObject `json:"content"`
}

// ServerClassObject represents the content of a server class
type ServerClassObject struct {
	RestartSplunkWeb bool     `json:"restartSplunkWeb,omitempty" url:"restartSplunkWeb,omitempty"`
	RestartSplunkd   bool     `json:"restartSplunkd,omitempty" url:"restartSplunkd,omitempty"`
	Whitelist        []string `json:"whitelist,omitempty" url:"whitelist,omitempty"`
	Blacklist        []string `json:"blacklist,omitempty" url:"blacklist,omitempty"`
	Apps             []string `json:"apps,omitempty" url:"apps,omitempty"`
	Name             string   `json:"name,omitempty" url:"name,omitempty"`
	Description      string   `json:"description,omitempty" url:"description,omitempty"`
}
