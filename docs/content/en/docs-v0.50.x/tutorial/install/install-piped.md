---
title: "3. Install Piped"
linkTitle: "3. Install Piped"
weight: 3
description: >
  Run a Piped agent on local as a single binary.
---

# 3. Install Piped

In this page, you will run a Piped agent on local as a single binary.

## Installation

1. Create a Piped ID on the console.

   1-1. Access the Piped settings page on the console: [http://localhost:8080/settings/piped](http://localhost:8080/settings/piped)

   1-2. Click `+ ADD`.
    ![ADD Piped](/images/tutorial/install/piped-add-button.png)

   1-3. Enter any `Name` and `Description`, and then click `SAVE`.
    ![Input piped](/images/tutorial/install/piped-add-input.png)

   1-4. You will see the Piped's ID and Key. **You will use them in the following steps.**
    ![Piped is registered](/images/tutorial/install/piped-registered.png)

2. Move on to the directory you cloned in [1. Setup Git Repository](../setup-git-repository/).

   ```console
   cd <YOUR_CLONED_TUTORIAL_DIR>
   cd src/install/piped/
   ```

3. Replace `<YOUR_PIPED_KEY>` below to the `Piped Key` value you got in `1-4.` and create a file of the key.[^1]

   ```console
   echo -n <YOUR_PIPED_KEY> > .pipedKey
   ```

4. Edit your `src/install/piped/piped.yaml` as below.
   - `spec.pipedID`: `Piped Id` you got in `1-4.`
   - `spec.pipedKeyFile`: `<YOUR_CLONED_TUTORIAL_DIR>/src/install/piped/.pipedKey`
   - `spec.repositories[0].remote`: The URL of your repository on GitHub you created in [1. Setup Git Repository](../setup-git-repository/).
   - `spec.platformProviders`: Uncomment one item. You will deploy to the environment.

5. Get the official binary of Piped from [pipecd/releases](https://github.com/pipe-cd/pipecd/releases/latest). Choose your OS and CPU architecture.

   ```console
   # OS=darwin or linux
   # CPU=arm64 or amd64
   curl -Lo ./piped  https://github.com/pipe-cd/pipecd/releases/download/v0.50.1/piped_v0.50.1_${OS}_${CPU}
   ```

6. Make the Piped executable.

   ```console
   chmod +x ./piped
   ```

7. Run the Piped.[^2]

   ```console
   ./piped piped --config-file=./piped.yaml --insecure
   ```

   <details>
   <summary>Troubleshooting (Click to expand)</summary>

   - `failed to create api client	{"error": "context deadline exceeded"}`
     - It means the piped failed to connect to the control plane.
     - Check `apiAddress` in `piped.yaml`
     - Check whether the control plane is serving on `apiAddress`.

   - `rpc error: code = Unauthenticated desc = Unauthenticated`
     - It means the piped connected to the control plane but failed to authenticate.
     - Check `projectID`, `pipedID`, `pipedKeyFile` in `piped.yaml`.
     - Check the value in the `pipedKeyFile`.

      
   - `failed to clone from remote`
     - It means the piped successfully connected to the control plane with authentication but failed to clone from GitHub.
     - Check  `spec.git.repositories`>`remote`,`branch`.

   </details>


## Confirmation

1. Access the Piped settings page on the console: [http://localhost:8080/settings/piped](http://localhost:8080/settings/piped)

    If successful, you will see a green mark(`Online`) in the `Name` area.
    If not green, please wait for a few minutes.
      ![Piped Status](/images/tutorial/install/piped-status.png)

## See Also

- [Install Piped](https://pipecd.dev/docs/installation/install-piped/)
- [Managing Piped](https://pipecd.dev/docs/user-guide/managing-piped/)
- [Configuration Reference of Piped](https://pipecd.dev/docs/user-guide/managing-piped/configuration-reference/)

[^1]: Do **NOT** push a Piped Key to remote because it is a credential. In actual use, you should store it and `piped.yaml` in a secret management tool.

[^2]: `--insecure` option disables TLS of the communication from the Piped to the Control Plane. (cf. [Piped's Runtime Options](https://pipecd.dev/docs/user-guide/managing-piped/runtime-options/))

---

[Next: Deploy >](../../tutorial/deploy/)

[< Previous: 2. Install Control Plane](../install-control-plane/)
