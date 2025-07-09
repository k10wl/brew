package services

import (
	"context"
	"time"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
)

type SessionService struct {
	sessionRepo ports.SessionRepository
}

func NewSessionService(sessionRepo ports.SessionRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionService) CreateSession(
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

func (s *SessionService) GetSessionByID(
	ctx context.Context,
	id string,
) (*domain.Session, error) {
	return s.sessionRepo.GetByID(ctx, id)
}

func (s *SessionService) UpdateSession(
	ctx context.Context,
	session *domain.Session,
) error {
	return s.sessionRepo.Update(ctx, session)
}

func (s *SessionService) DeleteSession(
	ctx context.Context,
	id string,
) error {
	return s.sessionRepo.Delete(ctx, id)
}

func (s *SessionService) UpdateLastAccessed(
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
