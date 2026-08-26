package service

import (
	"context"
	"runtime"

	"task247-reactcons/internal/model"
	"task247-reactcons/internal/store"
)

// Service orchestrates networks, species, reactions, conservation analysis and
// immutable versioning on top of the persistence layer.
type Service struct {
	store *store.Store
}

func New(repository *store.Store) *Service {
	return &Service{store: repository}
}

// SelfCheck reports liveness and aggregate counts.
func (s *Service) SelfCheck(ctx context.Context) (map[string]any, error) {
	stats, err := s.store.Stats(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"module":     "task247-reactcons",
		"go_version": runtime.Version(),
		"sqlite":     "3.46.1",
		"status":     "ok",
		"counts":     stats,
	}, nil
}

// Recover resumes any network left mid-solve and verifies the store is usable.
func (s *Service) Recover(ctx context.Context) error {
	nets, err := s.store.ListNetworks(ctx)
	if err != nil {
		return err
	}
	for _, n := range nets {
		if n.Status == model.NetworkSolving {
			if _, err := s.Solve(ctx, n.ID); err != nil {
				return err
			}
		}
	}
	return nil
}
