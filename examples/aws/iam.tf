// --------
// Lambda transformation
# Role for the transformation Lambda function attached to the kinesis stream
resource "aws_iam_role" "kinesis_firehose_lambda" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY
}

data "aws_iam_policy_document" "lambda_policy_doc" {
  statement {
    actions = [
      "logs:GetLogEvents",
    ]

    resources = [
      "*",
    ]

    effect = "Allow"
  }

  statement {
    actions = [
      "firehose:PutRecordBatch",
    ]

    resources = [
      "*"
    ]
  }

  statement {
    actions = [
      "logs:PutLogEvents",
    ]

    resources = [
      "*",
    ]

    effect = "Allow"
  }

  statement {
    actions = [
      "logs:CreateLogGroup",
    ]

    resources = [
      "*",
    ]

    effect = "Allow"
  }

  statement {
    actions = [
      "logs:CreateLogStream",
    ]

    resources = [
      "*",
    ]

    effect = "Allow"
  }
}

resource "aws_iam_policy" "lambda_transform_policy" {
  name   = "lambda_function_policy"
  policy = data.aws_iam_policy_document.lambda_policy_doc.json
}

resource "aws_iam_role_policy_attachment" "lambda_policy_role_attachment" {
  role       = aws_iam_role.kinesis_firehose_lambda.name
  policy_arn = aws_iam_policy.lambda_transform_policy.arn
}

// Firehose
# Role for Kinesis Firehose and CWL subscription
resource "aws_iam_role" "kinesis_firehose_stream_assume_role" {
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Action": "sts:AssumeRole",
      "Effect": "Allow"
    }
  ]
}
POLICY
}

data "aws_iam_policy_document" "kinesis_firehose_policy_document" {
  statement {
    actions = [
      "s3:AbortMultipartUpload",
      "s3:GetBucketLocation",
      "s3:GetObject",
      "s3:ListBucket",
      "s3:ListBucketMultipartUploads",
      "s3:PutObject",
    ]

    resources = [
      aws_s3_bucket.kinesis_firehose_stream_bucket.arn,
      "${aws_s3_bucket.kinesis_firehose_stream_bucket.arn}/*",
    ]

    effect = "Allow"
  }

  statement {
    actions = [
      "lambda:InvokeFunction",
      "lambda:GetFunctionConfiguration",
    ]

    resources = [
      "${aws_lambda_function.lambda_kinesis_firehose_data_transformation.arn}:$LATEST",
    ]
  }

  statement {
    actions = [
      "logs:PutLogEvents",
    ]

    resources = [
      "*"
    ]

    effect = "Allow"
  }
}

resource "aws_iam_policy" "kinesis_firehose_iam_policy" {
  policy = data.aws_iam_policy_document.kinesis_firehose_policy_document.json
}

resource "aws_iam_role_policy_attachment" "kinesis_fh_role_attachment" {
  role       = aws_iam_role.kinesis_firehose_stream_assume_role.name
  policy_arn = aws_iam_policy.kinesis_firehose_iam_policy.arn
}

resource "aws_iam_role" "cloudwatch_to_firehose_trust" {
  name        = "CloudWatchToSplunkFirehoseTrust"
  description = "Role for CloudWatch Log Group subscription"

  assume_role_policy = <<ROLE
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "logs.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
ROLE

}

data "aws_iam_policy_document" "cloudwatch_to_fh_access_policy" {
  statement {
    actions = [
      "firehose:*",
    ]

    effect = "Allow"

    resources = [
      aws_kinesis_firehose_delivery_stream.firehose-stream-splunk.arn,
    ]
  }

  statement {
    actions = [
      "iam:PassRole",
    ]

    effect = "Allow"

    resources = [
      aws_iam_role.cloudwatch_to_firehose_trust.arn,
    ]
  }
}

resource "aws_iam_policy" "cloudwatch_to_fh_access_policy" {
  description = "Cloudwatch to Firehose Subscription Policy"
  policy      = data.aws_iam_policy_document.cloudwatch_to_fh_access_policy.json
}

resource "aws_iam_role_policy_attachment" "cloudwatch_to_fh" {
  role       = aws_iam_role.cloudwatch_to_firehose_trust.name
  policy_arn = aws_iam_policy.cloudwatch_to_fh_access_policy.arn
}
