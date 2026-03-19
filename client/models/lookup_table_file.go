package models

type LookupTableFile struct {
	App          string `json:"namespace,omitempty" url:"namespace,omitempty"`
	Owner        string `json:"owner,omitempty" url:"owner,omitempty"`
	FileName     string `json:"lookup_file,omitempty" url:"lookup_file,omitempty"`
	FileContents string `json:"contents,omitempty" url:"contents,omitempty"`
}
