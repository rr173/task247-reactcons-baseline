package service

import (
	"context"
	"path/filepath"
	"testing"

	"task247-reactcons/internal/model"
	"task247-reactcons/internal/store"
)

func TestTask247Bug06ConstraintPersistence(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "reactcons.db")
	repository, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	svc := New(repository)
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "constraint-persistence"})
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range []string{"H2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, symbol, ""); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.AddReaction(ctx, net.ID, "H2 -> H2O", false); err != nil {
		t.Fatal(err)
	}
	constraint, err := svc.AddConstraint(ctx, net.ID, model.MoietyClosure{Members: []model.MoietyMember{{Symbol: "H2", Coeff: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	result, err := svc.Solve(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != model.NetworkConflict || result.ConstraintChecks[0].Status != "violated" {
		t.Fatalf("unexpected solve result: %+v", result)
	}
	if err := repository.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	constraints, err := New(reopened).ListConstraints(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(constraints) != 1 || constraints[0].ID != constraint.ID || constraints[0].Status != "violated" {
		t.Fatalf("constraint status was not recovered: %+v", constraints)
	}
}
