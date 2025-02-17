# Releases
This document explains the process to release new versions.

## Versioning
Versions are expressed as `vX.Y.Z`;

- `X` is the major version
- `Y` is the minor version
- `Z` is the patch version

## Release Cycle
- Major releases: arbitrary timing
- Minor releases: roughly every 2 months
- Patch releases: roughly every 2-3 weeks

Note: The team can release Release candidates(vX.Y.Z-rcXYZ) for versions at any time for early access/testing.

## Major release
This refers to the release of new features or breaking changes.

### Confirm the changelog and generate docs
- Run `make release`.

  This example assumes that `vX.Y.Z` will be released:
  ```shell
  make release version=vX.Y.Z
  ```

  The `RELEASE` file will be updated and docs `vX.Y.x` will be generated.

- Push the above changes and create a pull request to confirm the changelog.
  You can confirm the changelog through the reviewing comment in the pull request by [actions-gh-release](https://github.com/pipe-cd/actions-gh-release).

- Get reviews and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on the master branch.

- Create a tagged release from the master branch. The release should start with "v" and be followed by the version number.

## Minor release
This refers to the release of new features (breaking change may be included).

Please refer to [Major release](https://github.com/pipe-cd/pipecd/blob/master/RELEASES.md#major-release) for the processes.

## Patch release
This refers to the release of critical bug fixes. \
A bugfix for a functional issue (not a data loss or security issue) that only affects an alpha feature does not qualify as a critical bug fix.

This may also contain some minor features, but ensure that it does NOT contain any breaking changes.

### Prerequisites
- `gh` is needed to be installed and ran `gh auth login`. Please refer to [cli/cli](https://github.com/cli/cli).

### Fix bugs
- Create a pull request to fix a bug on the `master` branch.

- Get reviews and merge.

### Confirm the changelog and create Release Note

- Run `make release/init`.
  ```shell
  make release/init version=vX.Y.Z
  ```

  The `RELEASE` file will be updated.

- (Optional) if the patch contains changes with docs update, also need to run `make release/docs`
  ```shell
  make release/docs version=vX.Y.Z
  ```

  Note: You can use `make release version=vX.Y.Z` command to perform both init and docs sync tasks.

- Push the above changes and create a pull request to `master` to confirm the changelog.

- Get a review and merge.

### Backport fixes and Release Note

- Put labels of `cherry-pick` and `vX.Y.Z` to the PR of updating the `RELEASE` file to prevent conflicts.
- Run `cherry_pick` workflow
  - Label the merged PR you want to cherry pick with `cherry-pick` , `vX.Y.Z`
    (e.g. v0.48.6 https://github.com/pipe-cd/pipecd/pulls?q=is%3Apr+label%3Acherry-pick+is%3Aclosed+label%3Av0.48.6)
  - Execute the `cherry_pick` GitHub workflow with `release branch` and `release version` on master branch.
    (e.g. if you want to release v0.48.6, `release branch` is `release-v0.48.x` and `release version` is `v0.48.6`)

- If you have some trouble with the above, run release pick commits on local machine.

  This example assumes that the name of a release branch is `release-vX.Y.x` and the numbers of pull request are `#1234` and `#5678`:
  ```shell
  make release/pick branch=release-vX.Y.x pull_numbers="1234 5678"
  ````

- Get a review and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on the release branch.

- Create a tagged release from the release branch `release-vX.Y.x`. The release should start with "v" and be followed by the version number.

## RC Release

1. Prepare: Ensure all changes you want to attach to the release are available in the release target branch (`master` for the minor, `release-vX.Y.x` for the patch). For the patch, please refer to [Backport fixes and Release note](#backport-fixes-and-release-note)
2. Move on to [Releases > Draft a New Release](https://github.com/pipe-cd/pipecd/releases/new).
3. Set values as below:
   1. `Choose a tag`: Create a new tag `vX.Y.Z-rcN`
   2. `Target`(branch): use `master` for the minor rc, use `release-vX.Y.x` for the patch rc
   3. `Release title`: `Release vX.Y.Z-rcN`
   4. Body area
      1. Copy from the previous rc-release note.
      2. Modify the version of `> Note: This is a candidate release of vX.Y.Z.` if needed.
      3. Modify the version of  `## Changes since`.
      4. List Changes:
         1. Extract commits to include by the below commands.

            ```zsh
            $ PREVIOUS_TAG=v0.46.0-rc1 # Set the previous release tag
            $ git log $PREVIOUS_TAG..HEAD --oneline  | awk '{$1=""; print substr($0, 2)}'
            ```

            output(from newer to older):

            ```
            Add reference to the blog that shows how to install control plane on ECS (#4746)
            Update copyright (#4745)
            ...
            Add docs for SCRIPT_RUN stage (#4734)
            ```

         2. Classify the changes into 'Notable Changes' and 'Internal Changes'.
         3. Write them to the body area.
   5. **Select `Set as a pre-release`**, not `Set as the latest release`.
4. Push `Publish Release`.
