resource "docker_container" "helloworld" {
  name = "gcr.io/pipecd/helloworld:${var.image_version}"
  ports {
    internal = 9376
    external = "${var.external_port}"
  }
}

variable "external_port" {
    default = 80
}

variable "image_version" {
  default = "latest"
}

output "container" {
  value = docker_image.helloworld
}
