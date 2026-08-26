package store

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Store is the SQLite-backed persistence layer.
type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }
func (s *Store) DB() *sql.DB  { return s.db }

func (s *Store) migrate(ctx context.Context) error {
	statements := []string{
		`PRAGMA foreign_keys=ON`,
		`CREATE TABLE IF NOT EXISTS networks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ext_key TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS species (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			network_id INTEGER NOT NULL,
			symbol TEXT NOT NULL,
			name TEXT NOT NULL DEFAULT '',
			charge INTEGER NOT NULL DEFAULT 0,
			composition_json TEXT NOT NULL DEFAULT '{}',
			status TEXT NOT NULL,
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			UNIQUE(network_id, symbol),
			FOREIGN KEY(network_id) REFERENCES networks(id))`,
		`CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			network_id INTEGER NOT NULL,
			equation TEXT NOT NULL,
			reversible INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL,
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			UNIQUE(network_id, equation),
			FOREIGN KEY(network_id) REFERENCES networks(id))`,
		`CREATE TABLE IF NOT EXISTS reaction_participants (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			reaction_id INTEGER NOT NULL,
			species_id INTEGER NOT NULL,
			role TEXT NOT NULL,
			coeff REAL NOT NULL,
			UNIQUE(reaction_id, species_id, role),
			FOREIGN KEY(reaction_id) REFERENCES reactions(id),
			FOREIGN KEY(species_id) REFERENCES species(id))`,
		`CREATE TABLE IF NOT EXISTS conserved_pools (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			network_id INTEGER NOT NULL,
			label TEXT NOT NULL,
			vector_json TEXT NOT NULL,
			rank INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY(network_id) REFERENCES networks(id))`,
		`CREATE TABLE IF NOT EXISTS conflict_sets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			network_id INTEGER NOT NULL,
			kind TEXT NOT NULL,
			target TEXT NOT NULL,
			reaction_ids_json TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY(network_id) REFERENCES networks(id))`,
		`CREATE TABLE IF NOT EXISTS boundaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			network_id INTEGER NOT NULL,
			species_id INTEGER NOT NULL,
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			UNIQUE(network_id, species_id),
			FOREIGN KEY(network_id) REFERENCES networks(id),
			FOREIGN KEY(species_id) REFERENCES species(id))`,
		`CREATE TABLE IF NOT EXISTS constraints (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			network_id INTEGER NOT NULL,
			kind TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			description TEXT NOT NULL DEFAULT '',
			payload_json TEXT NOT NULL DEFAULT '{}',
			created_at TEXT NOT NULL,
			FOREIGN KEY(network_id) REFERENCES networks(id))`,
		`CREATE TABLE IF NOT EXISTS versions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			network_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			content_hash TEXT NOT NULL,
			payload_json TEXT NOT NULL,
			created_at TEXT NOT NULL,
			published_at TEXT,
			UNIQUE(network_id, content_hash),
			FOREIGN KEY(network_id) REFERENCES networks(id))`,
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func now() string                               { return time.Now().UTC().Format(time.RFC3339Nano) }
func parseTime(value string) (time.Time, error) { return time.Parse(time.RFC3339Nano, value) }
