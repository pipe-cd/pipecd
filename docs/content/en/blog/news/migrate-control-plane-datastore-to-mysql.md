---
date: 2021-04-01
title: "Migrate Control Plane datastore to MySQL"
linkTitle: "Migrate Control Plane datastore from MongoDB to MySQL"
weight: 999
description: "This page describes how to migrate Control Plane' datastore from MongoDB to MySQL."
author: Khanh Tran ([@khanhtc1202](https://twitter.com/khanhtc1202))
toc_hide: true
---

Since PipeCD release [v0.9.8](/blog/2021/03/25/release-v0.9.8) which introduces MySQL as PipeCD control-plane datastore, we plan to drop the support for MongoDB datastore in the near future.
Consider the supports of cloud providers (GCP, AWS, Azure, etc), MySQL has a higher status over MongoDB is one of the most reasons for this support-dropping conclusion.

If you are using MongoDB as your PipeCD control-plane datastore, the following guide is manipulation to migrate your datastore to MySQL.

### Before start

You need to install the [pipectl](/docs/user-guide/command-line-tool/#installation) version [v0.9.10-1-4ab28c0](https://github.com/pipe-cd/pipecd/releases/tag/v0.9.10-1-4ab28c0) command-line tool.

Validate your installed `pipectl` command

```console
$ pipectl datastore -h
Manage control-plane datastore resource.

Usage:
  pipectl datastore [command]

Available Commands:
  migrate     Migrate data to MySQL datastore.
```

### Migration

Step by step guide for migration

#### 1. Prepare your MySQL datastore instance.

Note:
- __MySQL v8.0 or later is required__.
- To enable `pipectl` to migrate data to the new MySQL datastore, your new MySQL data instance has to be connected from `pipectl` running environment.\
The current implementation of `pipectl datastore` subcommand connects directly to datastore, you do not need to authenticate your `pipectl` binary with `API Key` as for other subcommands. Just make sure `pipectl` running environment is in the same network with your MySQL instance is enough.
- Your MySQL instance has to be initialized with the going to be used `database`, make sure you create it before move to the next step.

For example, if you use docker to create your new MySQL instance, the command should be
```console
$ docker run -d \
    --name test-db \
    -e MYSQL_ROOT_PASSWORD=XXX \
    -e MYSQL_DATABASE=*database-name* \
    mysql:8.0
```

#### 2. Stop PipeCD control-plane

In case your PipeCD control-plane is installed by `helm`, simply run `helm uninstall` command would help.

#### 3. Migrate using `pipectl`

Migrate using the following command (replace `upstream-data-src`, `downstream-data-src` and `database` with your corresponding values)

```console
$ pipectl datastore migrate \
    --upstream-data-src="mongodb://127.0.0.1:27017/quickstart" \
    --downstream-data-src="root:test@tcp(127.0.0.1:3306)" \
    --database=quickstart
```

Note:
- Make sure your `data-src`s value are formatted as same as the above example.
- If you want to migrate only specific data models (not all at once), use the `--models` flag as follow `--models=Application,Project`. (Use `pipectl datastore migrate -h` to get the list of migratable models)

#### 4. Start PipeCD control-plane with the new configuration

Your new control-plane's configuration should be updated as follow:

```diff
 apiVersion: "pipecd.dev/v1beta1"
 kind: ControlPlane
 spec:
   datastore:
-    type: MONGODB
+    type: MYSQL
     config:
-      url: mongodb://127.0.0.1:27017/quickstart
+      url: root:test@tcp(127.0.0.1:3306)
       database: quickstart
...
```
See [ConfigurationReference](/docs/operator-manual/control-plane/configuration-reference/) for the full configuration.

Restart PipeCD control-plane as same as you start it before, your PipeCD should be ready 🚀.
