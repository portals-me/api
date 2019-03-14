
variable "service" {
  default = "portals-me"
}

variable "stage" {}

module "iam" {
  source = "../iam"
  stage = "${var.stage}"
  service = "${var.service}"
}

module "apigateway" {
  source = "../apigateway"

  service = "${var.service}"
  stage = "${var.stage}"
  authorizer_arn = "${var.apex_function_authorizer}"
  hello_arn = "${var.apex_function_hello}"
  user_arn = "${var.apex_function_user}"
  collection_arn = "${var.apex_function_collection}"
  article_arn = "${var.apex_function_article}"
  authenticator_arn = "${var.apex_function_authenticator}"
  timeline_arn = "${var.apex_function_timeline}"
}

module "dynamodb" {
  source = "../dynamodb"
  stage = "${var.stage}"
  service = "${var.service}"

  entity-stream_arn = "${var.apex_function_entity-stream}"
  stream-activity-feed_arn = "${var.apex_function_stream-activity-feed}"
  stream-timeline-feed_arn = "${var.apex_function_stream-timeline-feed}"
}
