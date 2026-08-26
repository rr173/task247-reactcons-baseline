package store

import (
	"context"
	"database/sql"
	"encoding/json"

	"task247-reactcons/internal/model"
)

func (s *Store) CreateSpecies(ctx context.Context, sp model.Species) (model.Species, error) {
	if sp.Status == "" {
		sp.Status = model.SpeciesValid
	}
	compJSON, err := json.Marshal(sp.Composition)
	if err != nil {
		return model.Species{}, err
	}
	ts := now()
	const q = `INSERT INTO species (network_id, symbol, name, charge, composition_json, status, note, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, q, sp.NetworkID, sp.Symbol, sp.Name, sp.Charge, string(compJSON), sp.Status, sp.Note, ts)
	if err != nil {
		if isUniqueErr(err) {
			return model.Species{}, model.ErrConflict
		}
		return model.Species{}, err
	}
	id, _ := res.LastInsertId()
	sp.ID = id
	sp.CreatedAt = ts
	return sp, nil
}

func (s *Store) GetSpecies(ctx context.Context, id int64) (model.Species, error) {
	const q = `SELECT id, network_id, symbol, name, charge, composition_json, status, note, created_at
		FROM species WHERE id = ?`
	row := s.db.QueryRowContext(ctx, q, id)
	return scanSpecies(row)
}

func (s *Store) ListSpecies(ctx context.Context, networkID int64) ([]model.Species, error) {
	const q = `SELECT id, network_id, symbol, name, charge, composition_json, status, note, created_at
		FROM species WHERE network_id = ? ORDER BY id`
	rows, err := s.db.QueryContext(ctx, q, networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Species
	for rows.Next() {
		sp, err := scanSpecies(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sp)
	}
	return out, rows.Err()
}

func (s *Store) UpdateSpeciesStatus(ctx context.Context, id int64, status model.SpeciesStatus, note string) error {
	const q = `UPDATE species SET status = ?, note = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, q, status, note, id)
	return err
}

func (s *Store) SpeciesBySymbol(ctx context.Context, networkID int64, symbol string) (model.Species, bool, error) {
	const q = `SELECT id, network_id, symbol, name, charge, composition_json, status, note, created_at
		FROM species WHERE network_id = ? AND symbol = ?`
	row := s.db.QueryRowContext(ctx, q, networkID, symbol)
	sp, err := scanSpecies(row)
	if err == model.ErrNotFound {
		return model.Species{}, false, nil
	}
	if err != nil {
		return model.Species{}, false, err
	}
	return sp, true, nil
}

func scanSpecies(row scanner) (model.Species, error) {
	var sp model.Species
	var compJSON string
	var status string
	err := row.Scan(&sp.ID, &sp.NetworkID, &sp.Symbol, &sp.Name, &sp.Charge, &compJSON, &status, &sp.Note, &sp.CreatedAt)
	if err == sql.ErrNoRows {
		return model.Species{}, model.ErrNotFound
	}
	if err != nil {
		return model.Species{}, err
	}
	if sp.Composition == nil {
		sp.Composition = map[string]int{}
	}
	if err := json.Unmarshal([]byte(compJSON), &sp.Composition); err != nil {
		return model.Species{}, err
	}
	sp.Status = model.SpeciesStatus(status)
	return sp, nil
}
