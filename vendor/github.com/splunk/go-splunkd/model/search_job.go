package model

import (
	"time"
)

type searchJobService interface {
	NewSearch(config SearchConfig) (*SearchJob, error)
	GetSearchEventsBySid(sid string, config EventFetchConfig) (SearchResponse, error)
	GetSearchResultsBySid(sid string, config ResultFetchConfig) (SearchResponse, error)
	GetSearchJobs(params *PaginationParams) ([]SearchJob, error)
	GetSearchJobBySid(sid string) (SearchJobsResponse, error)
	UpdateSearchJob(sid string, key string, value string) error
	DeleteSearchJob(sid string) error
	ControlSearchJob(sid string, action Action) error
}

// States
// TODO: this could still be a type, just need to implement json unmarshalling
const (
	QUEUED     string = "QUEUED"
	PARSING    string = "PARSING"
	RUNNING    string = "RUNNING"
	PAUSED     string = "PAUSED"
	FINALIZING string = "FINALIZING"
	FAILED     string = "FAILED"
	DONE       string = "DONE"
)

// SearchContext specifies the user and app context for a search job
type SearchContext struct {
	User string
	App  string
}

// SearchJob specifies the fields returned for a /search/jobs/ entry for a specific job
type SearchJob struct {
	Sid           string           `json:"sid"`
	Content       SearchJobContent `json:"content"`
	SearchService searchJobService
	Context       *SearchContext
}

// Action controls a search
type Action string

// Valid values for search control actions
const (
	Pause          Action = "pause"
	Unpause        Action = "unpause"
	Finalize       Action = "finalize"
	Cancel         Action = "cancel"
	Touch          Action = "touch"
	SetTTL         Action = "setttl"
	SetPriority    Action = "setpriority"
	EnablePreview  Action = "enablepreview"
	DisablePreview Action = "disablepreview"
)

// NewSearchJob creates a SearchJob
func NewSearchJob(service searchJobService) *SearchJob {
	return &SearchJob{SearchService: service}
}

// NewNamespacedSearchJob creates a SearchJob with user and app context
func NewNamespacedSearchJob(service searchJobService, user, app string) *SearchJob {
	return &SearchJob{SearchService: service, Context: &SearchContext{User: user, App: app}}
}

// WaitForCompletion polls the job until it's completed or errors out
func (job *SearchJob) WaitForCompletion() error {
	var err error
	var done bool
	for !done {
		done, err = job.IsDone()
		if err != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}

// IsDone checks if the SearchJob is in a DONE state
func (job *SearchJob) IsDone() (bool, error) {
	if job.GetDispatchState() != DONE && job.GetDispatchState() != FAILED {
		err := job.Refresh()
		if err != nil {
			return false, err
		}
	}
	return job.GetDispatchState() == DONE, nil
}

// Refresh refreshes the SearchJob state
func (job *SearchJob) Refresh() error {
	jobResponse, err := job.SearchService.GetSearchJobBySid(job.Sid)
	if err != nil {
		return err
	}
	if len(jobResponse.Entry) > 0 {
		job.Content = jobResponse.Entry[0].Content
	}
	return nil
}

// Update updates properties on a SearchJob
func (job *SearchJob) Update(key string, value string) error {
	return job.SearchService.UpdateSearchJob(job.Sid, key, value)
}

// Delete deletes the underlying search job
func (job *SearchJob) Delete() error {
	return job.SearchService.DeleteSearchJob(job.Sid)
}

// Pause pauses the underlying search job
func (job *SearchJob) Pause() error {
	return job.SearchService.ControlSearchJob(job.Sid, Pause)
}

// GetEvents retrieves events for the given search job
func (job *SearchJob) GetEvents(count uint, offset uint) (SearchResponse, error) {
	pagination := PaginationParams{
		Count:  count,
		Offset: offset,
	}
	return job.SearchService.GetSearchEventsBySid(job.Sid, EventFetchConfig{Pagination: pagination})
}

// GetResults retrieves results for the given search job
func (job *SearchJob) GetResults(count uint, offset uint) (SearchResponse, error) {
	pagination := PaginationParams{
		Count:  count,
		Offset: offset,
	}
	return job.SearchService.GetSearchResultsBySid(job.Sid, ResultFetchConfig{Pagination: pagination})
}

// GetEventCount retrieves the event count for the SearchJob
func (job *SearchJob) GetEventCount() int {
	return job.Content.EventCount
}

// GetDispatchState retrieves the dispatch state for the SearchJob
func (job *SearchJob) GetDispatchState() string {
	if job.Content.DispatchState == "" {
		return QUEUED
	}
	return job.Content.DispatchState
}

// Cancel cancels the job
func (job *SearchJob) Cancel() error {
	return job.SearchService.ControlSearchJob(job.Sid, Cancel)
}
