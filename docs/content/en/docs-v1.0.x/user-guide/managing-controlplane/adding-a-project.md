---
title: "Adding a project"
linkTitle: "Adding a project"
weight: 2
description: >
  This page describes how to set up a new project.
---

The control plane ops can add a new project for a team.
Project adding can be simply done from an internal web page prepared for the ops.
Because that web service is running in an `ops` pod, so in order to access it, using `kubectl port-forward` command to forward a local port to a port on the `ops` pod as following:

``` console
kubectl port-forward service/pipecd-ops 9082 --namespace={NAMESPACE}
```

Then, access to [http://localhost:9082](http://localhost:9082).

On that page, you will see the list of registered projects and a link to register new projects.
Registering a new project requires only a unique ID string and an optional description text.

Once a new project has been registered, a static admin (username, password) will be automatically generated for the project admin. You can send that information to the project admin. The project admin first uses the provided static admin information to log in to PipeCD. After that, they can change the static admin information, configure the SSO, RBAC or disable static admin user.

__Caution:__ The Role-Based Access Control (RBAC) setting is required to enable your team login using SSO, please make sure you have that setup before disable static admin user.