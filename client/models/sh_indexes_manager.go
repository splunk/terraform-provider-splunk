package models

type ShIndexesManagerResponse struct {
	Entry    []ShIndexesManagerEntry `json:"entry"`
	Messages []ErrorMessage          `json:"messages"`
}

type ShIndexesManagerEntry struct {
	Name    string                 `json:"name"`
	Id      string                 `json:"id"`
	ACL     ACLObject              `json:"acl"`
	Content ShIndexesManagerObject `json:"content"`
}

type ShIndexesManagerObject struct {
	Datatype               string `json:"datatype,omitempty" url:"datatype,omitempty"`
	FrozenTimePeriodInSecs string `json:"frozenTimePeriodInSecs,omitempty" url:"frozenTimePeriodInSecs,omitempty"`
	MaxGlobalRawDataSizeMB string `json:"maxGlobalRawDataSizeMB,omitempty" url:"maxGlobalRawDataSizeMB,omitempty"`
}
