variable "timeline_arn" {}

module "timeline" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
  path_part = "timeline"
  methods_count = 1
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"

  methods = [
    {
      http_method = "GET"
      function_arn = "${var.timeline_arn}"
    },
  ]
}

module "timeline-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.timeline.id}"
}

