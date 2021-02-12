package models

// DashboardResponse HTTP Input Response Schema
type DashboardResponse struct {
	Entry    []DashboardEntry `json:"entry"`
	Messages []ErrorMessage   `json:"messages"`
}

type DashboardEntry struct {
	Name    string                 `json:"name"`
	ACL     ACLObject              `json:"acl"`
	Content SplunkDashboardsObject `json:"content"`
}

type SplunkDashboardsObject struct {
	Name    string `json:"name,omitempty" url:"name,omitempty"`
	EAIData string `json:"eai:data,omitempty" url:"eai:data,omitempty"`
}
