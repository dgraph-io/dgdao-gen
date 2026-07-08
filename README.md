# dgdao-gen

Code generation and the wrapper-entity runtime for
[dgdao](https://github.com/dgraph-io/dgdao). Define Go structs, run
`go generate`, and get a fully typed client, query builders, auto-paging iterators, and
a CLI — all derived from your struct definitions.

`dgdao-gen` was extracted from dgdao. The generic typed client and
query primitives stay in dgdao; this project owns the generator and the
`entity` wrapper base that generated code embeds.

## Install

```
go get github.com/dgraph-io/dgdao-gen
```

## Usage

Add a `go:generate` directive next to your schema structs:

```go
//go:generate go run github.com/dgraph-io/dgdao-gen/cmd/dgdao-gen -entities
```

then run `go generate ./...`.

Generated code imports the generic primitives from `dgdao/typed` and the wrapper
base from `dgdao-gen/wrap`:

```go
import (
    "github.com/dgraph-io/dgdao/typed"
    "github.com/dgraph-io/dgdao-gen/wrap"
)
```

<!-- Struct-tag reference and full CLI-flag table land before the first release. -->

## Related projects

Part of the [dgdao](https://github.com/dgraph-io/dgdao) family:

- [dgdao](https://github.com/dgraph-io/dgdao) — core typed client and DAO library for Dgraph
- [dgdao-migrate](https://github.com/dgraph-io/dgdao-migrate) — struct-first schema migration framework for Dgraph
- [dgdao-telemetry](https://github.com/dgraph-io/dgdao-telemetry) — OpenTelemetry instrumentation for the dgdao typed client

## License

Apache-2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
