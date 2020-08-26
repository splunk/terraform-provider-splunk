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
	Disabled            bool   `json:"disabled,string,omitempty" url:"disabled"`
	EnableSSL           bool   `json:"enableSSL,string,omitempty" url:"enableSSL"`
	Port                int    `json:"port,string,omitempty" url:"port,omitempty"`
	DedicatedIoThreads  int    `json:"dedicatedIoThreads,string,omitempty" url:"dedicatedIoThreads,omitempty"`
	MaxSockets          int    `json:"maxSockets,string,omitempty" url:"maxSockets,omitempty"`
	MaxThreads          int    `json:"maxThreads,string,omitempty" url:"maxThreads,omitempty"`
	UseDeploymentServer bool   `json:"useDeploymentServer,string,omitempty" url:"useDeploymentServer"`
}
