{
	name: 'simple',
	image: error 'image is required',
	labels: {
	  'app': $.name
	},
	port: 9085,
	deployment: {
	  apiVersion: 'apps/v1',
	  kind: 'Deployment',
	  metadata: {
		name: $.name,
		labels: $.labels,
	  },
	  spec: {
		selector: {
		  matchLabels: $.labels
		},
		replicas: 1,
		template: {
		  metadata: {
			name: $.name,
			labels: $.labels,
		  },
		  spec: {
			containers: [
			  {
				name: 'helloworld',
				image: $.image,
				args: [
                   'server'
				],
				ports: [
				  {
					containerPort: $.port,
				  }
				],
			  }
			]
		  }
		}
	  }
	},
	service: {
	  apiVersion: 'v1',
	  kind: 'Service',
	  metadata: {
		name: $.name,
	  },
	  spec: {
		ports: [
		  {
		    protocol: 'TCP',
			port: 9085,
			targetPort: 9085
		  }
		],
		selector: $.labels
	  }
	},
	all: [$.deployment, $.service]
}
