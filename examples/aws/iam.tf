resource "aws_iam_role" "firehose-stream-role" {
  name = "firehose-stream-role"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "firehose.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
}
EOF
}


resource "aws_iam_role_policy" "firehose-stream-policy" {
  name = "firehose-stream-policy"
  role = aws_iam_role.firehose-stream-role.id
  depends_on = [
	aws_iam_role.firehose-stream-role
  ]

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "kinesis:*",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": [
        "arn:aws:s3:::firehose-to-splunk-splash-bucket",
        "arn:aws:s3:::firehose-to-splunk-splash-bucket/*"
      ]
    }
  ]
}
EOF
}
