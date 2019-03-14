variable "service" {}
variable "stage" {}


resource "aws_api_gateway_rest_api" "restapi" {
  name = "${var.service}-${var.stage}"
}

resource "aws_api_gateway_deployment" "restapi" {
  depends_on = [
    "aws_api_gateway_rest_api.restapi",
    "module.hello",
    "module.collections",
    "module.auth",
    "module.timeline",
    "module.users",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.restapi.id}"
  stage_name = "${var.stage}"

  variables {
    deployed_at = "${timestamp()}"
  }

  lifecycle {
    create_before_destroy = true
  }
}
