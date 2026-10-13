package splunk

import (
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// search head cluster members replicate objects asynchronously, so a request right after a write can fail
const defaultReplicationTimeout = time.Minute

func retryOnStatus(timeout time.Duration, statusCode int, request func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	err := resource.Retry(timeout, func() *resource.RetryError {
		r, err := request()
		if err == nil {
			resp = r
			return nil
		}
		if hasStatus(r, statusCode) {
			log.Printf("[DEBUG] %s. Retrying...", err)
			r.Body.Close()
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// readKnowledgeObject retries 404 only after create, on refresh 404 means the object was deleted
func readKnowledgeObject(d *schema.ResourceData, read func() (*http.Response, error)) (*http.Response, error) {
	if !d.IsNewResource() {
		return read()
	}
	return retryOnStatus(d.Timeout(schema.TimeoutCreate), http.StatusNotFound, read)
}

func hasStatus(resp *http.Response, statusCode int) bool {
	return resp != nil && resp.StatusCode == statusCode
}
