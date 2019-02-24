variable "service" {
  default = "portals-me"
}

variable "stage" {
  default = "dev"
}

variable "aws_region" {}
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
