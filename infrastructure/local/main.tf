provider "aws" {
  region = "ap-northeast-1"
  endpoints {
    dynamodb = "http://localhost:8000"
  }
}

module "dynamodb" {
  source = "../modules/dynamodb"
  stage = "test"
  service = "portals-me"

  entity-stream_arn = ""
  stream-count = 0
}
