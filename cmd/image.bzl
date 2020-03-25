def all_images():
    # TODO: Auto detect the list of images.
    cmds = {
        "api": "pipe-api",
        "livelog": "pipe-livelog",
        "runner": "pipe-runner",
    }
    images = {}

    for cmd, repo in cmds.items():
        images["$(DOCKER_REGISTRY)/%s:{STABLE_GIT_COMMIT_FULL}" % repo] = "//cmd/%s:image" % cmd
        images["$(DOCKER_REGISTRY)/%s:{STABLE_GIT_COMMIT}" % repo] = "//cmd/%s:image" % cmd

    return images
