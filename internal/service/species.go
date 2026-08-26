package service

import (
	"context"
	"fmt"

	"task247-reactcons/internal/chem"
	"task247-reactcons/internal/model"
)

// AddSpecies parses the chemical formula, derives element composition and charge,
// and stores a new species for the network.
func (s *Service) AddSpecies(ctx context.Context, networkID int64, symbol, name string) (model.Species, error) {
	comp, charge, err := chem.ParseFormula(symbol)
	if err != nil {
		return model.Species{}, fmt.Errorf("%w: %v", model.ErrFormula, err)
	}
	exists, err := s.store.NetworkExists(ctx, networkID)
	if err != nil {
		return model.Species{}, err
	}
	if !exists {
		return model.Species{}, model.ErrNotFound
	}
	return s.store.CreateSpecies(ctx, model.Species{
		NetworkID:   networkID,
		Symbol:      symbol,
		Name:        name,
		Charge:      charge,
		Composition: comp,
		Status:      model.SpeciesValid,
	})
}

func (s *Service) GetSpecies(ctx context.Context, id int64) (model.Species, error) {
	return s.store.GetSpecies(ctx, id)
}

func (s *Service) ListSpecies(ctx context.Context, networkID int64) ([]model.Species, error) {
	return s.store.ListSpecies(ctx, networkID)
}

func (s *Service) ReviseSpecies(ctx context.Context, id int64, status model.SpeciesStatus, note string) error {
	return s.store.UpdateSpeciesStatus(ctx, id, status, note)
}
