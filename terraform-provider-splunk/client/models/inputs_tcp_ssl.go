package models

// Scripted Input Response Schema
type InputsTCPSSLResponse struct {
	Entry    []InputsTCPSSLEntry `json:"entry"`
	Messages []ErrorMessage      `json:"messages"`
}

type InputsTCPSSLEntry struct {
	Name    string             `json:"name"`
	ACL     ACLObject          `json:"acl"`
	Content InputsTCPSSLObject `json:"content"`
}

type InputsTCPSSLObject struct {
	Disabled          bool   `json:"disabled,omitempty" url:"disabled"`
	RequireClientCert bool   `json:"requireClientCert,omitempty" url:"requireClientCert"`
	Password          string `json:"password,omitempty" url:"password,omitempty"`
	RootCA            string `json:"rootCA,omitempty" url:"rootCA,omitempty"`
	ServerCert        string `json:"serverCert,omitempty" url:"serverCert,omitempty"`
}
