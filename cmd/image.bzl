def all_images():
    targets = {
        "//cmd/piped:piped_app_image": "piped",
        "//cmd/piped:piped_okd_app_image": "piped-okd",
        "//cmd/pipecd:pipecd_app_image": "pipecd",
        "//cmd/pipectl:pipectl_app_image": "pipectl",
        "//cmd/helloworld:helloworld_app_image": "helloworld",
    }
    images = {}

    for target, repo in targets.items():
        images["$(DOCKER_REGISTRY)/%s:{STABLE_VERSION}" % repo] = target

    return images
