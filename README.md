[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/7489/badge)](https://www.bestpractices.dev/projects/7489)
[![Build](https://github.com/pipe-cd/pipecd/actions/workflows/build.yaml/badge.svg)](https://github.com/pipe-cd/pipecd/actions/workflows/build.yaml)
[![Test](https://github.com/pipe-cd/pipecd/actions/workflows/test.yaml/badge.svg)](https://github.com/pipe-cd/pipecd/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipe-cd/pipecd)](https://goreportcard.com/report/github.com/pipe-cd/pipecd)
[![codecov](https://codecov.io/gh/pipe-cd/pipecd/branch/master/graph/badge.svg)](https://codecov.io/gh/pipe-cd/pipecd)
[![Release](https://img.shields.io/github/v/release/pipe-cd/pipecd?label=Release)](https://github.com/pipe-cd/pipecd/releases/latest)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](/LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-pipecd-informational.svg)](https://pipecd.dev/docs/)
[![Slack](https://img.shields.io/badge/Slack-%23pipecd-informational.svg)](https://app.slack.com/client/T08PSQ7BQ/C01B27F9T0X)
[![Twitter](https://img.shields.io/twitter/url/https/twitter.com/pipecd_dev.svg?style=social&label=Follow%20%40pipecd_dev)](https://twitter.com/pipecd_dev)
[![CNCF Status](https://img.shields.io/badge/cncf%20status-sandbox-lightgrey.svg)](https://www.cncf.io/projects/pipecd/)

<p align="center">
  <img src="https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/logo.png" width="180"/>
</p>

<p align="center">
  A GitOps style continuous delivery platform that provides consistent deployment and operations experience for any applications
  <br/>
  <a href="https://pipecd.dev"><strong>Explore PipeCD docs »</strong></a>
  <a href="https://play.pipecd.dev?project=play"><strong>Try on Playground »</strong></a>
</p>

---

![](https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/rolled-back-deployment.png)

### Overview

PipeCD provides __a unified continuous delivery solution for multiple application kinds on multi-cloud__ that empowers engineers to deploy faster with more confidence, a GitOps tool that enables doing deployment operations by pull request on Git.

![](https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/pipecd-explanation.png)

---

### Why PipeCD?

- Simple, unified and easy to use but powerful pipeline definition to construct your deployment
- Same deployment interface to deploy applications of any platform, including Kubernetes, Terraform, GCP Cloud Run, AWS Lambda, AWS ECS
- No CRD or applications' manifest changes are required; Only need a pipeline definition along with your application manifests
- No deployment credentials are exposed or required outside the application cluster
- Built-in deployment analysis as part of the deployment pipeline to measure impact based on metrics, logs, emitted requests
- Easy to interact with any CI; The CI tests and builds artifacts, PipeCD takes the rest
- Insights show metrics like lead time, deployment frequency, MTTR and change failure rate to measure delivery performance
- Designed to manage thousands of cross-platform applications in multi-cloud for company scale but also work well for small projects

For more details, explore the [PipeCD Documentation](https://pipecd.dev/docs).

---

### Getting Started

- The [quickstart guide](https://pipecd.dev/docs/quickstart/) shows how to set up PipeCD components and deploy a hello-world application with PipeCD for testing purposes.
- The [tutorial](https://github.com/pipe-cd/tutorial) shows how to run PipeCD locally for introduction.
- The [installation guide](https://pipecd.dev/docs/installation/) explains and helps set up PipeCD for your real-life use environment.

---

### Community and development

- Check out the [PipeCD website](https://pipecd.dev) for the complete documentation and helpful links.
- Join the [`#pipecd` channel](https://cloud-native.slack.com/archives/C01B27F9T0X) in the [CNCF Slack Workspace](https://communityinviter.com/apps/cloud-native/cncf) to get help and to discuss features.
- Tweet [@pipecd_dev](https://twitter.com/pipecd_dev) on Twitter/X.
- Create Github [Issues](https://github.com/pipe-cd/pipecd/issues) or [Discussions](https://github.com/pipe-cd/pipecd/discussions/) to report bugs or request features.
- Join the [PipeCD Development and Community Meetings](https://zoom-lfx.platform.linuxfoundation.org/meeting/96831504919?password=2f60b8ec-5896-40c8-aa1d-d551ab339d00) where we share the latest project news, demos, answer questions, and help triage issues. You can also view the [Meeting Notes and Agenda](https://bit.ly/pipecd-mtg-notes).

Participation in PipeCD project is governed by the CNCF [Code of Conduct](CODE_OF_CONDUCT.md).

---

### Contributing

We'd love you to join us! Please see the [Contributing Guide](CONTRIBUTING.md) to get started.

---

### Adopters

You can find a list of publicly recognized users of the PipeCD in the [ADOPTERS.md](ADOPTERS.md) file. We strongly encourage all PipeCD users to add their names to this list, as we love to see the community's growing success!

---

### Thanks to the contributors of PipeCD ❤️

<div align="center">
  <a href="https://github.com/pipe-cd/pipecd/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=pipe-cd/pipecd" alt="PipeCD Contributors"/>
  </a>
</div>

---

**We are a [Cloud Native Computing Foundation](https://cncf.io/) Sandbox Project.**

<img src="https://www.cncf.io/wp-content/uploads/2022/07/cncf-color-bg.svg" width=300 />

The Linux Foundation® (TLF) has registered trademarks and uses trademarks. For a list of TLF trademarks, see [Trademark Usage](https://www.linuxfoundation.org/trademark-usage/).

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpipe-cd%2Fpipecd.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpipe-cd%2Fpipecd?ref=badge_large)
