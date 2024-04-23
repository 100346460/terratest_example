module "bucket_name" {
  source            = "../string_formatter"
  name              = var.name
  add_random_suffix =  var.add_random_suffix
}

module "pubsub_topic" {
  source  = "../pubsub"
  topic_name = var.topic_name
  project = var.project
}

resource "google_storage_bucket" "storage_bucket" {
  project       = var.project
  name          = module.bucket_name.result
  location      = upper(var.location)
  force_destroy = var.force_destroy
  uniform_bucket_level_access = true
}

resource "google_storage_notification" "notification" {
  bucket         = google_storage_bucket.storage_bucket.name
  payload_format = "JSON_API_V1"
  topic          = "projects/temporal-field-419719/topics/${module.pubsub_topic.topic_name}"
  event_types    = ["OBJECT_FINALIZE", "OBJECT_METADATA_UPDATE"]

}

resource "google_pubsub_subscription" "example_subscription" {
  project = var.project
  name  = var.subscription_name
  topic = module.pubsub_topic.topic_name
}


