---
title: "Building Plugins in Other Languages"
linkTitle: "Other Languages (Python/Node.js)"
weight: 23
description: >
  Demonstrating the language-agnostic nature of the PipeCD gRPC plugin architecture by creating minimal examples in Python and Node.js.
---

One of the key design achievements of **PipeCD v1** is its **language-agnostic plugin architecture**.

While the official SDK and the core plugin development book focus on Go, you can write a PipeCD plugin in **any language** that supports gRPC. The PipeCD Agent (`piped`) communicates with plugins entirely over gRPC using standard Protocol Buffers definitions.

This guide provides minimal, functional "Hello World" plugin snippets in **Python** and **Node.js** to show how easy it is to extend PipeCD beyond the Go ecosystem.

---

## The gRPC Interface

A PipeCD deployment plugin implements the `DeploymentService` defined in [api.proto](https://github.com/pipe-cd/pipecd/blob/master/pkg/plugin/api/v1alpha1/deployment/api.proto).

To qualify as a valid deployment plugin, your gRPC server must implement the following six RPCs:

| RPC Method                | Request Type                     | Response Type                     | Description                                                         |
| ------------------------- | -------------------------------- | --------------------------------- | ------------------------------------------------------------------- |
| `FetchDefinedStages`      | `FetchDefinedStagesRequest`      | `FetchDefinedStagesResponse`      | Declares which custom stage names this plugin supports.             |
| `DetermineVersions`       | `DetermineVersionsRequest`       | `DetermineVersionsResponse`       | Analyzes deployment inputs to determine the artifact versions used. |
| `DetermineStrategy`       | `DetermineStrategyRequest`       | `DetermineStrategyResponse`       | Decides the sync strategy (e.g., pipeline sync or quick sync).      |
| `BuildPipelineSyncStages` | `BuildPipelineSyncStagesRequest` | `BuildPipelineSyncStagesResponse` | Builds the sequence of stages for a pipeline execution.             |
| `BuildQuickSyncStages`    | `BuildQuickSyncStagesRequest`    | `BuildQuickSyncStagesResponse`    | Builds the sequence of stages for a quick sync execution.           |
| `ExecuteStage`            | `ExecuteStageRequest`            | `ExecuteStageResponse`            | Executes a specific stage during deployment and returns the status. |

---

## 1. Python "Hello World" Plugin Snippet

Python is highly popular for infrastructure scripting, custom automated analyses, and developer tooling. Below is a minimal gRPC server implementing the PipeCD `DeploymentService` using the standard `grpcio` library.

### Installation

Ensure you have the required gRPC dependencies installed:

```bash
pip install grpcio grpcio-tools
```

### Python Implementation (`plugin.py`)

```python
import time
from concurrent import futures
import grpc

# Import your generated proto definitions here.
# Typically generated from:
# - pkg/plugin/api/v1alpha1/deployment/api.proto
# - pkg/model/deployment.proto
# - pkg/model/common.proto
# - pkg/plugin/api/v1alpha1/common/common.proto
import deployment_pb2
import deployment_pb2_grpc

class HelloWorldPlugin(deployment_pb2_grpc.DeploymentServiceServicer):
    def FetchDefinedStages(self, request, context):
        """Declare that this plugin supports a custom stage named 'HELLO_WORLD'."""
        print("Received FetchDefinedStages request")
        return deployment_pb2.FetchDefinedStagesResponse(
            stages=["HELLO_WORLD"]
        )

    def DetermineVersions(self, request, context):
        """We do not require version detection in this minimal example."""
        print("Received DetermineVersions request")
        return deployment_pb2.DetermineVersionsResponse(versions=[])

    def DetermineStrategy(self, request, context):
        """We tell Piped to bypass strategy determination."""
        print("Received DetermineStrategy request")
        return deployment_pb2.DetermineStrategyResponse(
            unsupported=True,
            summary="HelloWorld plugin does not determine strategy"
        )

    def BuildPipelineSyncStages(self, request, context):
        """Accept the input stages as-is and wrap them in standard PipelineStage formats."""
        print("Received BuildPipelineSyncStages request")
        response_stages = []
        for stage_config in request.stages:
            response_stages.append(
                deployment_pb2.PipelineStage(
                    id=f"stage-{stage_config.index}",
                    name=stage_config.name,
                    desc=stage_config.desc or "Hello World Custom Sync",
                    index=stage_config.index,
                    status=0,  # STAGE_NOT_STARTED_YET
                    rollback=request.rollback
                )
            )
        return deployment_pb2.BuildPipelineSyncStagesResponse(stages=response_stages)

    def BuildQuickSyncStages(self, request, context):
        """Provide a default single-stage sync pipeline for quick syncs."""
        print("Received BuildQuickSyncStages request")
        response_stages = [
            deployment_pb2.PipelineStage(
                id="stage-quick-0",
                name="HELLO_WORLD",
                desc="Hello World Quick Sync",
                index=0,
                status=0,  # STAGE_NOT_STARTED_YET
                rollback=request.rollback
            )
        ]
        return deployment_pb2.BuildQuickSyncStagesResponse(stages=response_stages)

    def ExecuteStage(self, request, context):
        """Execute the HELLO_WORLD stage, print the greeting, and return SUCCESS."""
        stage_name = request.input.stage.name
        print(f"Executing stage: {stage_name}")

        if stage_name == "HELLO_WORLD":
            print("====================================")
            print("   👋 Hello, World from Python!     ")
            print("====================================")
            # StageStatus: STAGE_SUCCESS = 2
            return deployment_pb2.ExecuteStageResponse(
                status=2,
                message="Hello World stage executed successfully!"
            )

        # StageStatus: STAGE_FAILURE = 3
        return deployment_pb2.ExecuteStageResponse(
            status=3,
            message=f"Unsupported stage: {stage_name}"
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    deployment_pb2_grpc.add_DeploymentServiceServicer_to_server(
        HelloWorldPlugin(), server
    )
    # Listen on port 50051 (or a custom port specified in piped.yaml)
    server.add_insecure_port("[::]:50051")
    print("Python HelloWorld Plugin Server is running on port 50051...")
    server.start()
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()
```

---

## 2. Node.js "Hello World" Plugin Snippet

Node.js offers excellent async processing and dynamic loading. Below is a minimal gRPC server implementing the exact same PipeCD plugin service using `@grpc/grpc-js` and `@grpc/proto-loader`.

### Installation

Ensure you have the required gRPC and proto-loader packages installed in your project:

```bash
npm install @grpc/grpc-js @grpc/proto-loader
```

### Node.js Implementation (`plugin.js`)

```javascript
const path = require("path");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");

// Path to the service.proto or unified proto definition
const PROTO_PATH = path.join(__dirname, "deployment_api.proto");

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
  // Make sure to configure includeDirs if referencing external imports
  includeDirs: [__dirname],
});

const pluginProto =
  grpc.loadPackageDefinition(packageDefinition).grpc.plugin.deploymentapi
    .v1alpha1;

/**
 * Implements the DeploymentService gRPC service.
 */
const helloWorldPlugin = {
  fetchDefinedStages: (call, callback) => {
    console.log("Received fetchDefinedStages request");
    callback(null, {
      stages: ["HELLO_WORLD"],
    });
  },

  determineVersions: (call, callback) => {
    console.log("Received determineVersions request");
    callback(null, { versions: [] });
  },

  determineStrategy: (call, callback) => {
    console.log("Received determineStrategy request");
    callback(null, {
      unsupported: true,
      summary: "HelloWorld plugin does not determine strategy",
    });
  },

  buildPipelineSyncStages: (call, callback) => {
    console.log("Received buildPipelineSyncStages request");
    const { stages, rollback } = call.request;
    const responseStages = stages.map((stage) => ({
      id: `stage-${stage.index}`,
      name: stage.name,
      desc: stage.desc || "Hello World Custom Sync",
      index: stage.index,
      status: "STAGE_NOT_STARTED_YET",
      rollback: rollback,
    }));
    callback(null, { stages: responseStages });
  },

  buildQuickSyncStages: (call, callback) => {
    console.log("Received buildQuickSyncStages request");
    const { rollback } = call.request;
    const responseStages = [
      {
        id: "stage-quick-0",
        name: "HELLO_WORLD",
        desc: "Hello World Quick Sync",
        index: 0,
        status: "STAGE_NOT_STARTED_YET",
        rollback: rollback,
      },
    ];
    callback(null, { stages: responseStages });
  },

  executeStage: (call, callback) => {
    const stageName = call.request.input.stage.name;
    console.log(`Executing stage: ${stageName}`);

    if (stageName === "HELLO_WORLD") {
      console.log("====================================");
      console.log("   👋 Hello, World from Node.js!    ");
      console.log("====================================");
      // StageStatus: STAGE_SUCCESS
      callback(null, {
        status: "STAGE_SUCCESS",
        message: "Hello World stage executed successfully!",
      });
    } else {
      // StageStatus: STAGE_FAILURE
      callback(null, {
        status: "STAGE_FAILURE",
        message: `Unsupported stage: ${stageName}`,
      });
    }
  },
};

/**
 * Starts the gRPC Server.
 */
function main() {
  const server = new grpc.Server();
  server.addService(pluginProto.DeploymentService.service, helloWorldPlugin);

  // Listen on port 50051 (or custom port configured in piped.yaml)
  server.bindAsync(
    "0.0.0.0:50051",
    grpc.ServerCredentials.createInsecure(),
    (err, port) => {
      if (err) {
        console.error(`Failed to start server: ${err}`);
        return;
      }
      console.log(
        `Node.js HelloWorld Plugin Server is running on port ${port}...`,
      );
    },
  );
}

main();
```

---

## 3. Configuring Piped to Use Your Custom Plugin

To let PipeCD's `piped` daemon use your custom Python or Node.js plugin, you need to register it in `piped.yaml`.

```yaml
apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  # ... other configurations
  plugins:
    - name: helloworld-plugin
      # Since we ran our custom server on port 50051, we point piped directly to it.
      # You can also run it via a binary using the standard plugin execution model.
      address: localhost:50051
```

Once registered, you can start writing application configurations (`.pipe.yaml`) that leverage the custom stage:

```yaml
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: example-app
  labels:
    platform: kubernetes
  pipeline:
    stages:
      # Invoke your custom HELLO_WORLD stage running inside your custom gRPC server
      - name: HELLO_WORLD
        plugin: helloworld-plugin
        desc: "Execute custom hello world logic"
```
