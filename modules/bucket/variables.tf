variable "project" {}

variable "name" {}

variable "location" {}

variable "force_destroy" {
  default = false
}

variable "topic_name" {}

variable "subscription_name" {
  default = "test_subscription"
}

variable "add_random_suffix" {
  default = true
}
