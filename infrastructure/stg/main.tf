terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "portals-me-tfstate"
    key = "stg/terraform.tfstate"
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

module "portals-me" {
  source = "../modules/portals-me"

  stage = "${var.apex_environment}"
  apex_function_hello = "${var.apex_function_hello}"
  apex_function_user = "${var.apex_function_user}"
  apex_function_collection = "${var.apex_function_collection}"
  apex_function_article = "${var.apex_function_article}"
  apex_function_authenticator = "${var.apex_function_authenticator}"
  apex_function_authorizer = "${var.apex_function_authorizer}"
  apex_function_entity-stream = "${var.apex_function_entity-stream}"
  apex_function_stream-activity-feed = "${var.apex_function_stream-activity-feed}"
  apex_function_stream-timeline-feed = "${var.apex_function_stream-timeline-feed}"
  apex_function_timeline = "${var.apex_function_timeline}"
}

variable "apex_environment" {}
variable "apex_function_hello" {}
variable "apex_function_authorizer" {}
variable "apex_function_user" {}
variable "apex_function_collection" {}
variable "apex_function_article" {}
variable "apex_function_authenticator" {}
variable "apex_function_entity-stream" {}
variable "apex_function_stream-activity-feed" {}
variable "apex_function_stream-timeline-feed" {}
variable "apex_function_timeline" {}
