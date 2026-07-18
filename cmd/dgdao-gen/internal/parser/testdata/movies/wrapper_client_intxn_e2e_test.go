// Package movies_test provides end-to-end tests that exercise the generated
// top-level movies.Client's transaction delegators, NewTxnContext and InTxn,
// against a real, local file-backed dgdao client. The text-generation tests in
// the generator suite prove the two methods are emitted; these tests prove
// InTxn scopes a typed entity sub-client (movies.Client.Film) to the
// transaction — the property the wrapper_client.go.tmpl change exists to
// deliver: NewClient(conn) builds every per-entity sub-client purely from
// conn, so calling NewClient again on the txn-scoped conn scopes them all in
// one step. Like wrapper_query_e2e_test.go, this file lives inside the
// testdata tree because the generated package imports dgdao, which would
// cause an import cycle from the root test package, and reuses that file's
// newConn/addFilm helpers.
//
// These tests do not assert isolation semantics — that an uncommitted write
// staged inside the txn is invisible to a read outside it, or that Discard
// rolls a valid write back. dgdao's own suite notes that its local,
// file-backed embedded engine commits each successful mutation as it is
// staged rather than deferring to the client-side Commit call, so that
// guarantee is observable only against a real Dgraph cluster. What these
// tests prove deterministically on the embedded engine: staging a write
// through client.InTxn(tx).Film lands after tx.Commit(), and a query issued
// through that same scoped Film sub-client finds it — InTxn wires the
// generated entity sub-client to the transaction end to end.
//
// None of these tests call t.Parallel(): the dgdao engine is a strict
// process-wide singleton (only one client may exist at a time), so the tests
// must run sequentially. Each test gets its own t.TempDir()-backed client that
// t.Cleanup closes before the next test starts.
package movies_test

import (
	"context"
	"testing"

	movies "github.com/dgraph-io/dgdao-gen/cmd/dgdao-gen/internal/parser/testdata/movies"
)

// TestWrapperClient_InTxn_WritesStageAndCommit stages a Film write through
// client.InTxn(tx).Film, commits, and verifies the write landed via the
// original client. It proves Client.InTxn scopes a generated entity
// sub-client's Add to the transaction, not just the untyped surface.
func TestWrapperClient_InTxn_WritesStageAndCommit(t *testing.T) {
	ctx := context.Background()
	client := movies.NewClient(newConn(t))

	// A transaction's staged writes do not run autoSchema, so establish the
	// Film schema with a prior single-shot write before opening the txn.
	addFilm(ctx, t, client, "schema-seed", 2000)

	tx := client.NewTxnContext(ctx)
	defer tx.Discard()
	scoped := client.InTxn(tx)

	added := movies.NewFilm(movies.WithFilmName("staged"))
	if err := scoped.Film.Add(ctx, added); err != nil {
		t.Fatalf("scoped Film.Add: %v", err)
	}
	if added.UID() == "" {
		t.Fatal("scoped Film.Add should populate the UID")
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("tx.Commit: %v", err)
	}

	got, err := client.Film.Get(ctx, added.UID())
	if err != nil {
		t.Fatalf("Get after commit: %v", err)
	}
	if got == nil || got.Name() != "staged" {
		t.Fatalf("Get after commit = %+v, want a Film named %q", got, "staged")
	}
}

// TestWrapperClient_InTxn_QueryReadsThroughTxn proves that a query issued
// through the txn-scoped Film sub-client (client.InTxn(tx).Film.Query) finds
// a write staged through that same scoped client, before the txn is
// committed. It exercises the query path of the generated entity sub-client
// under InTxn, complementing the write-path coverage above.
func TestWrapperClient_InTxn_QueryReadsThroughTxn(t *testing.T) {
	ctx := context.Background()
	client := movies.NewClient(newConn(t))

	addFilm(ctx, t, client, "schema-seed", 2000)

	tx := client.NewTxnContext(ctx)
	defer tx.Discard()
	scoped := client.InTxn(tx)

	added := movies.NewFilm(movies.WithFilmName("findme"))
	if err := scoped.Film.Add(ctx, added); err != nil {
		t.Fatalf("scoped Film.Add: %v", err)
	}

	got, err := scoped.Film.Query(ctx).Filter(`eq(name, "findme")`).Nodes()
	if err != nil {
		t.Fatalf("scoped Query.Nodes: %v", err)
	}
	if len(got) != 1 || got[0].Name() != "findme" {
		t.Fatalf("scoped Query(name=findme) returned %d films, want exactly one named %q", len(got), "findme")
	}
}
