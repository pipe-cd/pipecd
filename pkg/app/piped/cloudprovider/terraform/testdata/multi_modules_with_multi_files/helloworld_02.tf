provider "docker" {
}

module "helloworld_02" {
  source = "helloworld"
  version = "v0.9.0"
  image_version = "v0.9.0"
  external_port = 8081
}
