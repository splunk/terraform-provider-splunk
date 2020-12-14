package models

type ClusterManagerResponse struct {
	Entry    []ClusterManagerEntry `json:"entry"`
	Messages []ErrorMessage        `json:"messages"`
}

type ClusterManagerEntry struct {
	Name    string               `json:"name"`
	ACL     ACLObject            `json:"acl"`
	Content ClusterManagerObject `json:"content"`
}

type ClusterManagerObject struct {
	AvailableSites      string `json:"available_sites,omitempty" url:"available_sites,omitempty"`
	ClusterLabel        string `json:"cluster_label,omitempty" url:"cluster_label,omitempty"`
	ConnectionTimeout   int    `json:"cxn_timeout,omitempty" url:"cxn_timeout,omitempty"`
	HeartbeatTimeout    int    `json:"heartbeat_timeout,omitempty" url:"heartbeat_timeout,omitempty"`
	Mode                string `json:"mode,omitempty" url:"mode,omitempty"`
	Multisite           bool   `json:"multisite,omitempty" url:"multisite,omitempty"`
	ReplicationFactor   int    `json:"replication_factor,omitempty" url:"replication_factor,omitempty"`
	RestartTimeout      int    `json:"restart_timeout,omitempty" url:"restart_timeout,omitempty"`
	SearchFactor        int    `json:"search_factor,omitempty" url:"search_factor,omitempty"`
	UseBatchMaskChanges bool   `json:"use_batch_mask_changes,omitempty" url:"use_batch_mask_changes,omitempty"`
}
