variable "project" {}

provider "google" {
  project     = var.project
  credentials = ".credentials/service-account.json"
}

terraform {
  backend "gcs" {
    bucket      = "pipecd-terraform-examples"
    prefix      = "tfstates/secret-management"
    credentials = ".credentials/service-account.json"
  }
}

variable "content" {}

resource "google_storage_bucket_object" "object" {
  name    = "examples/secret-management/${terraform.workspace}.txt"
  bucket  = "pipecd-terraform-examples"
  content = var.content
}
