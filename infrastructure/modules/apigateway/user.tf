variable "user_arn" {}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_api_gateway_resource" "users" {
  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
  path_part = "users"
}

resource "aws_api_gateway_method" "users" {
  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  resource_id = "${aws_api_gateway_resource.users.id}"
  http_method = "GET"
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"
}

resource "aws_api_gateway_integration" "users" {
  depends_on = [
    "aws_api_gateway_method.users",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  resource_id = "${aws_api_gateway_resource.users.id}"
  http_method = "GET"
  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${var.user_arn}/invocations"
}

resource "aws_lambda_permission" "users" {
  statement_id = "AllowExecutionFromAPIGateway"
  action = "lambda:InvokeFunction"
  function_name = "${var.user_arn}"
  principal = "apigateway.amazonaws.com"
  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.restapi.id}/*/GET${aws_api_gateway_resource.users.path}"
}


output "id" {
  value = "${aws_api_gateway_resource.users.id}"
}

output "path" {
  value = "${aws_api_gateway_resource.users.path}"
}
