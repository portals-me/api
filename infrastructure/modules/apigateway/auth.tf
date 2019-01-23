variable "auth_arn" {}

module "auth" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
  path_part = "auth"

  methods_count = 0
  methods = []
}
