variable "topic_name" {}
variable "project" {}

resource "google_pubsub_topic" "pubsub_topic" {
  name = var.topic_name
  project = var.project

  labels = {
    foo = "bar"
  }

  message_retention_duration = "86600s"
}

output "topic_name"{
    value = google_pubsub_topic.pubsub_topic.name
}