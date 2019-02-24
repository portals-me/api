resource "aws_dynamodb_table" "feeds" {
  name = "${var.service}-${var.stage}-feeds"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "user_id"
  range_key = "timestamp"

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  attribute {
    name = "item_id"
    type = "S"
  }

  global_secondary_index {
    name = "ItemID"
    hash_key = "item_id"
    projection_type = "KEYS_ONLY"
  }
}

resource "aws_dynamodb_table" "timeline" {
  name = "${var.service}-${var.stage}-timeline"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"
  range_key = "timestamp"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  attribute {
    name = "item_id"
    type = "S"
  }

  global_secondary_index {
    name = "ItemID"
    hash_key = "item_id"
    projection_type = "KEYS_ONLY"
  }
}
