resource "aws_iam_role" "portals-me-default" {
  name = "portals-me-default"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

data "aws_iam_policy_document" "portals-me-default-policy-doc" {
  statement {
    actions = [
      "cognito-identity:*",
      "dynamodb:*",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role_policy" "portals-me-default-policy" {
  name = "portals-me-default-policy"
  role = "${aws_iam_role.portals-me-default.name}"
  policy = "${data.aws_iam_policy_document.portals-me-default-policy-doc.json}"
}

output "portals-me-default-role" {
  value = "${aws_iam_role.portals-me-default.arn}"
}
