package store

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
)

func TestConservedPoolRoundTripKeepsMemberSymbols(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	net, err := s.CreateNetwork(ctx, model.Network{Name: "pool-roundtrip"})
	if err != nil {
		t.Fatal(err)
	}
	want := []model.ConservedPool{{Label: "element:H", Members: []model.PoolMember{{SpeciesID: 11, Symbol: "H2", Coeff: "2"}}}}
	if err := s.UpsertConservedPools(ctx, net.ID, want); err != nil {
		t.Fatal(err)
	}
	got, err := s.ListConservedPools(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].Members) != 1 || got[0].Members[0].Symbol != "H2" {
		t.Fatalf("pool symbols were not preserved: %+v", got)
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir() + "/reactcons.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}
