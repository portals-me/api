variable "authenticator_arn" {}

module "auth" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
  path_part = "auth"

  methods_count = 0
  methods = []
}

module "auth-signUp" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.auth.id}"
  path_part = "signUp"

  methods_count = 1
  methods = [
    {
      http_method = "POST"
      function_arn = "${var.authenticator_arn}"
    }
  ]
}

module "auth-signUp-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.auth-signUp.id}"
}

module "auth-signIn" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.auth.id}"
  path_part = "signIn"

  methods_count = 1
  methods = [
    {
      http_method = "POST"
      function_arn = "${var.authenticator_arn}"
    }
  ]
}

module "auth-signIn-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.auth-signIn.id}"
}

module "auth-twitter" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.auth.id}"
  path_part = "twitter"

  methods_count = 1
  methods = [
    {
      http_method = "POST"
      function_arn = "${var.authenticator_arn}"
    },
  ]
}

module "auth-twitter-get" {
  source = "lambda_api_req"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.auth.id}"
  resource_id = "${module.auth-twitter.id}"
  resource_path = "${module.auth-twitter.path}"
  path_part = "twitter"

  methods_count = 1
  methods = [
    {
      http_method = "GET"
      function_arn = "${var.authenticator_arn}"
    },
  ]
  request_parameters = {
    "method.request.querystring.oauth_token" = true
    "method.request.querystring.oauth_verifier" = true
  }
}

module "auth-twitter-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.auth-twitter.id}"
}
