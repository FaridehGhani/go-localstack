provider "aws" {
  access_key    = "dummy"
  secret_key = "dummy"
  region        = "us-east-1"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    sqs = "http://localhost:4566"
    dynamodb       = "http://localhost:4566"
  }
}

resource "aws_sqs_queue" "my_queue" {
  name                      = var.sqs
  visibility_timeout_seconds = 30
}

resource "aws_dynamodb_table" "messages" {
  name           = var.dynamodb
  read_capacity  = "20"
  write_capacity = "20"
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "Description"
    type = "S"
  }

  attribute {
    name = "CreatedAt"
    type = "S"
  }

  global_secondary_index {
    name               = "DescriptionIndex"
    hash_key           = "Description"
    projection_type    = "ALL"
    read_capacity      = 5
    write_capacity     = 5
  }

  global_secondary_index {
    name               = "CreatedAtIndex"
    hash_key           = "CreatedAt"
    projection_type    = "ALL"
    read_capacity      = 5
    write_capacity     = 5
  }
}
