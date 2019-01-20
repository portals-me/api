variable "collection_arn" {}

module "collections" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${aws_api_gateway_rest_api.restapi.root_resource_id}"
  path_part = "collections"
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"

  methods_count = 2
  methods = [
    {
      http_method = "GET"
      function_arn = "${var.collection_arn}"
    },
    {
      http_method = "POST"
      function_arn = "${var.collection_arn}"
    },
  ]
}

module "collections-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.collections.id}"
}

module "collections-collection" {
  source = "lambda_api"

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  parent_id = "${module.collections.id}"
  path_part = "{collectionId}"
  methods_count = 1
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.lambda_authorizer.id}"

  methods = [
    {
      http_method = "GET"
      function_arn = "${var.collection_arn}"
    }
  ]
}

module "collections-collection-cors" {
  source = "github.com/squidfunk/terraform-aws-api-gateway-enable-cors"
  version = "0.2.0"

  api_id          = "${aws_api_gateway_rest_api.restapi.id}"
  api_resource_id = "${module.collections-collection.id}"
}
