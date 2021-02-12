package models

// Outputs TCP Syslog Response Schema
type OutputsTCPSyslogResponse struct {
	Entry    []OutputsTCPSyslogEntry `json:"entry"`
	Messages []ErrorMessage          `json:"messages"`
}

type OutputsTCPSyslogEntry struct {
	Name    string                 `json:"name"`
	ACL     ACLObject              `json:"acl"`
	Content OutputsTCPSyslogObject `json:"content"`
}

type OutputsTCPSyslogObject struct {
	Disabled         bool   `json:"disabled,omitempty" url:"disabled"`
	Priority         int    `json:"priority,string,omitempty" url:"priority,omitempty"`
	Server           string `json:"server,omitempty" url:"server,omitempty"`
	SyslogSourceType string `json:"syslogSourceType,omitempty" url:"syslogSourceType,omitempty"`
	TimestampFormat  string `json:"timestampformat,omitempty" url:"timestampformat,omitempty"`
	Type             string `json:"type,omitempty" url:"type,omitempty"`
}
