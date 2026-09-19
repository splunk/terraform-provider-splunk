package models

// FederatedProviderResponse is the JSON response returned by the federated
// provider REST endpoint.
type FederatedProviderResponse struct {
	Entry    []FederatedProviderEntry `json:"entry"`
	Messages []ErrorMessage           `json:"messages"`
}

// FederatedProviderEntry represents one federated provider definition.
type FederatedProviderEntry struct {
	Name    string                  `json:"name"`
	Content FederatedProviderObject `json:"content"`
}

// FederatedProviderObject contains the Enterprise Splunk-to-Splunk provider
// settings supported by this provider.
type FederatedProviderObject struct {
	Type           string `json:"type,omitempty" url:"type,omitempty"`
	Mode           string `json:"mode,omitempty" url:"mode,omitempty"`
	AppContext     string `json:"appContext,omitempty" url:"appContext,omitempty"`
	HostPort       string `json:"hostPort,omitempty" url:"hostPort,omitempty"`
	ServiceAccount string `json:"serviceAccount,omitempty" url:"serviceAccount,omitempty"`
	Password       string `json:"password,omitempty" url:"password,omitempty"`
}
