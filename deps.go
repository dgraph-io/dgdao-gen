//go:build never

// Package dgdaogen carries no code; this file anchors module requirements.
package dgdaogen

import (
	// Anchor the dgdao requirement for the generated test fixtures under
	// cmd/dgdao-gen/internal/parser/testdata/, which the go tool's package
	// walker skips. The generator emits code that imports dgdao (and, via
	// dgdao's module graph, dgo and dgman); without a real import here,
	// go mod tidy would prune the modules those fixtures compile against.
	_ "github.com/dgraph-io/dgdao"
)
