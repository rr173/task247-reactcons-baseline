// Command reactcons exposes the chemical reaction network conservation diagnostic service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"task247-reactcons/internal/httpapi"
	"task247-reactcons/internal/model"
	"task247-reactcons/internal/service"
	"task247-reactcons/internal/store"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	dbPath := flag.String("db", "reactcons.db", "SQLite database path")
	smoke := flag.Bool("smoke-test", false, "run deterministic end-to-end self test")
	flag.Parse()
	if *smoke {
		if err := runSmoke(); err != nil {
			log.Fatalf("smoke test failed: %v", err)
		}
		fmt.Println("SMOKE TEST PASSED")
		return
	}
	repository, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	svc := service.New(repository)
	if err := svc.Recover(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Printf("reactcons listening on %s (db=%s)", *addr, *dbPath)
	if err := http.ListenAndServe(*addr, httpapi.New(svc).Handler()); err != nil {
		log.Fatal(err)
	}
}

func runSmoke() error {
	dir, err := os.MkdirTemp("", "reactcons-smoke-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, "reactcons.db")
	ctx := context.Background()
	repository, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	svc := service.New(repository)

	// Balanced network: combustion of hydrogen; must be publishable.
	net, err := svc.CreateNetwork(ctx, model.Network{Name: "smoke-balanced", Description: "balanced combustion"})
	if err != nil {
		return err
	}
	for _, sym := range []string{"H2", "O2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net.ID, sym, ""); err != nil {
			return err
		}
	}
	if _, err := svc.AddReaction(ctx, net.ID, "2 H2 + O2 -> 2 H2O", false); err != nil {
		return err
	}
	res, err := svc.Solve(ctx, net.ID)
	if err != nil {
		return err
	}
	if res.Status != model.NetworkPublishable {
		return fmt.Errorf("balanced network expected publishable, got %s", res.Status)
	}
	if len(res.ConservedPools) < 1 {
		return fmt.Errorf("expected at least one conserved pool")
	}
	ver, err := svc.PublishVersion(ctx, net.ID)
	if err != nil {
		return err
	}
	if ver.Status != model.VersionPublished || ver.ContentHash == "" {
		return fmt.Errorf("version not published")
	}
	verID := ver.ID

	// Conflict network: a hydrogen-oxygen reaction missing its product water; must report conflict.
	net2, err := svc.CreateNetwork(ctx, model.Network{Name: "smoke-conflict", Description: "missing product"})
	if err != nil {
		return err
	}
	for _, sym := range []string{"H2", "O2", "H2O"} {
		if _, err := svc.AddSpecies(ctx, net2.ID, sym, ""); err != nil {
			return err
		}
	}
	if _, err := svc.AddReaction(ctx, net2.ID, "2 H2 + O2 -> H2", false); err != nil {
		return err
	}
	res2, err := svc.Solve(ctx, net2.ID)
	if err != nil {
		return err
	}
	if res2.Status != model.NetworkConflict {
		return fmt.Errorf("conflict network expected conflict, got %s", res2.Status)
	}
	if len(res2.ConflictSets) == 0 {
		return fmt.Errorf("expected conflict sets for the unbalanced reaction")
	}
	if _, err := svc.PublishVersion(ctx, net2.ID); err == nil {
		return fmt.Errorf("publishing a conflicting network must fail")
	}

	// Restart recovery: close and reopen the database, confirm the frozen snapshot persists.
	if err := repository.Close(); err != nil {
		return err
	}
	reopened, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer reopened.Close()
	svc2 := service.New(reopened)
	if err := svc2.Recover(ctx); err != nil {
		return err
	}
	restored, err := svc2.GetVersion(ctx, verID)
	if err != nil {
		return err
	}
	if restored.Status != model.VersionPublished || restored.ContentHash != ver.ContentHash {
		return fmt.Errorf("snapshot recovery mismatch")
	}
	verify, err := svc2.VerifyVersion(ctx, verID)
	if err != nil {
		return err
	}
	if match, _ := verify["match"].(bool); !match {
		return fmt.Errorf("version hash verification failed after restart")
	}
	return nil
}
