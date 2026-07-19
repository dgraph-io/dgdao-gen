// End-to-end coverage for the generated edge option-functions: constructing
// an entity with With<X><Edge>(...) must wire the edge exactly as a post-hoc
// Set<Edge> call would, so construction sites need no follow-up setter. Lives
// in the testdata tree for the same import-cycle reason as its siblings, and
// reuses wrapper_query_e2e_test.go's newConn/addFilm helpers.
package movies_test

import (
	"context"
	"testing"

	movies "github.com/dgraph-io/dgdao-gen/cmd/dgdao-gen/internal/parser/testdata/movies"
)

// TestEdgeOption_SetsEdgeAtConstruction inserts a Director, constructs a Film
// with the directors edge supplied via the generated WithFilmDirectors option,
// and proves the edge landed by querying films through WhereDirectors.
func TestEdgeOption_SetsEdgeAtConstruction(t *testing.T) {
	ctx := context.Background()
	client := movies.NewClient(newConn(t))

	// Establish the Film schema with a single-shot write.
	addFilm(ctx, t, client, "schema-seed", 2000)

	d := movies.NewDirector(movies.WithDirectorName("Edgewise"))
	if err := client.Director.Insert(ctx, d); err != nil {
		t.Fatalf("Director.Insert: %v", err)
	}

	f := movies.NewFilm(
		movies.WithFilmName("edge-option-film"),
		movies.WithFilmDirectors(d),
	)
	if err := client.Film.Insert(ctx, f); err != nil {
		t.Fatalf("Film.Insert: %v", err)
	}

	got, err := client.Film.Query(ctx).WhereDirectors(`eq(name, $1)`, "Edgewise").Nodes()
	if err != nil {
		t.Fatalf("WhereDirectors query: %v", err)
	}
	if len(got) != 1 || got[0].Name() != "edge-option-film" {
		t.Fatalf("WhereDirectors(Edgewise) returned %d films, want exactly one named %q", len(got), "edge-option-film")
	}
}
