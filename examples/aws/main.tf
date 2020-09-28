# This example terraform configuration will combine both AWS and Splunk providers to help ingest AWS services data into Splunk.
# The Splunk provider creates a new index and an HEC input.
# The AWS provider creates a firehose delivery stream with splunk as destination using the above HEC input token and
# S3 splash bucket for failed events.

provider "aws" {
  // Environment variables used:
  // * Access Key ID:     AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
  // * Secret Access Key: AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
  region = "us-west-2"
}

provider "splunk" {
  url                  = "localhost:8089"
  username             = "admin"
  password             = "changeme"
  insecure_skip_verify = true
  // Or use environment variables used:
  // SPLUNK_USERNAME
  // SPLUNK_PASSWORD
  // SPLUNK_URL
  // SPLUNK_INSECURE_SKIP_VERIFY (Defaults to true)
}

resource "splunk_indexes" "aws-firehose-index" {
  name                   = "aws-firehose"
  max_hot_buckets        = 6
  max_total_data_size_mb = 1000000
}

resource "splunk_global_http_event_collector" "http" {
  disabled   = false
  enable_ssl = true
  port       = 8088
}

resource "splunk_inputs_http_event_collector" "firehose" {
  name       = "firehose"
  index      = "aws-firehose"
  sourcetype = "_json"
  disabled   = false
  use_ack    = 1
  token      = ""
  depends_on = [
    splunk_indexes.aws-firehose-index,
    splunk_global_http_event_collector.http,
  ]
}

resource "aws_kinesis_firehose_delivery_stream" "firehose-stream-splunk" {
  destination = "splunk"
  name        = "firehose-to-splunk"
  splunk_configuration {
    hec_endpoint               = "https://localhost:8088"
    hec_token                  = splunk_inputs_http_event_collector.firehose.token
    hec_acknowledgment_timeout = 180
    retry_duration             = 300
    s3_backup_mode             = "FailedEventsOnly"
  }

  s3_configuration {
    bucket_arn = aws_s3_bucket.firehose-splunk-s3.arn
    role_arn   = aws_iam_role.firehose-stream-role.arn
  }

  depends_on = [
    aws_s3_bucket.firehose-splunk-s3,
    aws_iam_role.firehose-stream-role,
    splunk_inputs_http_event_collector.firehose
  ]
}

resource "aws_s3_bucket" "firehose-splunk-s3" {
  bucket        = "firehose-to-splunk-splash-bucket"
  force_destroy = true
}
