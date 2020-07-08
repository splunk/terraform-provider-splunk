package models

// Outputs TCP Server Response Schema
type OutputsTCPServerResponse struct {
	Entry    []OutputsTCPServerEntry `json:"entry"`
	Messages []ErrorMessage          `json:"messages"`
}

type OutputsTCPServerEntry struct {
	Name    string                 `json:"name"`
	ACL     ACLObject              `json:"acl"`
	Content OutputsTCPServerObject `json:"content"`
}

type OutputsTCPServerObject struct {
	Disabled             bool   `json:"disabled,omitempty" url:"disabled"`
	Method               string `json:"method,omitempty" url:"method,omitempty"`
	SSLAltNameToCheck    string `json:"sslAltNameToCheck,omitempty" url:"sslAltNameToCheck,omitempty"`
	SSLCertPath          string `json:"sslCertPath,omitempty" url:"sslCertPath,omitempty"`
	SSLCipher            string `json:"sslCipher,omitempty" url:"sslCipher,omitempty"`
	SSLCommonNameToCheck string `json:"sslCommonNameToCheck,omitempty" url:"sslCommonNameToCheck,omitempty"`
	SSLPassword          string `json:"sslPassword,omitempty" url:"sslPassword,omitempty"`
	SSLRootCAPath        string `json:"sslRootCAPath,omitempty" url:"sslRootCAPath,omitempty"`
	SSLVerifyServerCert  bool   `json:"sslVerifyServerCert,omitempty" url:"sslVerifyServerCert"`
}
