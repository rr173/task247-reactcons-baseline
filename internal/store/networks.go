package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"task247-reactcons/internal/model"
)

func (s *Store) CreateNetwork(ctx context.Context, net model.Network) (model.Network, error) {
	if net.ExtKey == "" {
		net.ExtKey = uuid.NewString()
	}
	if net.Status == "" {
		net.Status = model.NetworkDraft
	}
	ts := now()
	const q = `INSERT INTO networks (ext_key, name, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, q, net.ExtKey, net.Name, net.Description, net.Status, ts, ts)
	if err != nil {
		if isUniqueErr(err) {
			return model.Network{}, model.ErrConflict
		}
		return model.Network{}, err
	}
	id, _ := res.LastInsertId()
	net.ID = id
	net.CreatedAt = ts
	net.UpdatedAt = ts
	return net, nil
}

func (s *Store) GetNetwork(ctx context.Context, id int64) (model.Network, error) {
	const q = `SELECT id, ext_key, name, description, status, created_at, updated_at FROM networks WHERE id = ?`
	row := s.db.QueryRowContext(ctx, q, id)
	return scanNetwork(row)
}

func (s *Store) GetNetworkByExtKey(ctx context.Context, extKey string) (model.Network, error) {
	const q = `SELECT id, ext_key, name, description, status, created_at, updated_at FROM networks WHERE ext_key = ?`
	row := s.db.QueryRowContext(ctx, q, extKey)
	return scanNetwork(row)
}

func (s *Store) ListNetworks(ctx context.Context) ([]model.Network, error) {
	const q = `SELECT id, ext_key, name, description, status, created_at, updated_at FROM networks ORDER BY id`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Network
	for rows.Next() {
		n, err := scanNetwork(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (s *Store) UpdateNetworkStatus(ctx context.Context, id int64, status model.NetworkStatus) error {
	const q = `UPDATE networks SET status = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, q, status, now(), id)
	return err
}

func (s *Store) SetNetworkMeta(ctx context.Context, id int64, name, description string) error {
	const q = `UPDATE networks SET name = ?, description = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, q, name, description, now(), id)
	return err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanNetwork(row scanner) (model.Network, error) {
	var n model.Network
	var status string
	err := row.Scan(&n.ID, &n.ExtKey, &n.Name, &n.Description, &status, &n.CreatedAt, &n.UpdatedAt)
	if err == sql.ErrNoRows {
		return model.Network{}, model.ErrNotFound
	}
	if err != nil {
		return model.Network{}, err
	}
	n.Status = model.NetworkStatus(status)
	return n, nil
}

func isUniqueErr(err error) bool {
	if err == nil {
		return false
	}
	return containsAny(err.Error(), "UNIQUE constraint failed", "duplicate")
}
