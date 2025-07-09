package mocks

import (
	"context"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
)

var _ ports.JarRepository = (*JarRepository)(nil)

type JarRepository struct {
	SaveFunc           func(ctx context.Context, jar *domain.Jar) error
	GetByIDFunc        func(ctx context.Context, id string) (*domain.Jar, error)
	GetBySessionIDFunc func(ctx context.Context, sessionID string, pointer *string, limit int) (*ports.PaginatedResult[*domain.Jar], error)
	UpdateFunc         func(ctx context.Context, jar *domain.Jar) error
	ExistsFunc         func(ctx context.Context, id string) (bool, error)
}

func (m *JarRepository) Save(ctx context.Context, jar *domain.Jar) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, jar)
	}
	return nil
}

func (m *JarRepository) GetByID(ctx context.Context, id string) (*domain.Jar, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *JarRepository) GetBySessionID(ctx context.Context, sessionID string, pointer *string, limit int) (*ports.PaginatedResult[*domain.Jar], error) {
	if m.GetBySessionIDFunc != nil {
		return m.GetBySessionIDFunc(ctx, sessionID, pointer, limit)
	}
	return &ports.PaginatedResult[*domain.Jar]{}, nil
}

func (m *JarRepository) Update(ctx context.Context, jar *domain.Jar) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, jar)
	}
	return nil
}

func (m *JarRepository) Exists(ctx context.Context, id string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, id)
	}
	return false, nil
}
