package service

import (
	"context"
	"encoding/json"
	"fmt"

	"task247-reactcons/internal/model"
)

// AddConstraint declares an experimental constraint, currently a moiety-closure:
// a linear combination of species amounts that the network must conserve.
func (s *Service) AddConstraint(ctx context.Context, networkID int64, closure model.MoietyClosure) (model.Constraint, error) {
	if len(closure.Members) == 0 {
		return model.Constraint{}, fmt.Errorf("%w: closure has no members", model.ErrInvalidArgument)
	}
	existing, err := s.store.ListSpecies(ctx, networkID)
	if err != nil {
		return model.Constraint{}, err
	}
	present := map[string]bool{}
	for _, sp := range existing {
		present[sp.Symbol] = true
	}
	for _, m := range closure.Members {
		if !present[m.Symbol] {
			return model.Constraint{}, fmt.Errorf("%w: closure references unknown species %q", model.ErrInvalidArgument, m.Symbol)
		}
	}
	payload, err := json.Marshal(closure)
	if err != nil {
		return model.Constraint{}, err
	}
	return s.store.CreateConstraint(ctx, networkID, "moietyclosure", "declared conserved moiety", payload)
}

func (s *Service) ListConstraints(ctx context.Context, networkID int64) ([]model.Constraint, error) {
	return s.store.ListConstraints(ctx, networkID)
}
