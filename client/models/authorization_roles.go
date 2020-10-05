package models

// Scripted Input Response Schema
type AuthorizationRolesResponse struct {
	Entry    []AuthorizationRolesEntry `json:"entry"`
	Messages []ErrorMessage            `json:"messages"`
}

type AuthorizationRolesEntry struct {
	Name    string                   `json:"name"`
	ACL     ACLObject                `json:"acl"`
	Content AuthorizationRolesObject `json:"content"`
}

type AuthorizationRolesObject struct {
	Capabilities              []string `json:"capabilities,omitempty" url:"capabilities,omitempty"`
	CumulativeRTSrchJobsQuota int      `json:"cumulativeRTSrchJobsQuota,omitempty" url:"cumulativeRTSrchJobsQuota,omitempty"`
	CumulativeSrchJobsQuota   int      `json:"cumulativeSrchJobsQuota,omitempty" url:"cumulativeSrchJobsQuota,omitempty"`
	DefaultApp                string   `json:"defaultApp,omitempty" url:"defaultApp,omitempty"`
	ImportedRoles             []string `json:"imported_roles,omitempty" url:"imported_roles,omitempty"`
	Name                      string   `json:"name,omitempty" url:"name,omitempty"`
	RtSrchJobsQuota           int      `json:"rtSrchJobsQuota,omitempty" url:"rtSrchJobsQuota,omitempty"`
	SrchDiskQuota             int      `json:"srchDiskQuota,omitempty" url:"srchDiskQuota,omitempty"`
	SrchFilter                string   `json:"srchFilter,omitempty" url:"srchFilter,omitempty"`
	SrchIndexesAllowed        []string `json:"srchIndexesAllowed,omitempty" url:"srchIndexesAllowed,omitempty"`
	SrchIndexesDefault        []string `json:"srchIndexesDefault,omitempty" url:"srchIndexesDefault,omitempty"`
	SrchJobsQuota             int      `json:"srchJobsQuota,omitempty" url:"srchJobsQuota,omitempty"`
	SrchTimeWin               int      `json:"srchTimeWin,omitempty" url:"srchTimeWin,omitempty"`
}
