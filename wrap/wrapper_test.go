/*
 * SPDX-FileCopyrightText: © 2017-2026 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package wrap_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matthewmcneely/modusgraph"
	"github.com/mlwelles/modusgraph-gen/wrap"
)

// widget is a minimal schema struct used to exercise the entity wrapper.
type widget struct {
	UID   string   `json:"uid,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
	Name  string   `json:"name,omitempty" dgraph:"index=exact"`
	Qty   int      `json:"qty,omitempty" dgraph:"index=int"`
}

func TestWrapper_UnwrapReturnsBackingPointer(t *testing.T) {
	backing := &widget{Name: "inner"}
	w := wrap.WrapValue(backing)
	if w.Unwrap() != backing {
		t.Fatal("Unwrap did not return the same backing pointer")
	}
}

func TestWrapper_JSONRoundTrip(t *testing.T) {
	w := wrap.WrapValue(&widget{Name: "j", Qty: 5})
	data, err := json.Marshal(&w)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got wrap.Wrapper[widget]
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Unwrap().Name != "j" || got.Unwrap().Qty != 5 {
		t.Fatalf("round trip lost data: %+v", got.Unwrap())
	}
}

func TestWrapper_UnmarshalIntoZeroValue(t *testing.T) {
	// UnmarshalJSON must lazily allocate the backing struct.
	var w wrap.Wrapper[widget]
	if err := json.Unmarshal([]byte(`{"name":"z"}`), &w); err != nil {
		t.Fatalf("Unmarshal into zero value: %v", err)
	}
	if w.Unwrap() == nil || w.Unwrap().Name != "z" {
		t.Fatalf("Unmarshal did not allocate/populate backing struct: %+v", w.Unwrap())
	}
}

func TestWrapper_Validate(t *testing.T) {
	w := wrap.WrapValue(&widget{Name: "v"})
	if err := w.Validate(context.Background(), modusgraph.NewValidator()); err != nil {
		t.Fatalf("Validate on a tag-free struct should pass; got %v", err)
	}
}

// embedTestEntity mirrors how generated entity types embed Wrapper as a value
// field. The compile-time guards below pin the contract Wrapper exists for: a
// *Entity must satisfy json.Marshaler and json.Unmarshaler via promotion of
// Wrapper's pointer-receiver methods.
type embedTestEntity struct {
	wrap.Wrapper[widget]
}

var (
	_ json.Marshaler   = (*embedTestEntity)(nil)
	_ json.Unmarshaler = (*embedTestEntity)(nil)
)
