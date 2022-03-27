---
date: 2022-03-27
title: "The PipeCD play environment is here!!!"
linkTitle: "The PipeCD play environment"
weight: 994
description: "In this post, we will have a glance at the play environment of PipeCD, how to access and what you can get from the environment."
author: Khanh Tran ([@khanhtc1202](https://twitter.com/khanhtc1202))
---

Good news, PipeCD team's happy to bring you a place where you can have a look at the PipeCD platform in use. We call it the `play environment` - [https://play.pipecd.dev](https://play.pipecd.dev).

![play-environment-view](/images/play-environment-overview.png)
<p style="text-align: center;">
PipeCD play environment view
</p>

With this live demo, you can now have a look at the PipeCD platform in use without self preparing [quickstart](/docs/quickstart/) or such. It’s way easier (and faster) to make a try by clicking around and seeing how you can get when using PipeCD as the CD platform on your own.

### How to access

The PipeCD play environment console is available at [https://play.pipecd.dev](https://play.pipecd.dev). After following the link, you will get the login page as below.

![play-environment-login](/images/play-environment-login.png)

Type __play__ to the input box as the project name to login and click to `Continue` to go to the account sign in page. Followed by `LOGIN WITH GITHUB`.

![play-environment-login-github](/images/play-environment-github-login.png)

Then, that's it. Feel free to click around and see what PipeCD can bring to you in real-life use.

#### Some pages you may feel interest

![play-environment-application](/images/play-environment-application.png)
<p style="text-align: center;">
Application detail page - show the application's state and info
</p>

![play-environment-deployment](/images/play-environment-deployment.png)
<p style="text-align: center;">
Deployment detail page - show the deployment's stages and its log
</p>

### Notes for the PipeCD play environment

1. Since the limitation of the resources, you can only log in with the [Viewer permission](/docs/operator-manual/control-plane/auth/#role-based-access-control-rbac). This means you can only click around and see PipeCD team prepared examples, __triggering new deployments or creating new resources are disabled__.
2. Currently, only applications of kinds: `KUBERNETES`, `CLOUDRUN` and `TERRAFORM` are prepared and available to see on the play environment console. We will add example applications of other kinds (`LAMBDA`, `ECS`, etc.) later.
3. The [application configuration](/docs/user-guide/adding-an-application/) files for those examples on the play environment are published at [pipe-cd/examples](https://github.com/pipe-cd/examples) repository. Those configuration files are __real-life useable__ configurations, but only on our cluster, since the encrypted credentials placed in those files are ours. You can replace those and use on your own.

Happy Pipecding 👋
