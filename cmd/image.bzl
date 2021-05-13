def all_images():
    cmds = {
        "piped": "piped",
        "pipecd": "pipecd",
        "pipectl": "pipectl",
        "helloworld": "helloworld",
    }
    images = {}

    for cmd, repo in cmds.items():
        images["$(DOCKER_REGISTRY)/%s:{STABLE_VERSION}" % repo] = "//cmd/%s:%s_app_image" % (cmd, repo)

    images["$(DOCKER_REGISTRY)/piped_okd:{STABLE_VERSION}"] = "//cmd/piped:piped_okd_app_image"
    return images
