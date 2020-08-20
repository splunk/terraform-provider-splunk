package models

// Local App Schema
type AppsLocalResponse struct {
	Entry    []AppsLocalEntry `json:"entry"`
	Messages []ErrorMessage   `json:"messages"`
}
type AppsLocalEntry struct {
	Name    string          `json:"name"`
	ACL     ACLObject       `json:"acl"`
	Content AppsLocalObject `json:"content"`
}
type AppsLocalObject struct {
	Author          string `json:"author,omitempty" url:"author,omitempty"`
	Configured      bool   `json:"configured,omitempty" url:"configured"`
	Description     string `json:"description,omitempty" url:"description,omitempty"`
	ExplicitAppName string `json:"explicit_appname,omitempty" url:"explicit_appname,omitempty"`
	Filename        bool   `json:"filename,omitempty" url:"filename"`
	Label           string `json:"label,omitempty" url:"label,omitempty"`
	Name            string `json:"name,omitempty" url:"name,omitempty"`
	Version         string `json:"version,omitempty" url:"version,omitempty"`
	Visible         bool   `json:"visible,omitempty" url:"visible"`
	Update          bool   `json:"update,omitempty" url:"update"`
}
