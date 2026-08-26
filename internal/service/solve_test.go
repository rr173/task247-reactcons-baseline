package service

import (
	"context"
	"path/filepath"
	"testing"

	"task247-reactcons/internal/model"
	"task247-reactcons/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	dir := t.TempDir()
	s, err := store.Open(filepath.Join(dir, "reactcons_test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestSolvePublishable(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "t"})
	if err != nil {
		t.Fatal(err)
	}
	for _, sym := range []string{"H2", "O2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, sym, ""); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.AddReaction(ctx, net.ID, "2 H2 + O2 -> 2 H2O", false); err != nil {
		t.Fatal(err)
	}
	res, err := svc.Solve(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != model.NetworkPublishable {
		t.Fatalf("expected publishable, got %s", res.Status)
	}
	if len(res.ConservedPools) < 1 {
		t.Fatalf("expected conserved pools")
	}
	ver, err := svc.PublishVersion(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ver.Status != model.VersionPublished || ver.ContentHash == "" {
		t.Fatalf("version not published: %+v", ver)
	}
}

func TestSolveConflictAndConstraintCheck(t *testing.T) {
	ctx := context.Background()
	svc := New(tempStore(t))
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "t2"})
	if err != nil {
		t.Fatal(err)
	}
	for _, sym := range []string{"H2", "O2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, sym, ""); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.AddReaction(ctx, net.ID, "2 H2 + O2 -> H2", false); err != nil {
		t.Fatal(err)
	}
	res, err := svc.Solve(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != model.NetworkConflict {
		t.Fatalf("expected conflict, got %s", res.Status)
	}
	if len(res.ConflictSets) == 0 {
		t.Fatal("expected conflict sets")
	}
	if _, err := svc.PublishVersion(ctx, net.ID); err == nil {
		t.Fatal("publishing a conflicting network must fail")
	}
}
