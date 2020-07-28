def all_images():
    cmds = {
        "api": "api",
        "api-debug": "api-debug",
        "piped": "piped",
        "web": "web",
        "helloworld": "helloworld",
    }
    images = {}

    for cmd, repo in cmds.items():
        images["$(DOCKER_REGISTRY)/%s:{STABLE_VERSION}" % repo] = "//cmd/%s:image" % cmd
        images["$(DOCKER_REGISTRY)/%s:{STABLE_GIT_COMMIT}" % repo] = "//cmd/%s:image" % cmd

    return images
