package store

import (
	"context"
	"database/sql"
	"encoding/json"

	"task247-reactcons/internal/model"
)

// CreateVersion supersedes any previously published version of the network and
// inserts a new immutable snapshot. When status is published, published_at is set.
func (s *Store) CreateVersion(ctx context.Context, netID int64, status model.VersionStatus, hash string, payload []byte) (model.NetworkVersion, error) {
	// Publishing the same content is idempotent: the unique content hash is the
	// version identity, and an existing snapshot must not be superseded before
	// the lookup has had a chance to return it.
	if existing, err := s.findVersionByHash(ctx, netID, hash); err != nil {
		return model.NetworkVersion{}, err
	} else if existing != nil {
		return *existing, nil
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE versions SET status = ? WHERE network_id = ? AND status = ?`,
		model.VersionSuperseded, netID, model.VersionPublished); err != nil {
		return model.NetworkVersion{}, err
	}
	ts := now()
	var publishedAt sql.NullString
	if status == model.VersionPublished {
		publishedAt = sql.NullString{String: ts, Valid: true}
	}
	const q = `INSERT INTO versions (network_id, status, content_hash, payload_json, created_at, published_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, q, netID, status, hash, string(payload), ts, publishedAt)
	if err != nil {
		if isUniqueErr(err) {
			return model.NetworkVersion{}, model.ErrConflict
		}
		return model.NetworkVersion{}, err
	}
	id, _ := res.LastInsertId()
	return model.NetworkVersion{
		ID:          id,
		NetworkID:   netID,
		Status:      status,
		ContentHash: hash,
		Payload:     json.RawMessage(payload),
		CreatedAt:   ts,
		PublishedAt: publishedAt.String,
	}, nil
}

func (s *Store) findVersionByHash(ctx context.Context, netID int64, hash string) (*model.NetworkVersion, error) {
	const q = `SELECT id, network_id, status, content_hash, payload_json, created_at, published_at
		FROM versions WHERE network_id = ? AND content_hash = ?`
	row := s.db.QueryRowContext(ctx, q, netID, hash)
	v, err := scanVersion(row)
	if err == model.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Store) GetVersion(ctx context.Context, id int64) (model.NetworkVersion, error) {
	const q = `SELECT id, network_id, status, content_hash, payload_json, created_at, published_at
		FROM versions WHERE id = ?`
	row := s.db.QueryRowContext(ctx, q, id)
	return scanVersion(row)
}

func (s *Store) ListVersions(ctx context.Context, networkID int64) ([]model.NetworkVersion, error) {
	const q = `SELECT id, network_id, status, content_hash, payload_json, created_at, published_at
		FROM versions WHERE network_id = ? ORDER BY id`
	rows, err := s.db.QueryContext(ctx, q, networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.NetworkVersion
	for rows.Next() {
		v, err := scanVersion(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func scanVersion(row scanner) (model.NetworkVersion, error) {
	var v model.NetworkVersion
	var payload string
	var publishedAt sql.NullString
	err := row.Scan(&v.ID, &v.NetworkID, &v.Status, &v.ContentHash, &payload, &v.CreatedAt, &publishedAt)
	if err == sql.ErrNoRows {
		return model.NetworkVersion{}, model.ErrNotFound
	}
	if err != nil {
		return model.NetworkVersion{}, err
	}
	v.Payload = json.RawMessage(payload)
	v.PublishedAt = publishedAt.String
	return v, nil
}
