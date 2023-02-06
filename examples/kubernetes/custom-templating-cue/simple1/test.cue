package kube

deployment: simple1: {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name: "simple1"
		labels: app: "simple1"
	}
	spec: {
		replicas: 2
		selector: matchLabels: {
			app:                  "simple1"
			"pipecd.dev/variant": "primary"
		}
		template: {
			metadata: {
				labels: {
					app:                  "simple1"
					"pipecd.dev/variant": "primary"
				}
				annotations: "sidecar.istio.io/inject": "false"
			}
			spec: containers: [{
				name:  "helloworld"
				image: "ghcr.io/pipe-cd/helloworld:v0.32.0"
				args: [
					"server",
				]
				ports: [{
					containerPort: 9085
				}]
			}]
		}
	}
}
service: simple1: {
	apiVersion: "v1"
	kind:       "Service"
	metadata: name: "simple1"
	spec: {
		selector: app: "simple1"
		ports: [{
			protocol:   "TCP"
			port:       9085
			targetPort: 9085
		}]
	}
}
