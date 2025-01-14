provider "docker" {
}

module "helloworld" {
  source = "helloworld"
  image_version = "v1.0.0"
  external_port = 8080
}
