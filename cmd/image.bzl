def all_images():
    cmds = {
        "piped": "piped",
        "piped": "piped_okd",
        "pipecd": "pipecd",
        "pipectl": "pipectl",
        "helloworld": "helloworld",
    }
    images = {}

    for cmd, repo in cmds.items():
        images["$(DOCKER_REGISTRY)/%s:{STABLE_VERSION}" % repo] = "//cmd/%s:%s_app_image" % (cmd, repo)

    return images
