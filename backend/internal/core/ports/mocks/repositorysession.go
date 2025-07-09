package mocks

import (
	"context"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
)

var _ ports.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	SaveFunc    func(ctx context.Context, session *domain.Session) error
	GetByIDFunc func(ctx context.Context, id string) (*domain.Session, error)
	UpdateFunc  func(ctx context.Context, session *domain.Session) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *SessionRepository) Save(ctx context.Context, session *domain.Session) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, session)
	}
	return nil
}

func (m *SessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *SessionRepository) Update(ctx context.Context, session *domain.Session) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, session)
	}
	return nil
}

func (m *SessionRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
