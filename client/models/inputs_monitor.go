package models

// Scripted Input Response Schema
type InputsMonitorResponse struct {
	Entry    []InputsMonitorEntry `json:"entry"`
	Messages []ErrorMessage       `json:"messages"`
}

type InputsMonitorEntry struct {
	Name    string              `json:"name"`
	ACL     ACLObject           `json:"acl"`
	Content InputsMonitorObject `json:"content"`
}

type InputsMonitorObject struct {
	Host            string `json:"host,omitempty" url:"host,omitempty"`
	Index           string `json:"index,omitempty" url:"index,omitempty"`
	SourceType      string `json:"sourcetype,omitempty" url:"sourcetype,omitempty"`
	RenameSource    string `json:"rename-source,omitempty" url:"rename-source,omitempty"`
	Disabled        bool   `json:"disabled,omitempty" url:"disabled"`
	CrcSalt         string `json:"crcSalt,omitempty" url:"crcSalt,omitempty"`
	FollowTail      bool   `json:"followTail,omitempty" url:"followTail"`
	Recursive       bool   `json:"recursive,omitempty" url:"recursive"`
	HostRegex       string `json:"host_regex,omitempty" url:"host_regex,omitempty"`
	HostSegment     int    `json:"host_segment,omitempty" url:"host_segment,omitempty"`
	TimeBeforeClose int    `json:"time-before-close,omitempty" url:"time-before-close,omitempty"`
	IgnoreOlderThan string `json:"ignore-older-than,omitempty" url:"ignore-older-than,omitempty"`
	Blacklist       string `json:"blacklist,omitempty" url:"blacklist,omitempty"`
	Whitelist       string `json:"whitelist,omitempty" url:"whitelist,omitempty"`
}
