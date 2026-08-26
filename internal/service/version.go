package service

import (
	"context"
	"encoding/json"
	"fmt"

	"task247-reactcons/internal/model"
	"task247-reactcons/internal/version"
)

// SnapshotPayload is the frozen, content-addressed state of a network version.
type SnapshotPayload struct {
	Network     model.Network      `json:"network"`
	Species     []model.Species    `json:"species"`
	Reactions   []model.Reaction   `json:"reactions"`
	Boundaries  []model.Boundary   `json:"boundaries"`
	Constraints []model.Constraint `json:"constraints"`
	Solve       model.SolveResult  `json:"solve"`
}

// PublishVersion freezes the current network state as an immutable, content-addressed
// version. It fails if the network is not in a publishable (conservation-clean) state.
func (s *Service) PublishVersion(ctx context.Context, netID int64) (model.NetworkVersion, error) {
	result, _, _, _, _, err := s.analyze(ctx, netID)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	if result.Status != model.NetworkPublishable {
		return model.NetworkVersion{}, fmt.Errorf("%w: network has conservation conflicts", model.ErrInvalidState)
	}
	net, err := s.store.GetNetwork(ctx, netID)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	species, err := s.store.ListSpecies(ctx, netID)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	reactions, err := s.store.ListReactions(ctx, netID)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	boundaries, err := s.store.ListBoundaries(ctx, netID)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	constraints, err := s.store.ListConstraints(ctx, netID)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	payload := SnapshotPayload{
		Network:     net,
		Species:     species,
		Reactions:   reactions,
		Boundaries:  boundaries,
		Constraints: constraints,
		Solve:       *result,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	hash, err := version.ContentHash(data)
	if err != nil {
		return model.NetworkVersion{}, err
	}
	return s.store.CreateVersion(ctx, netID, model.VersionPublished, hash, data)
}

func (s *Service) ListVersions(ctx context.Context, netID int64) ([]model.NetworkVersion, error) {
	return s.store.ListVersions(ctx, netID)
}

func (s *Service) GetVersion(ctx context.Context, id int64) (model.NetworkVersion, error) {
	return s.store.GetVersion(ctx, id)
}

// VerifyVersion recomputes the snapshot hash from the stored payload and compares it
// to the recorded content hash, detecting any tampering with the frozen state.
func (s *Service) VerifyVersion(ctx context.Context, id int64) (map[string]any, error) {
	v, err := s.store.GetVersion(ctx, id)
	if err != nil {
		return nil, err
	}
	expected, err := version.ContentHash(append([]byte(v.Payload), '\n'))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"version_id":      v.ID,
		"stored_hash":     v.ContentHash,
		"recomputed_hash": expected,
		"match":           expected == v.ContentHash,
	}, nil
}
