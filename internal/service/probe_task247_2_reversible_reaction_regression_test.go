package service

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
)

func TestTask247Bug02ReversibleReaction(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "reversible"})
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range []string{"H2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, symbol, ""); err != nil {
			t.Fatal(err)
		}
	}
	created, err := svc.AddReaction(ctx, net.ID, "H2 <=> H2O", false)
	if err != nil {
		t.Fatal(err)
	}
	if !created.Reversible {
		t.Fatalf("reversible equation lost its directionality: %+v", created)
	}
	stored, err := svc.ListReactions(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(stored) != 1 || !stored[0].Reversible {
		t.Fatalf("reversible state was not persisted: %+v", stored)
	}
}
