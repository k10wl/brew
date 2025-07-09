package mocks

import (
	"context"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
)

var _ ports.BrewRepository = (*BrewRepository)(nil)

type BrewRepository struct {
	SaveFunc           func(ctx context.Context, brew *domain.Brew) error
	GetByIDFunc        func(ctx context.Context, id string) (*domain.Brew, error)
	GetBySessionIDFunc func(ctx context.Context, sessionID string, pointer *string, limit int) (*ports.PaginatedResult[*domain.Brew], error)
	UpdateFunc         func(ctx context.Context, brew *domain.Brew) error
	ExistsFunc         func(ctx context.Context, id string) (bool, error)
}

func (m *BrewRepository) Save(ctx context.Context, brew *domain.Brew) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, brew)
	}
	return nil
}

func (m *BrewRepository) GetByID(ctx context.Context, id string) (*domain.Brew, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *BrewRepository) GetBySessionID(ctx context.Context, sessionID string, pointer *string, limit int) (*ports.PaginatedResult[*domain.Brew], error) {
	if m.GetBySessionIDFunc != nil {
		return m.GetBySessionIDFunc(ctx, sessionID, pointer, limit)
	}
	return &ports.PaginatedResult[*domain.Brew]{}, nil
}

func (m *BrewRepository) Update(ctx context.Context, brew *domain.Brew) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, brew)
	}
	return nil
}

func (m *BrewRepository) Exists(ctx context.Context, id string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, id)
	}
	return false, nil
}
