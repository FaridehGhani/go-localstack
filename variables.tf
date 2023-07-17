variable "sqs" {
  description = "value of the name for aws sqs queue service"
  type = string
  default = "my-queue"
}

variable "dynamodb" {
  description = "value of the name for aws dynamodb"
  type = string
  default = "messages"
}