package services

import (
	"context"
	"fmt"
	"log/slog"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
	_ "brew/internal/util"
)

type BrewService struct {
	brewRepo      ports.BrewRepository
	sessionRepo   ports.SessionRepository
	identifierGen ports.IdentifierGenerator
}

func NewBrewService(
	brewRepo ports.BrewRepository,
	sessionRepo ports.SessionRepository,
	identifierGen ports.IdentifierGenerator,
) *BrewService {
	return &BrewService{
		brewRepo:      brewRepo,
		sessionRepo:   sessionRepo,
		identifierGen: identifierGen,
	}
}

func (s *BrewService) CreateBrew(
	ctx context.Context,
	name string,
	sessionID string,
) (*domain.Brew, error) {
	slog.Debug("Creating brew", "name", name, "session_id", sessionID)

	id, err := s.identifierGen.Generate(ctx, name)
	if err != nil {
		slog.Error("Failed to generate identifier", "error", err, "name", name)
		return nil, err
	}

	exists, err := s.brewRepo.Exists(ctx, id)
	if err != nil {
		slog.Error("Failed to check if brew exists", "error", err, "id", id)
		return nil, err
	}
	if exists {
		slog.Warn("Brew already exists", "id", id)
		return nil, fmt.Errorf("brew with id %s already exists", id)
	}

	brew := &domain.Brew{
		ID:   id,
		Name: name,
	}

	err = s.brewRepo.Save(ctx, brew)
	if err != nil {
		slog.Error("Failed to save brew", "error", err, "id", id)
		return nil, err
	}

	slog.Debug("Brew created successfully", "id", id, "name", name)
	return brew, nil
}
