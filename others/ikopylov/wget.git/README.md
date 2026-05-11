# wget (scaffold)

This repository is a Go project that will grow into a `wget`-like CLI. For now, the entrypoint is intentionally a **no-op** (it exits 0 and prints nothing).

## Runbook

### Requirements

- Go **1.21+** (see `go.mod`)
- Read the project brief: [`docs/requirements.md`](docs/requirements.md)
- See the verification checklist: [`docs/audit.md`](docs/audit.md)

### Install dependencies

This scaffold has **no external dependencies** (standard library only).

### Build `./wget`

- Go **1.21+**

```bash
go build -o wget .
```

### Run `./wget`

```bash
./wget
```

### Run without building (dev)

```bash
go run .
```

### Test

```bash
go test ./...
```

### Lint / staticcheck

No lint/staticcheck step is required yet for this scaffold. For now, `go test ./...` is the only supported quality gate.

