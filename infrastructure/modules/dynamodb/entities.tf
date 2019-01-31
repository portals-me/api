variable "stage" {}
variable "service" {}

resource "aws_dynamodb_table" "entities" {
  name = "${var.service}-${var.stage}-entities"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"
  range_key = "sort"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "sort"
    type = "S"
  }

  attribute {
    name = "data"
    type = "S"
  }

  global_secondary_index {
    name = "DataTable"
    hash_key = "sort"
    range_key = "data"
    projection_type = "ALL"
  }
}
