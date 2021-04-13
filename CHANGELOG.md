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
