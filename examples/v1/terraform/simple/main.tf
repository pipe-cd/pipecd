terraform {
  # piped runs every deployment in a temporary directory and removes it afterwards,
  # so the default local state next to the .tf files would be lost each time.
  # Keep the state and the generated file in a fixed path instead. For real
  # infrastructure use a remote backend such as GCS or S3.
  backend "local" {
    path = "/tmp/pipecd-terraform-simple/terraform.tfstate"
  }

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4"
    }
  }
}

variable "content" {
  description = "Text written to the generated file."
  type        = string
  default     = "Hello from the PipeCD v1 Terraform plugin!"
}

# A cloud-free resource so the example runs on any piped with the terraform
# plugin loaded - no provider credentials or remote backend required.
resource "local_file" "hello" {
  filename = "/tmp/pipecd-terraform-simple/hello.txt"
  content  = var.content
}

output "file_path" {
  value = local_file.hello.filename
}
