package models

// Outputs TCP Default Response Schema
type OutputsTCPDefaultResponse struct {
	Entry    []OutputsTCPDefaultEntry `json:"entry"`
	Messages []ErrorMessage           `json:"messages"`
}

type OutputsTCPDefaultEntry struct {
	Name    string                  `json:"name"`
	ACL     ACLObject               `json:"acl"`
	Content OutputsTCPDefaultObject `json:"content"`
}

type OutputsTCPDefaultObject struct {
	DefaultGroup          string `json:"defaultGroup,omitempty" url:"defaultGroup,omitempty"`
	DropEventsOnQueueFull int    `json:"dropEventsOnQueueFull,omitempty" url:"dropEventsOnQueueFull,omitempty"`
	HeartbeatFrequency    int    `json:"heartbeatFrequency,omitempty" url:"heartbeatFrequency,omitempty"`
	MaxQueueSize          int    `json:"maxQueueSize,omitempty" url:"maxQueueSize,omitempty"`
	Disabled              bool   `json:"disabled,omitempty" url:"disabled"`
	IndexAndForward       bool   `json:"indexAndForward,omitempty" url:"indexAndForward"`
	SendCookedData        bool   `json:"sendCookedData,omitempty" url:"sendCookedData"`
}
