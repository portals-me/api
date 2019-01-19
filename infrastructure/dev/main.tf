provider "aws" {
  region = "ap-northeast-1"
}

module "iam" {
  source = "../modules/iam"
}

module "dynamodb" {
  source = "../modules/dynamodb"
  service = "portals-me"
}

