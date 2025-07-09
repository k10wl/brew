package services

import (
	"context"
	"time"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
	"brew/internal/utils/logger"
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
	logger.Debug("Creating session", "id", id)

	session := &domain.Session{
		ID:           id,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		IsActive:     true,
	}

	err := s.sessionRepo.Save(ctx, session)
	if err != nil {
		logger.Error("Failed to save session", "error", err, "id", id)
		return nil, err
	}

	logger.Debug("Session created successfully", "id", id)
	return session, nil
}

func (s *SessionService) GetSessionByID(
	ctx context.Context,
	id string,
) (*domain.Session, error) {
	logger.Debug("Getting session by ID", "id", id)
	return s.sessionRepo.GetByID(ctx, id)
}

func (s *SessionService) UpdateSession(
	ctx context.Context,
	session *domain.Session,
) error {
	logger.Debug("Updating session", "id", session.ID)
	return s.sessionRepo.Update(ctx, session)
}

func (s *SessionService) DeleteSession(
	ctx context.Context,
	id string,
) error {
	logger.Debug("Deleting session", "id", id)
	return s.sessionRepo.Delete(ctx, id)
}

func (s *SessionService) UpdateLastAccessed(
	ctx context.Context,
	id string,
) error {
	logger.Debug("Updating last accessed", "id", id)

	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get session for updating last accessed", "error", err, "id", id)
		return err
	}

	session.LastAccessed = time.Now()
	return s.sessionRepo.Update(ctx, session)
}
