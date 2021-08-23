- Start Date: 2021-06-18
- Target Version: 0.11.0

# Summary

This RFC proposes adding the ability to preview the deployment plan against a given commit before merging.
That helps developers review the pull request with a better experience and more confidence to approve it to trigger the actual deployments.

# Motivation

The mission of a CI/CD system is providing to developers the feedback related to the changes as quickly as possible.
Until now, PipeCD was only focusing on the feedback when the deployment is started running and after its completion, it also means after the commit got merged into the deployment branch.
Recently, we have received a lot of requests from our users that they wish they could have more feedback at the time pull request is created to help them review the pull request with more confidence.
So PipeCD team thinks it is time to find a way to provide the feedback earlier: during PR reviewing, before merging.

# Detailed design

Before going into the designing section, we will take a look at what feedback the developers want to have while reviewing the pull request.
Probably, there are many things else, but the following ones are considered currently:
- what application will be deployed after the pull request got merged
- what deployment strategy (`QUICK_SYNC` or `PIPELINE_SYNC`) will be used
- what resources will be added, deleted, or modified

And here are some real scenarios:

- Developer makes a pull request to update infrastructure code which are managed by Terraform, then they want to preview the `terraform plan` result when the pull request is created. It shows them that once the pull request gets merged which infrastructure component will be updated, destroyed or what will be newly created.
- Developer makes a pull request to update the using version of a remote Helm chart. In that case, Git's diff just shows a single line of version number change. Reviewing some kind of that pull request becomes hader. And they wish they could see the full list of resource changes.

Interacting with the pull request is the responsibility of a CI system, not CD system like PipeCD. So PipeCD will not interact directly with the pull request but provide a way for the CI job to fetch the plan-preview result. Specifically, we will add a new `plan-preview` command to our `pipectl` tool. That command communicates with PipeCD's control plane to retrieve the plan-preview result.

Then, developers can add a new CI job into their Git repository. The CI system triggers that job when a new pull request is created or updated to run `pipectl plan-preview` against the head commit of the pull request.

Although the developer can check the plan-preview result from the log of the CI job, we think the best way to send the feedback to the developer is directly commenting its content to their pull request. So we recommend doing that in your CI job. For GitHub users, PipeCD team will also provide a GitHub action that helps you enable this feature quickly and seamlessly.

**Architecture design:**

![](https://github.com/pipe-cd/pipe/blob/master/docs/static/images/rfc-plan-preview-architecture.png)

1. CI job or user runs `pipectl plan-preview` to request a plan-preview result for the head commit of the pull request.
2. Server lists all pipeds that are configured to handle the repository and emits a command for each piped into the datastore.
3. Piped fetches the unhandled commands to handle.
4. Piped generates plan-preview result based on the change between head commit and base commit.
5. Piped reports to control plane that command was handled and sends the plan-preview result to control plane.
6. Server marks the command in the datastore as handled and stores the plan-preview result into the filestore. 
7. After sending the plan-preview request, CI job waits and periodically checks the command result.
8. CI job logs the plan-preview result and comments its content on the pull request.

# Questions:

Q: What kind of application can have this feature?

A: All. (Kubernetes, Terraform, CloudRun, Lambda, Amazon ECS...)

# Unresolved questions

None.
