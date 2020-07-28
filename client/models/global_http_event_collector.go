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
	EnableSSL           bool   `json:"enableSSL,omitempty" url:"enableSSL"`
	Port                int    `json:"port,omitempty" url:"port,omitempty"`
	DedicatedIoThreads  int    `json:"dedicatedIoThreads,omitempty" url:"dedicatedIoThreads,omitempty"`
	MaxSockets          int    `json:"maxSockets,omitempty" url:"maxSockets,omitempty"`
	MaxThreads          int    `json:"maxThreads,omitempty" url:"maxThreads,omitempty"`
	UseDeploymentServer bool   `json:"useDeploymentServer,omitempty" url:"useDeploymentServer"`
}
