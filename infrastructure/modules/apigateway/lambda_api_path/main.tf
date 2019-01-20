data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

variable "rest_api_id" {}
variable "parent_id" {}
variable "path_part" {}
variable "methods_count" {}

variable "authorization" {
  type = "string"
  default = "NONE"
}

variable "authorizer_id" {
  type = "string"
  default = ""
}

variable "request_parameters" {
  type = "map"
}


variable "methods" {
  type = "list"
}


resource "aws_api_gateway_resource" "main_path" {
  rest_api_id = "${var.rest_api_id}"
  parent_id = "${var.parent_id}"
  path_part = "${var.path_part}"
}

resource "aws_api_gateway_method" "main_path" {
  count = "${var.methods_count}"
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${aws_api_gateway_resource.main_path.id}"
  http_method = "${lookup(var.methods[count.index], "http_method")}"
  authorization = "${var.authorization}"
  authorizer_id = "${var.authorizer_id}"

  request_parameters = "${var.request_parameters}"
}

resource "aws_api_gateway_integration" "main_path" {
  depends_on = [
    "aws_api_gateway_method.main_path",
  ]

  count = "${var.methods_count}"
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${aws_api_gateway_resource.main_path.id}"
  http_method = "${lookup(var.methods[count.index], "http_method")}"
  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${lookup(var.methods[count.index], "function_arn")}/invocations"
}

resource "aws_lambda_permission" "main_path" {
  count = "${var.methods_count}"
  action = "lambda:InvokeFunction"
  function_name = "${lookup(var.methods[count.index], "function_arn")}"
  principal = "apigateway.amazonaws.com"
  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.rest_api_id}/*/${lookup(var.methods[count.index], "http_method")}${aws_api_gateway_resource.main_path.path}"
}

output "id" {
  value = "${aws_api_gateway_resource.main_path.id}"
}

output "path" {
  value = "${aws_api_gateway_resource.main_path.path}"
}
