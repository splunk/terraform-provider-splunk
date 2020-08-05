package models

// HTTP Input Response Schema
type GlobalHECResponse struct {
	Entry    []GlobalHECEntry `json:"entry"`
	Messages []ErrorMessage   `json:"messages"`
}

type GlobalHECEntry struct {
	Name    string                         `json:"name"`
	Content GlobalHttpEventCollectorObject `json:"content"`
}

type GlobalHttpEventCollectorObject struct {
	Host                string `json:"host,omitempty" url:"host,omitempty"`
	Index               string `json:"index,omitempty" url:"index,omitempty"`
	Disabled            bool   `json:"disabled,omitempty" url:"disabled"`
	EnableSSL           string `json:"enableSSL,omitempty" url:"enableSSL,omitempty"`
	Port                string `json:"port,omitempty" url:"port,omitempty"`
	DedicatedIoThreads  string `json:"dedicatedIoThreads,omitempty" url:"dedicatedIoThreads,omitempty"`
	MaxSockets          string `json:"maxSockets,omitempty" url:"maxSockets,omitempty"`
	MaxThreads          string `json:"maxThreads,omitempty" url:"maxThreads,omitempty"`
	UseDeploymentServer string `json:"useDeploymentServer,omitempty" url:"useDeploymentServer,omitempty"`
}
