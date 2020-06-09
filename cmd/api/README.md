# control-plane

## Development

### How to run a control-plane at your local environment

Prepare a configuration file in anywhere. The following is a sample configuration file.

``` yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  datastore:
    type: FIRESTORE
    config:
      namespace: sandbox
      project: pipecd
      credentialsFile: "/your-path-to-path/firestore-service-account-credential.json"
  filestore:
    type: GCS
    config:
      bucket: stage-logs-sandbox 
      credentialsFile: "/your-path-to-path/gcs-service-account-credential.json"
  cache:
    redisAddress: "localhost:6379"
    ttl: 5m
```

You can run control plane in local machine as follows:

``` console
bazelisk run //cmd/api:api -- server \
  --config-file=/your-path-to-path/control-plane.yaml
```

### How to run a control-plane with mock response mode

If you use web mock response, please write the following config.

``` yaml
apiVersion: "pipecd.dev/v1beta1"
kind: ControlPlane
spec:
  datastore: {}
  filestore: {}
  cache: {}
```

You can run mock control plane in local machine as follows:

``` console
bazelisk run //cmd/api:api -- server \
  --config-file=/your-path-to-path/control-plane-mock.yaml \
  --use-fake-response=true \
  --enable-grpc-reflection=true
```