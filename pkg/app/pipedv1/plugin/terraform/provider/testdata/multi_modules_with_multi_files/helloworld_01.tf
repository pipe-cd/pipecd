provider "docker" {
}

module "helloworld_01" {
  source = "helloworld"
  version = "v1.0.0"
  image_version = "v1.0.0"
  external_port = 8080
}
