data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

variable "rest_api_id" {}
variable "resource_id" {}
variable "resource_path" {}
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


variable "methods" {
  type = "list"
}

variable "request_parameters" {
  type = "map"
}


resource "aws_api_gateway_method" "main" {
  count = "${var.methods_count}"
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${lookup(var.methods[count.index], "http_method")}"
  authorization = "${var.authorization}"
  authorizer_id = "${var.authorizer_id}"
  request_parameters = "${var.request_parameters}"
}

resource "aws_api_gateway_integration" "main" {
  count = "${var.methods_count}"
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${lookup(var.methods[count.index], "http_method")}"
  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${lookup(var.methods[count.index], "function_arn")}/invocations"
}

resource "aws_lambda_permission" "main" {
  count = "${var.methods_count}"
  action = "lambda:InvokeFunction"
  function_name = "${lookup(var.methods[count.index], "function_arn")}"
  principal = "apigateway.amazonaws.com"
  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.rest_api_id}/*/${lookup(var.methods[count.index], "http_method")}${var.resource_path}"
}
