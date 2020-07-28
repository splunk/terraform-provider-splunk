package models

//https://docs.splunk.com/Documentation/Splunk/8.0.4/RESTUM/RESTusing#Access_Control_List
type ACLObject struct {
	App            string      `json:"app,omitempty" url:"app,omitempty"`
	CanChangePerms bool        `json:"can_change_perms,omitempty" url:"can_change_perms,omitempty"`
	CanList        bool        `json:"can_list,omitempty" url:"can_list,omitempty"`
	CanShareApp    bool        `json:"can_share_app,omitempty" url:"can_share_app,omitempty"`
	CanShareGlobal bool        `json:"can_share_global,omitempty" url:"can_share_global,omitempty"`
	CanShareUser   bool        `json:"can_share_user,omitempty" url:"can_share_user,omitempty"`
	CanWrite       bool        `json:"can_write,omitempty" url:"can_write,omitempty"`
	Modifiable     bool        `json:"modifiable,omitempty" url:"modifiable,omitempty"`
	Owner          string      `json:"owner,omitempty" url:"owner,omitempty"`
	Removable      bool        `json:"removable,omitempty" url:"removable,omitempty"`
	Sharing        string      `json:"sharing,omitempty" url:"sharing,omitempty"`
	Perms          Permissions `json:"perms,omitempty" url:"perms,omitempty"`
}

type Permissions struct {
	Read  []string `json:"read,omitempty" url:"read,omitempty"`
	Write []string `json:"write,omitempty" url:"write,omitempty"`
}

type ACLResponse struct {
	Entry    []ACLEntry     `json:"entry"`
	Messages []ErrorMessage `json:"messages"`
}

type ACLEntry struct {
	Name    string    `json:"name"`
	Content ACLObject `json:"acl"`
}
