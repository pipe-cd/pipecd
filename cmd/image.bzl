def all_images():
    cmds = {
        "api": "pipecd-api",
        "api-debug": "pipecd-api-debug",
        "piped": "pipecd-piped",
        "web": "pipecd-web",
        "helloworld": "pipecd-helloworld",
    }
    images = {}

    for cmd, repo in cmds.items():
        images["$(DOCKER_REGISTRY)/%s:{STABLE_GIT_COMMIT_FULL}" % repo] = "//cmd/%s:image" % cmd
        images["$(DOCKER_REGISTRY)/%s:{STABLE_GIT_COMMIT}" % repo] = "//cmd/%s:image" % cmd

    return images
