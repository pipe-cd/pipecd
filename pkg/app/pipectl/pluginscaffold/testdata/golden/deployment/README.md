# demo PipeCD plugin (deployment)

Scaffolded with `pipectl plugin init --kind deployment`.

## Stages
- `DEMO_SYNC`
- `DEMO_ROLLBACK`

## Build

```bash
go build .
```

## Next steps

Implement stage logic under `deployment/` and configure deploy targets in `config/`.
