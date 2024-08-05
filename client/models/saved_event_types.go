package models

type SavedEventTypesResponse struct {
	Entry    []SavedEventTypesEntry `json:"entry"`
	Messages []ErrorMessage         `json:"messages"`
}

type SavedEventTypesEntry struct {
	Name    string               `json:"name"`
	ACL     ACLObject            `json:"acl"`
	Content SavedEventTypeObject `json:"content"`
}

type SavedEventTypeObject struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Disabled    bool     `json:"disabled,omitempty" url:"disabled,omitempty"`
	Color       string   `json:"color,omitempty" url:"color,omitempty"`
	Priority    int      `json:"priority,omitempty" url:"priority,omitempty"`
	Search      string   `json:"search,omitempty" url:"search,omitempty"`
	Tags        []string `json:"tags,omitempty" url:"tags,omitempty"`
}
