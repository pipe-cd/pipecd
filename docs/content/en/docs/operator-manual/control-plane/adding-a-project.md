---
title: "Adding a project"
linkTitle: "Adding a project"
weight: 2
description: >
  This page describes how to set up a new project.
---

The control plane operator can add a new project for a team.
Project adding can be simply done from an internal web page prepared for the operator.
Because that web service is running in an `operator` pod, so in order to access it, using `kubectl port-forward` command to forward a local port to a port on the `operator` pod as following:

``` console
kubectl port-forward service/pipecd-operator 9082
```

Then, access to [http://localhost:9082](http://localhost:9082).

On that page, you will see the list of registered projects and a link to register new projects.
Registering a new project requires only a unique ID string and an optional description text.

Once a new project has been registered, a static admin (username, password) will be automatically generated for the project admin. You can send that information to the project admin. The project admin uses the provided static admin information to log in to PipeCD. After that, they can change static admin information, configure the SSO or disable static admin user.
