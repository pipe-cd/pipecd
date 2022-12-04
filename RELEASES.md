# Release Process
This document explains the process to release new version and to address patch release.

## Major release
This refers to the release of new features.

### Confirm the changelog and Create Release Note
- Run the release script

  This example assumes that `vX.Y.Z` will be released:
  ```shell
  make gen/release version=vX.Y.Z
  ````

  `RELEASE` file will be updated and `docs/content/en/blog/releases/vX.Y.Z.md` file will be created.

- Push the above changes and Create a pull request to confirm the changelog.
  You can confirm the changelog through the reviewing comment in pull request by GitHub Actions.
  For more information, Please see [actions-gh-release](https://github.com/pipe-cd/actions-gh-release).

- Update the content in `docs/content/en/blog/releases/vX.Y.Z.md` file based on changelog results.
  Please refer to [this example](https://github.com/pipe-cd/pipecd/pull/3792/commits/2b59f7f2a492405bf6df905b3823b417e4a10c3e).

  It is recommended to commit the above changes once here.

### Generate document for new version
- Run the release document script

  This example assumes that `vX.Y.Z` will be released:
  ```shell
  make gen/release-docs version=vX.Y.Z
  ````

- Make a pull request to `master` branch with the above changes and get reviews and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on the master branch.

- Create a tagged release from the master branch. The release should start with "v" and be followed by the version number.

## Patch release
This refers to the release of critical bug fixes. \
A bugfix for a functional issue (not a data loss or security issue) that only affects an alpha feature does not qualify as a critical bug fix.

### Prerequisites
- `gh` is needed to be installed and ran `gh auth login`. Please refer to [cli/cli](https://github.com/cli/cli).

### Fix bugs
- Create a pull request to fix a bug on the `master` branch.

- Get reviews and merge.

### Confirm the changelog and Create Release Note
- As well as [Major release](https://github.com/pipe-cd/pipecd/blob/master/RELEASES.md#confirm-the-changelog-and-create-release-note), create a pull request to create a release note on the `master` branch.

- Get a review and merge.

### Backport fixes and Release note
- Run the cherry-pick script

  This example assumes that the name of a release branch is `release-vX.Y.x` and the numbers of pull request are `#1234` and `#5678`:
  ```shell
  make gen/cherry-pick branch=release-vX.Y.x pull_numbers=1234 5678
  ````

- Get a review and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on the release branch.

- Create a tagged release from the release branch `release-vX.Y.x`. The release should start with "v" and be followed by the version number.
