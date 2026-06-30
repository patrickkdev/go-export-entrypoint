Extract the Go files required by an entrypoint into another directory.

## Install

```bash
go install github.com/patrickkdev/go-export-entrypoint@latest
```

or run directly:

```bash
go run github.com/patrickkdev/go-export-entrypoint@latest ./cmd/worker ../worker
```

## Usage

```bash
go-export-entrypoint <entrypoint> <destination>
```

Example:

```bash
go-export-entrypoint ./cmd/worker ../worker
```

The tool:

- copies `go.mod`
- copies `go.sum` (if present)
- finds every package needed to build the entrypoint
- copies every package that belongs to the current module
- copies only `.go` files while preserving the package layout

## Notes

- External modules are not copied.
- Packages are copied as whole directories (only `.go` files are included).
