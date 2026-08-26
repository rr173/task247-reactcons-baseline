package store

import (
	"context"
	"encoding/json"

	"task247-reactcons/internal/model"
)

func (s *Store) UpsertConflictSets(ctx context.Context, networkID int64, sets []model.ConflictSet) error {
	for _, cs := range sets {
		data, err := json.Marshal(cs.ReactionIDs)
		if err != nil {
			return err
		}
		const q = `INSERT INTO conflict_sets (network_id, kind, target, reaction_ids_json, created_at) VALUES (?, ?, ?, ?, ?)`
		if _, err := s.db.ExecContext(ctx, q, networkID, cs.Kind, cs.Target, string(data), now()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListConflictSets(ctx context.Context, networkID int64) ([]model.ConflictSet, error) {
	const q = `SELECT id, kind, target, reaction_ids_json FROM conflict_sets WHERE network_id = ? ORDER BY kind, target`
	rows, err := s.db.QueryContext(ctx, q, networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.ConflictSet
	for rows.Next() {
		var (
			id     int64
			kind   string
			target string
			ids    string
		)
		if err := rows.Scan(&id, &kind, &target, &ids); err != nil {
			return nil, err
		}
		var reactionIDs []int64
		if err := json.Unmarshal([]byte(ids), &reactionIDs); err != nil {
			return nil, err
		}
		out = append(out, model.ConflictSet{ID: id, Kind: kind, Target: target, ReactionIDs: reactionIDs})
	}
	return out, rows.Err()
}
