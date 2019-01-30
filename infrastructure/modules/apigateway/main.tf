variable "service" {}
variable "stage" {}


resource "aws_api_gateway_rest_api" "restapi" {
  name = "${var.service}-${var.stage}"
}

resource "aws_api_gateway_deployment" "restapi" {
  depends_on = [
    "module.hello",
    "module.users",
    "module.collections",
    "module.auth",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  stage_name = "${var.stage}"

  variables {
    deployed_at = "${timestamp()}"
  }
}

resource "aws_api_gateway_method_settings" "restapi-settings" {
  depends_on = [ "aws_api_gateway_deployment.restapi" ]
  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  stage_name = "${var.stage}"
  method_path = "*/*"

  settings {
    metrics_enabled = true
    logging_level = "ERROR"
  }
}
