package service

import (
	"context"
	"math"
	"math/big"
	"testing"

	"task247-reactcons/internal/model"
)

func TestTask247Bug05NullspaceConservation(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "nullspace"})
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range []string{"H2", "O2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, symbol, ""); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.AddReaction(ctx, net.ID, "H2 + O2 -> H2O", false); err != nil {
		t.Fatal(err)
	}
	pools, err := svc.ConservedPools(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(pools) < 2 {
		t.Fatalf("expected non-trivial linear nullspace pools, got %+v", pools)
	}
	reactions, err := svc.ListReactions(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, pool := range pools {
		netValue := 0.0
		for _, member := range pool.Members {
			coeff, ok := new(big.Rat).SetString(member.Coeff)
			if !ok {
				t.Fatalf("invalid rational pool coefficient %q", member.Coeff)
			}
			value, _ := coeff.Float64()
			for _, participant := range reactions[0].Participants {
				if participant.SpeciesID != member.SpeciesID {
					continue
				}
				if participant.Role == "product" {
					netValue += value * participant.Coeff
				} else {
					netValue -= value * participant.Coeff
				}
			}
		}
		if math.Abs(netValue) > 1e-9 {
			t.Fatalf("pool %s is not in the reaction nullspace: net=%g members=%+v", pool.Label, netValue, pool.Members)
		}
	}
}
