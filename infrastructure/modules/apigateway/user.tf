variable "user_arn" {}

module "users" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
  path_part = "users"
  methods_count = 1
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"

  methods = [
    {
      http_method = "GET"
      function_arn = "${var.user_arn}"
    }
  ]
}

module "users-user" {
  source = "lambda_api_path"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.users.id}"
  path_part = "{userId}"
  methods_count = 1
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"

  request_parameters = {
    "method.request.path.userId" = true
  }

  methods = [
    {
      http_method = "GET"
      function_arn = "${var.user_arn}"
    }
  ]
}