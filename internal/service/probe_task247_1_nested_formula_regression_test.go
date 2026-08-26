package service

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
)

func TestTask247Bug01NestedFormula(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "nested-formula"})
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range []string{"Mg(Al(OH)4)2", "MgO", "Al2O3", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, symbol, ""); err != nil {
			t.Fatal(err)
		}
	}
	sp, err := svc.ListSpecies(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range sp {
		if item.Symbol == "Mg(Al(OH)4)2" && (item.Composition["Mg"] != 1 || item.Composition["Al"] != 2 || item.Composition["O"] != 8 || item.Composition["H"] != 8) {
			t.Fatalf("nested formula composition was not expanded: %+v", item.Composition)
		}
	}
	if _, err := svc.AddReaction(ctx, net.ID, "Mg(Al(OH)4)2 -> MgO + Al2O3 + 4 H2O", false); err != nil {
		t.Fatal(err)
	}
	result, err := svc.Solve(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != model.NetworkPublishable {
		t.Fatalf("balanced nested-formula reaction was not publishable: %s", result.Status)
	}
}
