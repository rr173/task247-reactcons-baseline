package service

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
)

func TestTask247Bug04ExemptReaction(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "exempt-reaction"})
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range []string{"H2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, symbol, ""); err != nil {
			t.Fatal(err)
		}
	}
	rx, err := svc.AddReaction(ctx, net.ID, "H2 -> H2O", false)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.ExemptReaction(ctx, rx.ID, "open system step"); err != nil {
		t.Fatal(err)
	}
	result, err := svc.Solve(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != model.NetworkPublishable || len(result.ConflictSets) != 0 {
		t.Fatalf("exempt reaction still affected diagnosis: status=%s conflicts=%+v", result.Status, result.ConflictSets)
	}
}
