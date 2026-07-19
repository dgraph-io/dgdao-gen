// Package movies_test provides end-to-end tests that verify the generated
// entity types (movies.*) integrate correctly with dgdao.AsRecord.
// These tests live here rather than in the root record_test.go because the
// generated entity packages import dgdao, which would cause an import
// cycle from the root test package. Placing the tests inside the testdata
// tree satisfies Go's internal-package visibility rule.
package movies_test

import (
	"testing"

	dgdao "github.com/dgraph-io/dgdao"
	movies "github.com/dgraph-io/dgdao-gen/cmd/dgdao-gen/internal/parser/testdata/movies"
	moviesSchema "github.com/dgraph-io/dgdao-gen/cmd/dgdao-gen/internal/parser/testdata/movies/schema"
)

// TestAsRecord_RealEntityRoutesToRecord verifies the AsRecord reflection
// probe correctly substitutes a generated entity (movies.Studio) with its
// backing record struct (*moviesSchema.Studio).
func TestAsRecord_RealEntityRoutesToRecord(t *testing.T) {
	s := &moviesSchema.Studio{Name: "Pixar"}
	w := movies.NewStudioWithRecord(s)

	out := dgdao.AsRecord(w)
	got, ok := out.(*moviesSchema.Studio)
	if !ok {
		t.Fatalf("expected *moviesSchema.Studio after AsRecord, got %T", out)
	}
	if got != s {
		t.Fatalf("expected AsRecord to return the SAME backing pointer; got a different *moviesSchema.Studio")
	}
}

// TestAsRecord_RealRecordPassthrough verifies that a plain
// *moviesSchema.Studio (already implementing the Record interface via its
// generated RecordTypeName method) passes through AsRecord unchanged.
func TestAsRecord_RealRecordPassthrough(t *testing.T) {
	s := &moviesSchema.Studio{Name: "Pixar"}
	out := dgdao.AsRecord(s)
	if out != any(s) {
		t.Fatalf("expected record struct to pass through unchanged; got %T", out)
	}
}

// TestRecordInterface_RealRecordSatisfies verifies that the generated
// record struct satisfies the dgdao.Record interface via its
// generated RecordTypeName method.
func TestRecordInterface_RealRecordSatisfies(t *testing.T) {
	var _ dgdao.Record = (*moviesSchema.Studio)(nil)
}

// TestRecordTypeName_RealRecordReturnsCanonical verifies the generated
// RecordTypeName returns the canonical entity name.
func TestRecordTypeName_RealRecordReturnsCanonical(t *testing.T) {
	s := &moviesSchema.Studio{}
	if got := s.RecordTypeName(); got != "Studio" {
		t.Fatalf("expected RecordTypeName() == %q, got %q", "Studio", got)
	}
}
