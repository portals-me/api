variable "service" {
  default = "portals-me"
}

variable "stage" {
  default = "prod"
}

variable "aws_region" {}
variable "apex_environment" {}
variable "apex_function_hello" {}
variable "apex_function_authorizer" {}
variable "apex_function_user" {}
variable "apex_function_collection" {}
variable "apex_function_article" {}
variable "apex_function_authenticator" {}
