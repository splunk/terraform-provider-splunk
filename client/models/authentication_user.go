package models

// Scripted Input Response Schema
type AuthenticationUserResponse struct {
	Entry    []AuthenticationUserEntry `json:"entry"`
	Messages []ErrorMessage            `json:"messages"`
}

type AuthenticationUserEntry struct {
	Name    string                   `json:"name"`
	ACL     ACLObject                `json:"acl"`
	Content AuthenticationUserObject `json:"content"`
}

type AuthenticationUserObject struct {
	DefaultApp            string   `json:"defaultApp,omitempty" url:"defaultApp,omitempty"`
	Email                 string   `json:"email,omitempty" url:"email,omitempty"`
	ForceChangePass       bool     `json:"force-change-pass,omitempty" url:"force-change-pass"`
	Name                  string   `json:"name,omitempty" url:"name,omitempty"`
	Password              string   `json:"password,omitempty" url:"password,omitempty"`
	RealName              string   `json:"realname,omitempty" url:"realname,omitempty"`
	RestartBackgroundJobs bool     `json:"restart_background_jobs,omitempty" url:"restart_background_jobs"`
	Roles                 []string `json:"roles,omitempty" url:"roles,omitempty"`
	TZ                    string   `json:"tz,omitempty" url:"tz,omitempty"`
}
