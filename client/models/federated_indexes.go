package models

// FederatedIndexResponse is the JSON response returned by the federated index
// REST endpoint.
type FederatedIndexResponse struct {
	Entry    []FederatedIndexEntry `json:"entry"`
	Messages []ErrorMessage        `json:"messages"`
}

// FederatedIndexEntry represents one federated index definition.
type FederatedIndexEntry struct {
	Name    string               `json:"name"`
	Content FederatedIndexObject `json:"content"`
}

// FederatedIndexObject maps a local federated index to a provider dataset.
type FederatedIndexObject struct {
	FederatedProvider string `json:"federated.provider,omitempty" url:"federated.provider,omitempty"`
	FederatedDataset  string `json:"federated.dataset,omitempty" url:"federated.dataset,omitempty"`
}
