package splunk

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/splunk/terraform-provider-splunk/client"
)

// search head cluster members replicate objects asynchronously, so a request right after a write can fail
const defaultReplicationTimeout = time.Minute

// replicationRefreshTimeout is how long a refresh or delete retries a 404 before treating the object as gone.
// Short enough that a missing object still leaves plan, long enough to ride out one lagged cluster member.
// The wait ends at this deadline. Polls use a short delay capped by the time remaining.
var replicationRefreshTimeout = 15 * time.Second

const (
	replicationRetryMinDelay = 500 * time.Millisecond
	replicationRetryMaxDelay = 2 * time.Second
)

// createDeadlines records the create budget for the ResourceData currently being applied.
// Terraform applies resources concurrently, so the key is that ResourceData pointer.
var createDeadlines sync.Map

func startCreateDeadline(d *schema.ResourceData) func() time.Duration {
	deadline := time.Now().Add(d.Timeout(schema.TimeoutCreate))
	createDeadlines.Store(d, deadline)
	return func() time.Duration {
		if left := time.Until(deadline); left > 0 {
			return left
		}
		return 0
	}
}

func stopCreateDeadline(d *schema.ResourceData) {
	createDeadlines.Delete(d)
}

func replicationTimeout(d *schema.ResourceData) time.Duration {
	if !d.IsNewResource() {
		return replicationRefreshTimeout
	}
	if v, ok := createDeadlines.Load(d); ok {
		if left := time.Until(v.(time.Time)); left > 0 {
			return left
		}
		return 0
	}
	return d.Timeout(schema.TimeoutCreate)
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

func retryOnStatus(timeout time.Duration, statusCode int, request func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	err := retryFor(timeout, func() error {
		closeBody(resp)
		resp = nil

		r, reqErr := request()
		resp = r
		return reqErr
	}, func(error) bool {
		return hasStatus(resp, statusCode)
	})
	// The last response stays open so the caller can read the status and close the body.
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// retryUpdateAcl retries only a 404 from UpdateAcl's HTTPStatusError.
// A 400 or 403 from the ACL update fails the apply on the first response.
func retryUpdateAcl(timeout time.Duration, update func() error) error {
	return retryFor(timeout, update, func(err error) bool {
		var statusErr *client.HTTPStatusError
		return errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusNotFound
	})
}

// retryFor retries while retryable reports true and the deadline has not passed.
// The next wait is capped by the time remaining, so the loop returns at the deadline
// instead of sleeping through a longer backoff.
func retryFor(timeout time.Duration, call func() error, retryable func(error) bool) error {
	if timeout < 0 {
		timeout = 0
	}
	deadline := time.Now().Add(timeout)
	delay := replicationRetryMinDelay
	for {
		err := call()
		if err == nil || !retryable(err) {
			return err
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return err
		}
		wait := delay
		if wait > remaining {
			wait = remaining
		}
		log.Printf("[DEBUG] %s. Retrying...", err)
		time.Sleep(wait)
		if delay < replicationRetryMaxDelay {
			delay *= 2
			if delay > replicationRetryMaxDelay {
				delay = replicationRetryMaxDelay
			}
		}
	}
}

// readKnowledgeObject retries 404 after create for the remainder of timeouts.create.
// On refresh it retries for replicationRefreshTimeout, then the caller drops state if the object is still missing.
func readKnowledgeObject(d *schema.ResourceData, read func() (*http.Response, error)) (*http.Response, error) {
	return retryOnStatus(replicationTimeout(d), http.StatusNotFound, read)
}

func hasStatus(resp *http.Response, statusCode int) bool {
	return resp != nil && resp.StatusCode == statusCode
}
