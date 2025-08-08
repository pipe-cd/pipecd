---
title: "PipeCD v1 Documentation Index"
linkTitle: "Documentation v1"
weight: 5
description: >
  Complete documentation for PipeCD v1 plugin architecture
---

# PipeCD v1 Documentation

Welcome to the comprehensive documentation for PipeCD v1, featuring the revolutionary plugin-based architecture that transforms how deployments are executed and managed.

## What's New in v1?

PipeCD v1 introduces a plugin-based architecture that enables:
- **Extensible Deployments**: Deploy to any platform through plugins
- **Community Ecosystem**: Use and contribute plugins developed by the community  
- **Flexible Architecture**: Mix and match plugins for custom deployment workflows
- **Multi-Language Support**: Write plugins in any language with gRPC support

## Documentation Structure

### Getting Started

- **[PipeCD v1 Overview](overview/)** - Comprehensive introduction to v1 architecture, benefits, and key concepts
- **[Quick Start Guide](quickstart/)** - Get up and running with PipeCD v1 in minutes
- **[Installation Guide](installation/)** - Detailed installation and setup instructions

### Core Concepts

- **[Plugin Architecture](architecture/)** - Deep dive into the plugin-based architecture
- **[Configuration Reference](configuration-reference/)** - Complete configuration documentation
- **[Plugin System](plugin-system/)** - Understanding how plugins work and interact

### Migration and Upgrade

- **[Migration Guide](migration-guide/)** - Step-by-step migration from v0 to v1
- **[Compatibility Matrix](compatibility/)** - Version compatibility and support lifecycle
- **[Troubleshooting](troubleshooting/)** - Common issues and solutions

### Plugin Development

- **[Plugin Development Guide](plugin-development/)** - Complete guide to developing plugins
- **[SDK Reference](sdk-reference/)** - Official SDK documentation and API reference
- **[Plugin Examples](plugin-examples/)** - Real-world plugin implementation examples

### Built-in Plugins

- **[Kubernetes Plugin](plugins/kubernetes/)** - Container deployment with progressive delivery
- **[Terraform Plugin](plugins/terraform/)** - Infrastructure as Code deployments
- **[Lambda Plugin](plugins/lambda/)** - AWS Lambda function deployments
- **[Wait Plugin](plugins/wait/)** - Time-based deployment stages
- **[ScriptRun Plugin](plugins/scriptrun/)** - Custom script execution

### Advanced Topics

- **[Multi-Cluster Deployments](advanced/multi-cluster/)** - Deploy across multiple environments
- **[Custom Stages](advanced/custom-stages/)** - Implement custom deployment stages
- **[Security Best Practices](advanced/security/)** - Secure plugin and deployment configuration
- **[Performance Optimization](advanced/performance/)** - Optimize plugin and deployment performance

### Operations

- **[Monitoring and Observability](operations/monitoring/)** - Monitor v1 deployments and plugins
- **[Backup and Recovery](operations/backup/)** - Backup strategies for v1 configurations
- **[Scaling and High Availability](operations/scaling/)** - Scale Piped agents and plugins

### Community

- **[Community Plugins](community/plugins/)** - Discover and contribute community plugins
- **[Contributing Guide](community/contributing/)** - How to contribute to PipeCD v1
- **[Plugin Registry](community/registry/)** - Official plugin registry and marketplace

## Quick Navigation

### I want to...

**Get started with PipeCD v1**
→ [Quick Start Guide](quickstart/) → [Installation Guide](installation/)

**Migrate from v0 to v1**
→ [Migration Guide](migration-guide/) → [Compatibility Matrix](compatibility/)

**Develop a custom plugin**
→ [Plugin Development Guide](plugin-development/) → [SDK Reference](sdk-reference/)

**Deploy Kubernetes applications**
→ [Kubernetes Plugin](plugins/kubernetes/) → [Configuration Examples](configuration-reference/#kubernetes-plugin-configuration)

**Deploy infrastructure with Terraform**
→ [Terraform Plugin](plugins/terraform/) → [Infrastructure Examples](examples/terraform/)

**Understand the architecture**
→ [PipeCD v1 Overview](overview/) → [Plugin Architecture](architecture/)

**Troubleshoot issues**
→ [Troubleshooting Guide](troubleshooting/) → [Common Issues](troubleshooting/#common-issues)

## Version Information

- **Current Version**: v1.0.0
- **API Version**: v1alpha1
- **Compatibility**: Control Plane v0.52.0+
- **v0 Support**: Until December 31, 2025

## Key Features

### 🔌 Plugin Architecture
- **Extensible**: Add support for any deployment platform
- **Modular**: Use only the plugins you need
- **Isolated**: Plugin failures don't affect other deployments

### 🌍 Multi-Platform Support
- **Kubernetes**: Progressive delivery with canary/blue-green strategies
- **Terraform**: Infrastructure as Code with approval workflows
- **Cloud Platforms**: AWS Lambda, Google Cloud Run, Azure Functions
- **Custom Platforms**: Develop plugins for any platform

### 🛠️ Developer Experience
- **Configuration-as-Code**: Version control all deployment configurations
- **GitOps Workflow**: Trigger deployments through Git commits
- **Progressive Delivery**: Built-in canary and blue-green deployment strategies
- **Rollback Support**: Automatic and manual rollback capabilities

### 📊 Observability
- **Deployment Insights**: Track DORA metrics and deployment performance
- **Real-time Monitoring**: Monitor deployments and plugin health
- **Audit Logging**: Complete audit trail of all deployment activities
- **Integration Ready**: Works with existing monitoring and alerting systems

## Recent Updates

### v1.0.0 (Latest)
- ✅ Plugin architecture general availability
- ✅ Kubernetes plugin with full progressive delivery support
- ✅ Terraform plugin with approval workflows
- ✅ Migration tools for v0 to v1 transition
- ✅ Community plugin repository launch

### Upcoming Features
- 🔄 Plugin hot reloading
- 🔐 Plugin signature verification
- 📦 Enhanced plugin registry
- 🎯 Resource-based RBAC for plugins

## Getting Help

### Documentation
- **User Guides**: Step-by-step instructions for common tasks
- **API Reference**: Complete API and configuration documentation
- **Examples**: Real-world configuration examples and use cases

### Community Support
- **Slack**: [#pipecd](https://cloud-native.slack.com/archives/C01B27F9T0X) - General discussion and support
- **Slack**: [#pipecd-plugin-dev](https://cloud-native.slack.com/) - Plugin development discussion
- **GitHub**: [Issues](https://github.com/pipe-cd/pipecd/issues) - Bug reports and feature requests
- **GitHub**: [Discussions](https://github.com/pipe-cd/pipecd/discussions) - Community Q&A

### Professional Support
- **Office Hours**: Monthly community calls with maintainers
- **Training**: Official PipeCD training and certification programs
- **Consulting**: Professional services for enterprise deployments

## Contributing

PipeCD v1 is an open-source project and we welcome contributions:

- **Code**: Contribute to core platform and plugins
- **Documentation**: Improve guides and examples
- **Plugins**: Develop and share community plugins
- **Testing**: Help test new features and report issues
- **Feedback**: Share your experience and suggestions

See our [Contributing Guide](community/contributing/) for detailed information.

## License

PipeCD is released under the [Apache License 2.0](https://github.com/pipe-cd/pipecd/blob/master/LICENSE).

---

**Ready to get started?** Begin with our [Quick Start Guide](quickstart/) or explore the [PipeCD v1 Overview](overview/) to learn more about the architecture.

**Migrating from v0?** Check out our comprehensive [Migration Guide](migration-guide/) for a smooth transition.

**Want to develop plugins?** Dive into our [Plugin Development Guide](plugin-development/) and start building.
