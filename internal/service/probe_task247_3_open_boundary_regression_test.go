package service

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
)

func TestTask247Bug03OpenBoundary(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "open-boundary"})
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range []string{"H2", "O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, symbol, ""); err != nil {
			t.Fatal(err)
		}
	}
	species, err := svc.ListSpecies(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	var oxygen int64
	for _, item := range species {
		if item.Symbol == "O" {
			oxygen = item.ID
		}
	}
	if oxygen == 0 {
		t.Fatal("oxygen species was not created")
	}
	if _, err := svc.AddReaction(ctx, net.ID, "H2 -> H2 + O", false); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.MarkBoundary(ctx, net.ID, oxygen, "oxygen exchange"); err != nil {
		t.Fatal(err)
	}
	result, err := svc.Solve(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != model.NetworkPublishable {
		t.Fatalf("open-boundary reaction was treated as a conflict: %s", result.Status)
	}
}
