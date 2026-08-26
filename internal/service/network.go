package service

import (
	"context"

	"task247-reactcons/internal/model"
)

func (s *Service) CreateNetwork(ctx context.Context, net model.Network) (model.Network, error) {
	return s.store.CreateNetwork(ctx, net)
}

func (s *Service) GetNetwork(ctx context.Context, id int64) (model.Network, error) {
	return s.store.GetNetwork(ctx, id)
}

func (s *Service) ListNetworks(ctx context.Context) ([]model.Network, error) {
	return s.store.ListNetworks(ctx)
}
