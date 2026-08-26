package store

import (
	"context"
	"encoding/json"

	"task247-reactcons/internal/model"
)

type poolVectorEntry struct {
	SpeciesID int64  `json:"species_id"`
	Symbol    string `json:"symbol"`
	Coeff     string `json:"coeff"`
}

func (s *Store) UpsertConservedPools(ctx context.Context, networkID int64, pools []model.ConservedPool) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM conserved_pools WHERE network_id = ?`, networkID); err != nil {
		return err
	}
	for i, p := range pools {
		entries := make([]poolVectorEntry, 0, len(p.Members))
		for _, m := range p.Members {
			entries = append(entries, poolVectorEntry{SpeciesID: m.SpeciesID, Symbol: m.Symbol, Coeff: m.Coeff})
		}
		data, err := json.Marshal(entries)
		if err != nil {
			return err
		}
		const q = `INSERT INTO conserved_pools (network_id, label, vector_json, rank) VALUES (?, ?, ?, ?)`
		if _, err := s.db.ExecContext(ctx, q, networkID, p.Label, string(data), i); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListConservedPools(ctx context.Context, networkID int64) ([]model.ConservedPool, error) {
	const q = `SELECT id, label, vector_json FROM conserved_pools WHERE network_id = ? ORDER BY rank`
	rows, err := s.db.QueryContext(ctx, q, networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.ConservedPool
	for rows.Next() {
		var (
			id     int64
			label  string
			vector string
		)
		if err := rows.Scan(&id, &label, &vector); err != nil {
			return nil, err
		}
		var entries []poolVectorEntry
		if err := json.Unmarshal([]byte(vector), &entries); err != nil {
			return nil, err
		}
		pool := model.ConservedPool{ID: id, Label: label}
		// Resolve symbols lazily is not possible without species; store symbol in payload.
		for _, e := range entries {
			pool.Members = append(pool.Members, model.PoolMember{SpeciesID: e.SpeciesID, Symbol: e.Symbol, Coeff: e.Coeff})
		}
		out = append(out, pool)
	}
	return out, rows.Err()
}
