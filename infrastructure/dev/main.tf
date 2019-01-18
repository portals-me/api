provider "aws" {
  region = "ap-northeast-1"
}

module "iam" {
  source = "../modules/iam"
}
