provider "aws" {
  region = "ap-northeast-1"
}

module "iam" {
  source = "../modules/iam"
}

module "apigateway" {
  source = "../modules/apigateway"

  service = "${var.service}"
  stage = "${var.apex_environment}"
  hello_arn = "${var.apex_function_hello}"
}

module "dynamodb" {
  source = "../modules/dynamodb"
  service = "${var.service}"
}
