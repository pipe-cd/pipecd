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

### Confirm the changelog and Create Release Note
- Run the release script

  This example assumes that `vX.Y.Z` will be released:
  ```shell
  make release/init version=vX.Y.Z
  ````

  `RELEASE` file will be updated.

- Push the above changes and Create a pull request to confirm the changelog.
  You can confirm the changelog through the reviewing comment in pull request by GitHub Actions.
  For more information, Please see [actions-gh-release](https://github.com/pipe-cd/actions-gh-release).

### Generate document for new version
- Run the release document script

  This example assumes that `vX.Y.Z` will be released:
  ```shell
  make release/docs version=vX.Y.Z
  ````

- Make a pull request to `master` branch with the above changes and get reviews and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on the master branch.

- Create a tagged release from the master branch. The release should start with "v" and be followed by the version number.

## Minor release
This refers to the release of new features.

Please refer to [Major release](https://github.com/pipe-cd/pipecd/blob/master/RELEASES.md#major-release) for the processes.

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
- Run release pick commits

  This example assumes that the name of a release branch is `release-vX.Y.x` and the numbers of pull request are `#1234` and `#5678`:
  ```shell
  make release/pick branch=release-vX.Y.x pull_numbers="1234 5678"
  ````

- Get a review and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on the release branch.

- Create a tagged release from the release branch `release-vX.Y.x`. The release should start with "v" and be followed by the version number.
