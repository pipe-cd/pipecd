# pipe

- Configuration Definition
    - Application, Piped, Control Plane, Notification, Metrics Template
        - https://github.com/kapetaniosci/pipe/blob/master/pkg/config
    - Example
        - https://github.com/kapetaniosci/pipe/tree/master/pkg/config/testdata

- Model Definition
    - https://github.com/kapetaniosci/pipe/tree/master/pkg/model

- Piped Component
    - https://github.com/kapetaniosci/pipe/tree/master/pkg/app/piped

- Control Plane Component
    - https://github.com/kapetaniosci/pipe/tree/master/pkg/app/api
    - https://github.com/kapetaniosci/pipe/tree/master/pkg/app/web


**DRAFT**

Powerful, Easy to Use, Easy to Operate

**Powerful**
- Unifed tool for various deployment kinds: kubernetes (plain-yaml, helm, kustomize), terraform, lambda, cloudrun...
- Deployment strategies: canary, bluegreen, rolling update
- Analysis by metrics, log, smoke test
- Automatic rollback
- Configuration drift detection
- Insight shows delivery perfomance

**Easy to Use**
- Operations by Pull Request: scale, rolling update, rollback by PR
- Realtime visualization of application state
- Deployment pipeline to see what is happenning
- Intuitive UI

**Easy to Operate**
- Just 2 components: piped and control-plane
- Piped can be run on kubernetes, vm or even local machine
- Easy to operate multi-tenancy, multi-cluster
- Security
