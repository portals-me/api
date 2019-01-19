provider "aws" {
  region = "ap-northeast-1"
}

module "iam" {
  source = "../modules/iam"
}

module "apigateway" {
  source = "../modules/apigateway"

  aws_region = "${var.aws_region}"
  service = "${var.service}"
  stage = "${var.apex_environment}"
  authorizer_arn = "${var.apex_function_authorizer}"
  hello_arn = "${var.apex_function_hello}"
  user_arn = "${var.apex_function_user}"
}

module "dynamodb" {
  source = "../modules/dynamodb"
  service = "${var.service}"
}
