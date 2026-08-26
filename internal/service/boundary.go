package service

import (
	"context"

	"task247-reactcons/internal/model"
)

// MarkBoundary declares a species as an open-system exchange species; its material
// may enter or leave the system, so conservation is not required for it.
func (s *Service) MarkBoundary(ctx context.Context, networkID, speciesID int64, note string) (model.Boundary, error) {
	sp, err := s.store.GetSpecies(ctx, speciesID)
	if err != nil {
		return model.Boundary{}, err
	}
	if sp.NetworkID != networkID {
		return model.Boundary{}, model.ErrInvalidArgument
	}
	return s.store.CreateBoundary(ctx, model.Boundary{
		NetworkID: networkID,
		SpeciesID: speciesID,
		Symbol:    sp.Symbol,
		Note:      note,
	})
}

func (s *Service) ListBoundaries(ctx context.Context, networkID int64) ([]model.Boundary, error) {
	return s.store.ListBoundaries(ctx, networkID)
}

func (s *Service) RemoveBoundary(ctx context.Context, networkID, speciesID int64) error {
	return s.store.DeleteBoundary(ctx, networkID, speciesID)
}
