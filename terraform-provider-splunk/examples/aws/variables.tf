variable "region" {
  description = "AWS region (Ex: us-west-2)"
  type        = string
}

variable "splunk_hec_endpoint" {
  description = "The splunk instance HEC URL (Ex: https://my-splunk.com:8088)"
  type        = string
}

variable "vpc_flow_logs_cloudwatch_log_group" {
  description = "The cloudwatch log group to which vpc flow logs are being sent"
  type        = string
}
