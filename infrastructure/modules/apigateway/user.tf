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
  methods_count = 2
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"

  request_parameters = {
    "method.request.path.userId" = true
  }

  methods = [
    {
      http_method = "GET"
      function_arn = "${var.user_arn}"
    },
    {
      http_method = "PUT"
      function_arn = "${var.user_arn}"
    },
  ]
}

module "users-user-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.users-user.id}"
}

module "users-user-feed" {
  source = "lambda_api_path"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.users-user.id}"
  path_part = "feed"
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
    },
  ]
}

module "users-user-feed-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.users-user-feed.id}"
}

module "users-user-follow" {
  source = "lambda_api_path"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.users-user.id}"
  path_part = "follow"
  methods_count = 1
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"

  request_parameters = {
    "method.request.path.userId" = true
  }

  methods = [
    {
      http_method = "POST"
      function_arn = "${var.user_arn}"
    },
  ]
}

module "users-user-follow-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.users-user-follow.id}"
}
