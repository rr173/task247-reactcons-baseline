package service

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
	"task247-reactcons/internal/store"
)

func TestTask247Bug10RestartRecovery(t *testing.T) {
	ctx := context.Background()
	dbPath := t.TempDir() + "/reactcons.db"
	repository, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	svc := New(repository)
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "restart-recovery"})
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range []string{"H2", "O2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, symbol, ""); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.AddReaction(ctx, net.ID, "2 H2 + O2 -> 2 H2O", false); err != nil {
		t.Fatal(err)
	}
	if err := repository.UpdateNetworkStatus(ctx, net.ID, model.NetworkSolving); err != nil {
		t.Fatal(err)
	}
	if err := repository.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if err := New(reopened).Recover(ctx); err != nil {
		t.Fatal(err)
	}
	recovered, err := New(reopened).GetNetwork(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if recovered.Status != model.NetworkPublishable {
		t.Fatalf("solving network was not recovered: %+v", recovered)
	}
}
