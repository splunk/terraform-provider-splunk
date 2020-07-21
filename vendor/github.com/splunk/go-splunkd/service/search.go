package service

import (
	"github.com/splunk/go-splunkd/model"
	"github.com/splunk/go-splunkd/util"
)

// SearchService implements a new service type
type SearchService service

// NewSearch dispatches a new spl search
func (service *SearchService) NewSearch(config model.SearchConfig) (*model.SearchJob, error) {
	var job = model.NewSearchJob(service)
	jobURL := service.client.BuildSplunkdURL(nil, "services",
		"search", "jobs")
	response, err := service.client.Post(jobURL, config)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	err = util.ParseResponse(&job, response)
	return job, err
}

// NewNamespacedSearch dispatches a new spl search with user and app context
func (service *SearchService) NewNamespacedSearch(config model.SearchConfig, user, app string) (*model.SearchJob, error) {
	var job = model.NewNamespacedSearchJob(service, user, app)
	if len(user) == 0 {
		user = "nobody"
	}
	if len(app) == 0 {
		app = "search"
	}
	jobURL := service.client.BuildSplunkdURL(nil, "servicesNS", user, app, "search", "jobs")
	response, err := service.client.Post(jobURL, config)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	err = util.ParseResponse(job, response)
	return job, err
}

// NewOneShotSearch dispatches a new spl search and fetches result back synchronously
func (service *SearchService) NewOneShotSearch(config model.SearchConfig) (*model.SearchResponse, error) {
	var searchModel model.SearchResponse
	config.ExecuteMode = "oneshot"

	jobURL := service.client.BuildSplunkdURL(nil, "services",
		"search", "jobs")

	response, err := service.client.Post(jobURL, config)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	err = util.ParseResponse(&searchModel, response)
	return &searchModel, err
}

// GetSearchEventsBySid fetches events by search sid
func (service *SearchService) GetSearchEventsBySid(sid string, config model.EventFetchConfig) (model.SearchResponse, error) {
	var searchModel model.SearchResponse
	jobURL := service.client.BuildSplunkdURL(util.ParseURLParams(config), "services",
		"search", "jobs", sid, "events")
	response, err := service.client.Get(jobURL)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return searchModel, err
	}
	err = util.ParseResponse(&searchModel, response)
	return searchModel, err
}

// GetSearchResultsBySid fetches results by search sid
func (service *SearchService) GetSearchResultsBySid(sid string, config model.ResultFetchConfig) (model.SearchResponse,
	error) {
	var searchModel model.SearchResponse
	jobURL := service.client.BuildSplunkdURL(util.ParseURLParams(config), "services",
		"search", "jobs", sid, "results")
	response, err := service.client.Get(jobURL)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return searchModel, err
	}
	err = util.ParseResponse(&searchModel, response)
	return searchModel, err
}

// GetSearchJobs gets details of current searches for a given count and offset
func (service *SearchService) GetSearchJobs(params *model.PaginationParams) ([]model.SearchJob, error) {
	var jobs []model.SearchJob

	if params == nil {
		params = model.NewDefaultPaginationParams()
	}

	jobsURL := service.client.BuildSplunkdURL(util.ParseURLParams(*params), "services",
		"search", "jobs")

	response, err := service.client.Get(jobsURL)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	var jobsResponse model.SearchJobsResponse
	err = util.ParseResponse(&jobsResponse, response)

	if err == nil {
		for _, jobEntry := range jobsResponse.Entry {
			job := model.SearchJob{Sid: jobEntry.Content.Sid, Content: jobEntry.Content, SearchService: service}
			jobs = append(jobs, job)
		}
	}
	return jobs, err
}

// GetSearchJobBySid will retrieve a job's status by its sid
func (service *SearchService) GetSearchJobBySid(sid string) (model.SearchJobsResponse, error) {
	var jobsResponse model.SearchJobsResponse
	jobURL := service.client.BuildSplunkdURL(nil,
		"services", "search", "jobs", sid)

	response, err := service.client.Get(jobURL)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return jobsResponse, err
	}
	err = util.ParseResponse(&jobsResponse, response)
	return jobsResponse, err
}

// UpdateSearchJob will update the status of an existing search job
func (service *SearchService) UpdateSearchJob(sid string, key string, value string) error {
	jobURL := service.client.BuildSplunkdURL(nil,
		"services", "search", "jobs", sid)

	response, err := service.client.Post(jobURL, map[string]string{key: value})
	if response != nil {
		defer response.Body.Close()
	}
	return err
}

// DeleteSearchJob will delete an existing search job
func (service *SearchService) DeleteSearchJob(sid string) error {
	jobURL := service.client.BuildSplunkdURL(nil,
		"services", "search", "jobs", sid)

	response, err := service.client.Delete(jobURL)
	if response != nil {
		defer response.Body.Close()
	}
	return err
}

// ControlSearchJob runs a job control command for the search by sid
func (service *SearchService) ControlSearchJob(sid string, action model.Action) error {
	jobURL := service.client.BuildSplunkdURL(nil, "services",
		"search", "jobs", sid, "control")

	response, err := service.client.Post(jobURL, map[string]string{"action": string(action)})
	if response != nil {
		defer response.Body.Close()
	}
	return err
}
