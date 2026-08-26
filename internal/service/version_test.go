package service

import (
	"context"
	"testing"

	"task247-reactcons/internal/model"
	"task247-reactcons/internal/store"
)

func TestPublishSameContentIsIdempotent(t *testing.T) {
	repository, err := store.Open(t.TempDir() + "/reactcons.db")
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()

	svc := New(repository)
	ctx := context.Background()
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "idempotent-publish"})
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
	first, err := svc.PublishVersion(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.PublishVersion(ctx, net.ID)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID || first.ContentHash != second.ContentHash || second.Status != model.VersionPublished {
		t.Fatalf("same content was not idempotent: first=%+v second=%+v", first, second)
	}
}
