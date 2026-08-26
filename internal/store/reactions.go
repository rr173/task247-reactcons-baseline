package store

import (
	"context"
	"database/sql"

	"task247-reactcons/internal/model"
)

func (s *Store) CreateReaction(ctx context.Context, rx model.Reaction) (model.Reaction, error) {
	if rx.Status == "" {
		rx.Status = model.ReactionCandidate
	}
	ts := now()
	const q = `INSERT INTO reactions (network_id, equation, reversible, status, note, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, q, rx.NetworkID, rx.Equation, boolToInt(rx.Reversible), rx.Status, rx.Note, ts)
	if err != nil {
		if isUniqueErr(err) {
			return model.Reaction{}, model.ErrConflict
		}
		return model.Reaction{}, err
	}
	id, _ := res.LastInsertId()
	rx.ID = id
	rx.CreatedAt = ts
	for _, p := range rx.Participants {
		const pq = `INSERT INTO reaction_participants (reaction_id, species_id, role, coeff) VALUES (?, ?, ?, ?)`
		if _, err := s.db.ExecContext(ctx, pq, rx.ID, p.SpeciesID, p.Role, p.Coeff); err != nil {
			return model.Reaction{}, err
		}
	}
	return rx, nil
}

func (s *Store) GetReaction(ctx context.Context, id int64) (model.Reaction, error) {
	const q = `SELECT r.id, r.network_id, r.equation, r.reversible, r.status, r.note, r.created_at,
			p.species_id, p.role, p.coeff, sp.symbol
		FROM reactions r
		LEFT JOIN reaction_participants p ON p.reaction_id = r.id
		LEFT JOIN species sp ON sp.id = p.species_id
		WHERE r.id = ? ORDER BY p.id`
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return model.Reaction{}, err
	}
	defer rows.Close()
	return scanReaction(rows)
}

func (s *Store) ListReactions(ctx context.Context, networkID int64) ([]model.Reaction, error) {
	const q = `SELECT r.id, r.network_id, r.equation, r.reversible, r.status, r.note, r.created_at,
			p.species_id, p.role, p.coeff, sp.symbol
		FROM reactions r
		LEFT JOIN reaction_participants p ON p.reaction_id = r.id
		LEFT JOIN species sp ON sp.id = p.species_id
		WHERE r.network_id = ? ORDER BY r.id, p.id`
	rows, err := s.db.QueryContext(ctx, q, networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanReactions(rows)
}

func (s *Store) ExemptReaction(ctx context.Context, id int64, reason string) error {
	const q = `UPDATE reactions SET status = ?, note = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, q, model.ReactionExempt, reason, id)
	return err
}

// LoadReactionRefs returns reactions of a network as analysis references together
// with a symbol->speciesID map used for constraint resolution.
func (s *Store) LoadReactionRefs(ctx context.Context, networkID int64) ([]model.ReactionRef, map[string]int64, error) {
	reactions, err := s.ListReactions(ctx, networkID)
	if err != nil {
		return nil, nil, err
	}
	symbolToID := map[string]int64{}
	species, err := s.ListSpecies(ctx, networkID)
	if err != nil {
		return nil, nil, err
	}
	for _, sp := range species {
		symbolToID[sp.Symbol] = sp.ID
	}
	refs := make([]model.ReactionRef, 0, len(reactions))
	for _, rx := range reactions {
		ref := model.ReactionRef{ID: rx.ID, Equation: rx.Equation}
		for _, p := range rx.Participants {
			ref.Participants = append(ref.Participants, model.ParticipantRef{
				SpeciesID: p.SpeciesID,
				Role:      p.Role,
				Coeff:     p.Coeff,
			})
		}
		refs = append(refs, ref)
	}
	return refs, symbolToID, nil
}

func (s *Store) LoadSpeciesRefs(ctx context.Context, networkID int64, openSet map[int64]bool) ([]model.SpeciesRef, error) {
	species, err := s.ListSpecies(ctx, networkID)
	if err != nil {
		return nil, err
	}
	out := make([]model.SpeciesRef, 0, len(species))
	for _, sp := range species {
		out = append(out, model.SpeciesRef{
			ID:          sp.ID,
			Symbol:      sp.Symbol,
			Composition: sp.Composition,
			Charge:      sp.Charge,
			Open:        openSet[sp.ID],
		})
	}
	return out, nil
}

func scanReaction(rows *sql.Rows) (model.Reaction, error) {
	rxs, err := scanReactions(rows)
	if err != nil {
		return model.Reaction{}, err
	}
	if len(rxs) == 0 {
		return model.Reaction{}, model.ErrNotFound
	}
	return rxs[0], nil
}

func scanReactions(rows *sql.Rows) ([]model.Reaction, error) {
	var out []model.Reaction
	var cur *model.Reaction
	for rows.Next() {
		var (
			id         int64
			netID      int64
			equation   string
			rev        int
			status     string
			note       string
			createdAt  string
			spID       sql.NullInt64
			role       sql.NullString
			coeff      sql.NullFloat64
			symbol     sql.NullString
		)
		if err := rows.Scan(&id, &netID, &equation, &rev, &status, &note, &createdAt, &spID, &role, &coeff, &symbol); err != nil {
			return nil, err
		}
		if cur == nil || cur.ID != id {
			r := &model.Reaction{
				ID:         id,
				NetworkID:  netID,
				Equation:   equation,
				Reversible: rev != 0,
				Status:     model.ReactionStatus(status),
				Note:       note,
				CreatedAt:  createdAt,
			}
			out = append(out, *r)
			cur = &out[len(out)-1]
		}
		if spID.Valid {
			cur.Participants = append(cur.Participants, model.Participant{
				SpeciesID: spID.Int64,
				Symbol:    symbol.String,
				Role:      role.String,
				Coeff:     coeff.Float64,
			})
		}
	}
	return out, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
