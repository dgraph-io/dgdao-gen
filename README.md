# modusgraph-gen

Code generation and the wrapper-entity runtime for
[modusgraph](https://github.com/matthewmcneely/modusgraph). Define Go structs, run
`go generate`, and get a fully typed client, query builders, auto-paging iterators, and
a CLI — all derived from your struct definitions.

`modusgraph-gen` was extracted from a fork of modusGraph. The generic typed client and
query primitives stay in modusGraph; this project owns the generator and the
`entity` wrapper base that generated code embeds.

## Install

`modusgraph-gen` depends on a fork of modusGraph that is published under a different
import path. Go does not propagate `replace` directives to consumers, so **your project
must declare the same replace directive**:

```
// go.mod
require (
    github.com/mlwelles/modusgraph-gen v0.1.0
    github.com/matthewmcneely/modusgraph v0.0.0-00010101000000-000000000000
)

replace github.com/matthewmcneely/modusgraph => github.com/mlwelles/modusGraph v0.5.0-dev-mlwelles-20260604a
```

Without this directive, the build resolves `matthewmcneely/modusgraph` to upstream,
which lacks the typed client and bulk loader the generated code depends on.

## Usage

Add a `go:generate` directive next to your schema structs:

```go
//go:generate go run github.com/mlwelles/modusgraph-gen/cmd/modusgraph-gen -entities
```

then run `go generate ./...`.

Generated code imports the generic primitives from `modusgraph/typed` and the wrapper
base from `modusgraph-gen/entity` (aliased to avoid a package-name collision in
projects whose own package is `entity`):

```go
import (
    "github.com/matthewmcneely/modusgraph/typed"
    mgentity "github.com/mlwelles/modusgraph-gen/entity"
)
```

<!-- Struct-tag reference and full CLI-flag table land before the first release. -->

## License

Apache-2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
