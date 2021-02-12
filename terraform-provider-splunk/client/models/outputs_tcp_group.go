package models

// Outputs TCP Group Response Schema
type OutputsTCPGroupResponse struct {
	Entry    []OutputsTCPGroupEntry `json:"entry"`
	Messages []ErrorMessage         `json:"messages"`
}

type OutputsTCPGroupEntry struct {
	Name    string                `json:"name"`
	ACL     ACLObject             `json:"acl"`
	Content OutputsTCPGroupObject `json:"content"`
}

type OutputsTCPGroupObject struct {
	Compressed            bool     `json:"compressed,omitempty" url:"compressed"`
	Disabled              bool     `json:"disabled,omitempty" url:"disabled"`
	DropEventsOnQueueFull int      `json:"dropEventsOnQueueFull,omitempty" url:"dropEventsOnQueueFull,omitempty"`
	HeartbeatFrequency    int      `json:"heartbeatFrequency,omitempty" url:"heartbeatFrequency,omitempty"`
	MaxQueueSize          string   `json:"maxQueueSize,omitempty" url:"maxQueueSize,omitempty"`
	Method                string   `json:"method,omitempty" url:"method,omitempty"`
	SendCookedData        bool     `json:"sendCookedData,omitempty" url:"sendCookedData"`
	Servers               []string `json:"servers,omitempty" url:"servers,omitempty"`
	Token                 string   `json:"token,omitempty" url:"token,omitempty"`
}
