package store

import (
	"context"
)

// Stats returns aggregate counts used by the self-check endpoint.
func (s *Store) Stats(ctx context.Context) (map[string]int64, error) {
	out := map[string]int64{}
	tables := []string{"networks", "species", "reactions", "conserved_pools", "conflict_sets", "boundaries", "constraints", "versions"}
	for _, t := range tables {
		var n int64
		q := `SELECT COUNT(*) FROM ` + t
		if err := s.db.QueryRowContext(ctx, q).Scan(&n); err != nil {
			return nil, err
		}
		out[t] = n
	}
	return out, nil
}

// NetworkExists reports whether a network with the given id is present.
func (s *Store) NetworkExists(ctx context.Context, id int64) (bool, error) {
	var n int64
	const q = `SELECT COUNT(*) FROM networks WHERE id = ?`
	if err := s.db.QueryRowContext(ctx, q, id).Scan(&n); err != nil {
		return false, err
	}
	return n > 0, nil
}
