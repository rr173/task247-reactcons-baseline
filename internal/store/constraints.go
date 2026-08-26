package store

import (
	"context"
	"encoding/json"

	"task247-reactcons/internal/model"
)

func (s *Store) CreateConstraint(ctx context.Context, netID int64, kind, description string, payload []byte) (model.Constraint, error) {
	if len(payload) == 0 {
		payload = []byte("{}")
	}
	ts := now()
	const q = `INSERT INTO constraints (network_id, kind, status, description, payload_json, created_at)
		VALUES (?, ?, 'pending', ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, q, netID, kind, description, string(payload), ts)
	if err != nil {
		return model.Constraint{}, err
	}
	id, _ := res.LastInsertId()
	return model.Constraint{
		ID:          id,
		NetworkID:   netID,
		Kind:        kind,
		Status:      "pending",
		Description: description,
		Payload:     json.RawMessage(payload),
		CreatedAt:   ts,
	}, nil
}

func (s *Store) ListConstraints(ctx context.Context, networkID int64) ([]model.Constraint, error) {
	const q = `SELECT id, network_id, kind, status, description, payload_json, created_at
		FROM constraints WHERE network_id = ? ORDER BY id`
	rows, err := s.db.QueryContext(ctx, q, networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Constraint
	for rows.Next() {
		var c model.Constraint
		var payload string
		if err := rows.Scan(&c.ID, &c.NetworkID, &c.Kind, &c.Status, &c.Description, &payload, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.Payload = json.RawMessage(payload)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) UpdateConstraintStatus(ctx context.Context, id int64, status string) error {
	const q = `UPDATE constraints SET status = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, q, status, id)
	return err
}
