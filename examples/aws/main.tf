# This example terraform configuration will combine both AWS and Splunk providers to help ingest AWS services data into Splunk.
# The Splunk provider creates a new index, a HEC token and installs the AWS kinesis firehose TA with knowledge objects to detect AWS data sources.
# The AWS provider creates a firehose delivery stream with splunk as destination using the above HEC token, S3 configurations to create a splash bucket for failed events, lambda transformation required by firehose to format VPC flow logs correctly and a cloudwatch subscription filter.
# The iam roles required for all the AWS resources are in separate file. (iam.tf)
terraform {
  required_providers {
    splunk = {
      source  = "splunk/splunk"
      version = "1.4.25"
    }
  }
}

provider "aws" {
  // Environment variables used:
  // * Access Key ID:     AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
  // * Secret Access Key: AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
  region = var.region
}

provider "splunk" {
  // Provide splunk instance credentials and details either via resource block or env variables
  url                      = "localhost:8089"
  username                 = "admin"
  password                 = "changeme"
  insecure_skip_verify     = true
  ignore_schedule_priority = false
}

resource "splunk_indexes" "vpc-flow-logs-index" {
  name                   = "vpc-flow-logs"
  max_hot_buckets        = 6
  max_total_data_size_mb = 1000000
}

resource "splunk_apps_local" "Splunk_TA_aws-kinesis-firehose" {
  name             = "/splunk-add-on-for-amazon-kinesis-firehose_122.tgz" // Location of the app on the remote instance
  filename         = true
  explicit_appname = "Splunk_TA_aws-kinesis-firehose"
  update           = true
}

resource "splunk_global_http_event_collector" "http" {
  disabled   = false
  enable_ssl = true
  port       = 8088
}

resource "splunk_inputs_http_event_collector" "firehose" {
  name       = "firehose"
  index      = splunk_indexes.vpc-flow-logs-index.name
  sourcetype = "aws:firehose:json"
  disabled   = false
  use_ack    = 1
  token      = ""
  depends_on = [
    splunk_indexes.vpc-flow-logs-index,
    splunk_global_http_event_collector.http,
    splunk_apps_local.Splunk_TA_aws-kinesis-firehose
  ]
}

resource "aws_kinesis_firehose_delivery_stream" "firehose-stream-splunk" {
  name        = "firehose-to-splunk"
  destination = "splunk"

  splunk_configuration {
    hec_endpoint               = var.splunk_hec_endpoint
    hec_token                  = splunk_inputs_http_event_collector.firehose.token
    hec_acknowledgment_timeout = 180
    retry_duration             = 300
    s3_backup_mode             = "FailedEventsOnly"

    processing_configuration {
      enabled = true
      processors {
        type = "Lambda"
        parameters {
          parameter_name  = "LambdaArn"
          parameter_value = "${aws_lambda_function.lambda_kinesis_firehose_data_transformation.arn}:$LATEST"
        }
        parameters {
          parameter_name  = "RoleArn"
          parameter_value = aws_iam_role.kinesis_firehose_stream_assume_role.arn
        }
      }
    }
  }

  s3_configuration {
    bucket_arn = aws_s3_bucket.kinesis_firehose_stream_bucket.arn
    role_arn   = aws_iam_role.kinesis_firehose_stream_assume_role.arn
  }

  depends_on = [
    aws_s3_bucket.kinesis_firehose_stream_bucket,
    aws_iam_role.kinesis_firehose_stream_assume_role,
    splunk_inputs_http_event_collector.firehose
  ]
}

resource "aws_s3_bucket" "kinesis_firehose_stream_bucket" {
  bucket        = "firehose-to-splunk-splash-bucket"
  acl           = "private"
  force_destroy = true
}

resource "aws_cloudwatch_log_subscription_filter" "cloudwatch_subscription_filter" {
  name            = "vpc-flow-logs-to-cloudwatch-logs"
  log_group_name  = var.vpc_flow_logs_cloudwatch_log_group
  filter_pattern  = ""
  destination_arn = aws_kinesis_firehose_delivery_stream.firehose-stream-splunk.arn
  role_arn        = aws_iam_role.cloudwatch_to_firehose_trust.arn
  distribution    = "ByLogStream"
  depends_on = [
    aws_kinesis_firehose_delivery_stream.firehose-stream-splunk
  ]
}

resource "aws_lambda_function" "lambda_kinesis_firehose_data_transformation" {
  s3_bucket     = "trumpet-splunk-prod-${var.region}"
  s3_key        = "splunk_vpc_firehose_processor_v0.3.zip"
  function_name = "firehose-to-splunk-vpc-lambda-function"

  role    = aws_iam_role.kinesis_firehose_lambda.arn
  handler = "lambda_function.handler"
  runtime = "python3.7"
  timeout = 300
}
