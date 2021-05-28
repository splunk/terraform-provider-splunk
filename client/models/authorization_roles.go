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
	CumulativeRTSrchJobsQuota int      `json:"cumulativeRTSrchJobsQuota" url:"cumulativeRTSrchJobsQuota"`
	CumulativeSrchJobsQuota   int      `json:"cumulativeSrchJobsQuota" url:"cumulativeSrchJobsQuota"`
	DefaultApp                string   `json:"defaultApp,omitempty" url:"defaultApp,omitempty"`
	ImportedRoles             []string `json:"imported_roles,omitempty" url:"imported_roles,omitempty"`
	Name                      string   `json:"name,omitempty" url:"name,omitempty"`
	RtSrchJobsQuota           int      `json:"rtSrchJobsQuota" url:"rtSrchJobsQuota"`
	SrchDiskQuota             int      `json:"srchDiskQuota" url:"srchDiskQuota"`
	SrchFilter                string   `json:"srchFilter,omitempty" url:"srchFilter,omitempty"`
	SrchIndexesAllowed        []string `json:"srchIndexesAllowed,omitempty" url:"srchIndexesAllowed,omitempty"`
	SrchIndexesDefault        []string `json:"srchIndexesDefault,omitempty" url:"srchIndexesDefault,omitempty"`
	SrchJobsQuota             int      `json:"srchJobsQuota" url:"srchJobsQuota"`
	SrchTimeWin               int      `json:"srchTimeWin" url:"srchTimeWin"`
}
