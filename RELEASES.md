# Release Process
*This release process have not been fixed yet, hence this is subject to frequent change.*

This document explains the process to release new version and to address stable release.

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

### Create Document
- Run the release document script

  This example assumes that `vX.Y.Z` will be released:
  ```shell
  make release-docs version=vX.Y.Z
  ````

- Push the above changes and get a review and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on master branch.

- Create a tagged release. The release should start with "v" and be followed by the version number.

- Create a branch from the tagged release, e.g. "release-0.1". This will be used for the stable release.

## Stable release
This refers to the release of critical bug fixes.
A bugfix for a functional issue (not a data loss or security issue) that only affects an alpha feature does not qualify as a critical bug fix.

### Fix bugs
- Create a pull request to fix a bug on the `master` branch.

- Get a review and merge.

### Backport fixes
- Cherry-pick the original commit to the latest release branch.
  Release branches have a name of `release-MAJOR.MINO`.

- Get a review and merge.

### Confirm the changelog and Create Release Note
- As well as [Major release](https://github.com/pipe-cd/pipecd/blob/master/RELEASES.md#confirm-the-changelog-and-create-release-note), create a pull request to create a release note on the `master` branch.

- Get a review and merge.

### Backport Release Note
- Cherry-pick the original commit of creating release note to the latest release branch.
  Release branches have a name of `release-MAJOR.MINO`.

- Get a review and merge.

### Cut a new release
- Before cutting a new release, wait for all jobs in GitHub Actions to pass on release branch.

- Create a tagged release. The release should start with "v" and be followed by the version number.
