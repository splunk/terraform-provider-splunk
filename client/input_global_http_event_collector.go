package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func (client *Client) CreateGlobalHttpEventCollectorObject(httpInputConfigObj models.GlobalHttpEventCollectorObject) (*http.Response, error) {
	/*
		{"links":{"create":"/services/data/inputs/http/_new","_reload":"/services/data/inputs/http/_reload","_acl":"/services/data/inputs/http/_acl"},"origin":"https://localhost:8089/services/data/inputs/http","updated":"2020-07-27T13:57:48-07:00","generator":{"build":"a6754d8441bf","version":"8.0.3"},"entry":[{"name":"http","id":"https://localhost:8089/servicesNS/nobody/splunk_httpinput/data/inputs/http/http","updated":"2020-07-27T13:57:48-07:00","links":{"alternate":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http","list":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http","_reload":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http/_reload","edit":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http","disable":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http/disable"},"author":"nobody","acl":{"app":"splunk_httpinput","can_change_perms":true,"can_list":true,"can_share_app":true,"can_share_global":true,"can_share_user":false,"can_write":true,"modifiable":true,"owner":"nobody","perms":{"read":["*"],"write":["*"]},"removable":false,"sharing":"app"},"fields":{"required":[],"optional":["acceptFrom","ackIdleCleanup","allowQueryStringAuth","allowSslCompression","allowSslRenegotiation","caCertFile","caPath","channel_cookie","cipherSuite","crossOriginSharingHeaders","crossOriginSharingPolicy","dedicatedIoThreads","description","disabled","ecdhCurveName","ecdhCurves","enableSSL","forceHttp10","index","indexes","listenOnIPv6","maxEventSize","maxIdleTime","maxSockets","maxThreads","port","requireClientCert","sendStrictTransportSecurityHeader","serverCert","sourcetype","sslAltNameToCheck","sslCommonNameToCheck","sslKeysfile","sslKeysfilePassword","sslVersions","useACK","useDeploymentServer"],"wildcard":[".*"]},"content":{"_rcvbuf":1572864,"ackIdleCleanup":"true","allowSslCompression":"true","allowSslRenegotiation":"true","dedicatedIoThreads":"2","disabled":false,"eai:acl":null,"eai:appName":"splunk_httpinput","eai:userName":"admin","enableSSL":"1","enablessl":"false","host":"ajayaraman-MBP-6E14B","index":"default","indexes":[],"maxSockets":"0","maxThreads":"0","port":"8088","sslVersions":"*,-ssl2","useDeploymentServer":"0"}}],"paging":{"total":1,"perPage":30,"offset":0},"messages":[]}
	*/
	values, err := query.Values(&httpInputConfigObj)
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadGlobalHttpEventCollectorObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateGlobalHttpEventCollectorObject(httpInputConfigObj models.GlobalHttpEventCollectorObject) (*http.Response, error) {
	values, err := query.Values(&httpInputConfigObj)
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
