provider "aws" {
  access_key    = "dummy"
  secret_key = "dummy"
  region        = "us-east-1"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    sqs = "http://localhost:4566"
  }
}

resource "aws_sqs_queue" "my_queue" {
  name                      = "my-queue"
  visibility_timeout_seconds = 30
}
