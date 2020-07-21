package model

// SearchResponse represents the /search/jobs/{sid}/events or /search/jobs/{sid}/results response
type SearchResponse struct {
	Preview     bool                     `json:"preview"`
	InitOffset  int                      `json:"init_offset"`
	Messages    []interface{}            `json:"messages"`
	Results     []map[string]interface{} `json:"results"`
	Fields      []map[string]interface{} `json:"fields"`
	Highlighted map[string]interface{}   `json:"highlighted"`
}
