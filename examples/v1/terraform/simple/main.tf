terraform {
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
  filename = "${path.module}/hello-${terraform.workspace}.txt"
  content  = var.content
}

output "file_path" {
  value = local_file.hello.filename
}
