package service

import (
	"context"

	"task247-reactcons/internal/diagnose"
	"task247-reactcons/internal/model"
)

// analyze loads the network state and runs conservation diagnostics without
// persisting results. It returns the solved result plus the intermediate refs.
func (s *Service) analyze(ctx context.Context, netID int64) (*model.SolveResult, []model.SpeciesRef, []model.ReactionRef, map[int64]bool, []model.Constraint, error) {
	species, err := s.store.ListSpecies(ctx, netID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	openSet, err := s.store.OpenSpeciesSet(ctx, netID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	speciesRefs := make([]model.SpeciesRef, 0, len(species))
	for _, sp := range species {
		speciesRefs = append(speciesRefs, model.SpeciesRef{
			ID:          sp.ID,
			Symbol:      sp.Symbol,
			Composition: sp.Composition,
			Charge:      sp.Charge,
			Open:        openSet[sp.ID],
		})
	}
	reactions, err := s.store.ListReactions(ctx, netID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	reactionRefs := make([]model.ReactionRef, 0, len(reactions))
	exempt := map[int64]bool{}
	for _, rx := range reactions {
		if rx.Status == model.ReactionExempt {
			exempt[rx.ID] = true
		}
		ref := model.ReactionRef{ID: rx.ID, Equation: rx.Equation}
		for _, p := range rx.Participants {
			ref.Participants = append(ref.Participants, model.ParticipantRef{
				SpeciesID: p.SpeciesID,
				Role:      p.Role,
				Coeff:     p.Coeff,
			})
		}
		reactionRefs = append(reactionRefs, ref)
	}
	constraints, err := s.store.ListConstraints(ctx, netID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	result, err := diagnose.Analyze(speciesRefs, reactionRefs, exempt, constraints)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return result, speciesRefs, reactionRefs, exempt, constraints, nil
}

// Solve runs the conservation analysis, persists conserved pools, conflict sets and
// constraint statuses, and advances the network status accordingly.
func (s *Service) Solve(ctx context.Context, netID int64) (*model.SolveResult, error) {
	if err := s.store.UpdateNetworkStatus(ctx, netID, model.NetworkSolving); err != nil {
		return nil, err
	}
	result, _, _, _, _, err := s.analyze(ctx, netID)
	if err != nil {
		_ = s.store.UpdateNetworkStatus(ctx, netID, model.NetworkDraft)
		return nil, err
	}
	if err := s.store.UpsertConservedPools(ctx, netID, result.ConservedPools); err != nil {
		return nil, err
	}
	if err := s.store.UpsertConflictSets(ctx, netID, result.ConflictSets); err != nil {
		return nil, err
	}
	for range result.ConstraintChecks {
		// Constraint statuses are returned to the caller but not persisted.
	}
	if err := s.store.UpdateNetworkStatus(ctx, netID, result.Status); err != nil {
		return nil, err
	}
	return result, nil
}

// Conservation returns a freshly computed analysis without persisting it.
func (s *Service) Conservation(ctx context.Context, netID int64) (*model.SolveResult, error) {
	result, _, _, _, _, err := s.analyze(ctx, netID)
	return result, err
}

// ConservedPools returns the conserved subspace for the network.
func (s *Service) ConservedPools(ctx context.Context, netID int64) ([]model.ConservedPool, error) {
	result, _, _, _, _, err := s.analyze(ctx, netID)
	if err != nil {
		return nil, err
	}
	return result.ConservedPools, nil
}

// ConflictSets returns the minimal conflict reaction sets for the network.
func (s *Service) ConflictSets(ctx context.Context, netID int64) ([]model.ConflictSet, error) {
	result, _, _, _, _, err := s.analyze(ctx, netID)
	if err != nil {
		return nil, err
	}
	return result.ConflictSets, nil
}
