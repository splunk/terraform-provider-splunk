package service

import (
	"github.com/splunk/go-splunkd/model"
	"github.com/splunk/go-splunkd/util"
)

// AuthorizationService implements a new service type
type AuthorizationService service

// Login gets a session ID for use in subsequent API calls that require authentication
// It also sets SessionKey fields in Client
func (as *AuthorizationService) Login() (model.SessionKey, error) {
	var sessionKey model.SessionKey
	loginURL := as.client.BuildSplunkdURL(nil, "services", "auth", "login")
	bodyValues := map[string]string{
		"username": as.client.Auth[0],
		"password": as.client.Auth[1],
	}
	response, err := as.client.Post(loginURL, bodyValues)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return sessionKey, err
	}
	err = util.ParseResponse(&sessionKey, response)
	if err == nil {
		as.client.SessionKey = sessionKey.SessionKey
	}
	return sessionKey, err
}
