data "aws_region" "current" {}
variable "authorizer_arn" {}

resource "aws_api_gateway_authorizer" "lambda_authorizer" {
  name = "${var.service}-${var.stage}-lambda_authorizer"
  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  authorizer_uri = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${var.authorizer_arn}/invocations"
  authorizer_credentials = "${aws_iam_role.invocation_role.arn}"
}

resource "aws_iam_role" "invocation_role" {
  name = "${var.service}-${var.stage}-api_gateway_auth_invocation"
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "invocation_policy" {
  name = "${var.service}-${var.stage}-default"
  role = "${aws_iam_role.invocation_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "lambda:InvokeFunction",
      "Effect": "Allow",
      "Resource": "${var.authorizer_arn}"
    }
  ]
}
EOF
}
