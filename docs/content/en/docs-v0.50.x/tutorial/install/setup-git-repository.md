---
title: "1. Setup Git Repository"
linkTitle: "1. Setup Git Repository"
weight: 1
description: >
  Prepare the Git repository for the tutorial.
---

# 1. Setup Git Repository

In this page, you will prepare the Git repository for the tutorial.

## Steps

1. Fork the [pipe-cd/tutorial](https://github.com/pipe-cd/tutorial) repository on GitHub.

2. Clone the forked repository to your local machine.

   ```sh
   git clone https://github.com/<YOUR_GITHUB_USERNAME>/tutorial.git
   cd tutorial
   ```

3. Confirm that you can push to the repository.

   ```sh
   git push origin main
   ```

## What Is in the Repository

The repository contains:

- `src/install/`: Configuration files for installing Control Plane and Piped.
  - `control-plane/docker-compose.yaml`: Docker Compose file for the Control Plane.
  - `piped/piped.yaml`: Configuration file for the Piped.
- `src/deploy/`: Configuration files for deploying applications.
  - Platform-specific directories (`kubernetes/`, `cloudrun/`, `ecs/`, `lambda/`, `terraform/`) with deployment configurations.

## See Also

- [Git Documentation](https://git-scm.com/doc)

---

[Next: 2. Install Control Plane >](../install-control-plane/)

[< Previous: Install](../)
