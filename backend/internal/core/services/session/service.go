package session

import (
	"context"
	"time"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
)

type Service struct {
	sessionRepo ports.SessionRepository
}

func NewService(sessionRepo ports.SessionRepository) *Service {
	return &Service{
		sessionRepo: sessionRepo,
	}
}

func (s *Service) CreateSession(
	ctx context.Context,
	id string,
) (*domain.Session, error) {
	session := &domain.Session{
		ID:           id,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		IsActive:     true,
	}

	err := s.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) GetSessionByID(
	ctx context.Context,
	id string,
) (*domain.Session, error) {
	return s.sessionRepo.GetByID(ctx, id)
}

func (s *Service) UpdateSession(
	ctx context.Context,
	session *domain.Session,
) error {
	return s.sessionRepo.Update(ctx, session)
}

func (s *Service) DeleteSession(
	ctx context.Context,
	id string,
) error {
	return s.sessionRepo.Delete(ctx, id)
}

func (s *Service) UpdateLastAccessed(
	ctx context.Context,
	id string,
) error {
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	session.LastAccessed = time.Now()
	return s.sessionRepo.Update(ctx, session)
}
