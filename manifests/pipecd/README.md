## Maintainer guide
### How to update dependencies
1. To run Helm, replace the `version` and `appVersion` fields in Chart.yaml with temporary versions (whatever is fine)
2. Update the `dependencies[*].version` field in Chart.yaml with what you need
3. Update Chart.lock with:

```
helm dependency update manifests/pipecd
```

4. Remove charts under the `charts` directory.
