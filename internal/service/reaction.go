package service

import (
	"context"
	"fmt"

	"task247-reactcons/internal/chem"
	"task247-reactcons/internal/model"
)

// AddReaction parses an equation, resolves each species symbol against the network,
// and persists the reaction with its participants.
func (s *Service) AddReaction(ctx context.Context, networkID int64, equation string, reversible bool) (model.Reaction, error) {
	reactants, products, rev, err := chem.ParseEquation(equation)
	if err != nil {
		return model.Reaction{}, fmt.Errorf("%w: %v", model.ErrEquation, err)
	}
	if reversible {
		rev = true
	}
	exists, err := s.store.NetworkExists(ctx, networkID)
	if err != nil {
		return model.Reaction{}, err
	}
	if !exists {
		return model.Reaction{}, model.ErrNotFound
	}
	participants := make([]model.Participant, 0, len(reactants)+len(products))
	for _, t := range append(append([]chem.Term{}, reactants...), products...) {
		sp, found, err := s.store.SpeciesBySymbol(ctx, networkID, t.Symbol)
		if err != nil {
			return model.Reaction{}, err
		}
		if !found {
			return model.Reaction{}, fmt.Errorf("%w: unknown species %q in network %d", model.ErrInvalidArgument, t.Symbol, networkID)
		}
		participants = append(participants, model.Participant{
			SpeciesID: sp.ID,
			Symbol:    t.Symbol,
			Role:      t.Role,
			Coeff:     t.Coeff,
		})
	}
	return s.store.CreateReaction(ctx, model.Reaction{
		NetworkID:  networkID,
		Equation:   equation,
		Reversible: rev,
		Status:     model.ReactionCandidate,
		Participants: participants,
	})
}

func (s *Service) GetReaction(ctx context.Context, id int64) (model.Reaction, error) {
	return s.store.GetReaction(ctx, id)
}

func (s *Service) ListReactions(ctx context.Context, networkID int64) ([]model.Reaction, error) {
	return s.store.ListReactions(ctx, networkID)
}

// ExemptReaction marks a reaction as exempt from conservation checks (e.g. an
// open-system boundary reaction the researcher has judged acceptable).
func (s *Service) ExemptReaction(ctx context.Context, id int64, reason string) error {
	return s.store.ExemptReaction(ctx, id, reason)
}
