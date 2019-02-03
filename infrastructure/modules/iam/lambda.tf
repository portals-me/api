variable "stage" {}
variable "service" {}

resource "aws_iam_role" "default" {
  name = "${var.service}-${var.stage}-default"

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

data "aws_iam_policy_document" "default-policy-doc" {
  statement {
    actions = [
      "cognito-identity:*",
      "dynamodb:*",
      "lambda:InvokeFunction",
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

resource "aws_iam_role_policy" "default-policy" {
  name = "${var.service}-${var.stage}-default-policy"
  role = "${aws_iam_role.default.name}"
  policy = "${data.aws_iam_policy_document.default-policy-doc.json}"
}

output "default-role" {
  value = "${aws_iam_role.default.arn}"
}
