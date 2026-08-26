package service

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
)

func TestTask247Bug07ConflictRefresh(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "conflict-refresh"})
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
	if _, err := svc.Solve(ctx, net.ID); err != nil {
		t.Fatal(err)
	}
	conflicts, err := svc.store.ListConflictSets(ctx, net.ID)
	if err != nil || len(conflicts) == 0 {
		t.Fatalf("initial conflict was not recorded: %v %+v", err, conflicts)
	}
	if err := svc.ExemptReaction(ctx, rx.ID, "accepted open-system exception"); err != nil {
		t.Fatal(err)
	}
	result, err := svc.Solve(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != model.NetworkPublishable {
		t.Fatalf("repaired network was not publishable: %s", result.Status)
	}
	conflicts, err = svc.store.ListConflictSets(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("stale conflict sets remained after re-solve: %+v", conflicts)
	}
}
