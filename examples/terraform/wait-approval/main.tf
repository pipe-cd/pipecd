terraform {
  backend "gcs" {
    bucket      = "pipecd-play-terraform-examples-backend"
    prefix      = "wait-approval"
    credentials = ".credentials/service-account.json"
  }
}

variable "project" {}
variable "content" {}

provider "google" {
  project     = var.project
  credentials = ".credentials/service-account.json"
}

resource "google_storage_bucket_object" "object" {
  name    = "wait-approval/${terraform.workspace}.txt"
  bucket  = "pipecd-play-terraform-examples"
  content = var.content
}
