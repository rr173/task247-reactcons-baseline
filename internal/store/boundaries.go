package store

import (
	"context"

	"task247-reactcons/internal/model"
)

func (s *Store) CreateBoundary(ctx context.Context, b model.Boundary) (model.Boundary, error) {
	ts := now()
	const q = `INSERT INTO boundaries (network_id, species_id, note, created_at) VALUES (?, ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, q, b.NetworkID, b.SpeciesID, b.Note, ts)
	if err != nil {
		if isUniqueErr(err) {
			return model.Boundary{}, model.ErrConflict
		}
		return model.Boundary{}, err
	}
	id, _ := res.LastInsertId()
	b.ID = id
	b.CreatedAt = ts
	return b, nil
}

func (s *Store) ListBoundaries(ctx context.Context, networkID int64) ([]model.Boundary, error) {
	const q = `SELECT b.id, b.network_id, b.species_id, sp.symbol, b.note, b.created_at
		FROM boundaries b LEFT JOIN species sp ON sp.id = b.species_id
		WHERE b.network_id = ? ORDER BY b.id`
	rows, err := s.db.QueryContext(ctx, q, networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Boundary
	for rows.Next() {
		var b model.Boundary
		if err := rows.Scan(&b.ID, &b.NetworkID, &b.SpeciesID, &b.Symbol, &b.Note, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (s *Store) DeleteBoundary(ctx context.Context, networkID, speciesID int64) error {
	const q = `DELETE FROM boundaries WHERE network_id = ? AND species_id = ?`
	_, err := s.db.ExecContext(ctx, q, networkID, speciesID)
	return err
}

func (s *Store) OpenSpeciesSet(ctx context.Context, networkID int64) (map[int64]bool, error) {
	boundaries, err := s.ListBoundaries(ctx, networkID)
	if err != nil {
		return nil, err
	}
	set := make(map[int64]bool, len(boundaries))
	for _, b := range boundaries {
		set[b.SpeciesID] = true
	}
	return set, nil
}
