## 1.4.21
* Fix: pagerduty integration key, custom details.

## 1.4.20
* Support for defining pagerduty integration key, custom details.
* Fix: Better error handling for non 20X error codes.

## 1.4.19
* Support for Pager Duty fields in saved_searches

## 1.4.18
* Support for SNOW alert actions
* Prerequisite: Install ServiceNow Addon into Splunk instance.

## 1.4.17
* Support for XSOAR alert actions

## 1.4.16
* Fix incorrect revert in v1.4.15

## 1.4.15
* Support for jira service desk actions in saved_searches

## 1.4.14
* Fix: Omit auto_summarize field in saved_searches when empty

## 1.4.13
* Fix: configs_conf permits underscores in conf filename

## 1.4.12
* Fix: Don't read all searches just to find one search

## 1.4.11
* Fix: Don't read all views just to find one view

## 1.4.10
Role Capabilities are unordered (#95)

## 1.4.9 (Sep 29, 2021)
Handle missing SAML groups (#89)

## 1.4.8 (Aug 23, 2021)
* Added splunk_generic_acl resource

## 1.4.7 (Aug 06, 2021)
* Support for webhook alert action in saved_searches

## 1.4.6 (June 22, 2021)
* Fix:  Adds Cookie handling (fixes #49) (#75)
* Primarily helps with sending subsequent requests to the same SH when SH cluster is enabled with ELB.
* Example: With AWS, `lb_cookie_stickiness_policy` has to configured for requests to be sent to the same SH.

## 1.4.5 (June 14, 2021)
* Fix: Enabling to explicitly set values to roles attributes #76

## 1.4.4 (April 12, 2021)
* Fix: use_ack disable update failure

## 1.4.3 (April 6, 2021)
* Fix: State migration for alert_track #66

## 1.4.2 (April 1, 2021)
* Support to create indexes on Splunk Cloud (Beta)

## 1.4.1 (March 9, 2021)
* Fix: State not accurate w/r/t splunk_saved_searches->alert_track #65

## 1.4.0 (March 9, 2021)
* Data UI Views does not persist permissions on creation #59
* Updated Examples

## 1.3.9 (January 6, 2021)
* Support to create dashboards and views resource #45

## 1.3.8 (January 5, 2021)
* Fix: Unable to create a Metric index (#48)
* Fix: TestAccCreateSplunkIndex is randomly failing (#39)

## 1.3.7 (December 8, 2020)
* Fix: http client overrides default transport and no longer supports environment proxy settings. (#46)

## 1.3.6 (December 8, 2020)
* Fix: Pass HEC token as input with http event collector resource

## 1.3.5 (December 1, 2020)
Add email message field for reports & alerts #38

## 1.3.4 (November 18, 2020)
* Fix: Adding Slack actions to saved_searches resource #33
* Revert go mod path update

## 1.3.3 (November 13, 2020)
* Fix: Fix URL encoding for resource names #32
* Enhancements: Added linting `golangci-lint`
* go mod path update to `github.com/terraform-providers/terraform-provider-splunk`

## 1.3.2 (October 21, 2020)
* Fix: Feature Request for saved search to support additional attributes #24
* Fix: `saved_search` default is_visible to true #3

## 1.3.1 (October 6, 2020)
* Support for admin/SAML-groups API endpoint #23

## 1.3.0 (October 6, 2020)
* Fix: authorization/roles srchIndexesAllowed and srchIndexesDefault should be lists #20
* Fix: Support for saved_search argument dispatch.index_earliest #15
* github action workflow for integration tests
* Additional examples in the examples folder

## 1.2.1 (September 30, 2020)
* Bug fix for authorization header

## 1.2.0 (September 29, 2020)
* Change provider configuration - support for splunk auth token

## 1.1.0 (September 25, 2020)
* Change HTTP inputs resource attribute types (`use_ack` and `use_deployment_server`) to `int`
* Creating examples folder
* Adding AWS firehose + Splunk integration example under `examples/aws`

## 1.0.2 (September 21, 2020)
* Adding Github actions workflow

## 1.0.1 (September 21, 2020)
* Adding CHANGELOG.MD
