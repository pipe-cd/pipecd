variable "project" {}
variable "region" {}
variable "credentials_file" {}

provider "google" {
  project     = var.project
  region      = var.region
  credentials = file(var.credentials_file)
}

terraform {
  backend "gcs" {
    bucket      = "pipecd-terraform-sample"
    prefix      = "tfstates/simple"
    credentials = "/etc/terraform/terraform_example_service_account"
  }
}

variable "example_bucket_name" {}

resource "google_storage_bucket_object" "simple" {
  name    = "${terraform.workspace}-simple.txt"
  bucket  = var.example_bucket_name
  content = "simple"
}
