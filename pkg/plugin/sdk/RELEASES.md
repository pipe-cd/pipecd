# Releases

This document explains the process to release new versions.

## Versioning

The piped plugin SDK are versioned independently from [PipeCD](https://github.com/pipe-cd/pipecd). This means that the plugin version is not related to the PipeCD version.

The piped plugin SDK version follow semantic versioning (same as PipeCD).

## Release Cycle

There is no release cycle for piped plugin SDK. Maintainers team will release new versions when there are new features or bug fixes.

## Cut a new release

Piped plugin SDK releases is done by [release sdk github actions](https://github.com/pipe-cd/piped-plugin-sdk-go/blob/main/.github/workflows/release.yaml).

### Inputs

| Name       | Description                                   | Required |  Example      |
|------------|-----------------------------------------------|:--------:|:-------------:|
| version    | The version of the piped plugin SDK to release.        |    yes   | v0.1.0    |
