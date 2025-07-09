package services

import (
	"context"
	"fmt"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
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
	id, err := s.identifierGen.Generate(ctx, name)
	if err != nil {
		return nil, err
	}

	exists, err := s.brewRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("brew with id %s already exists", id)
	}

	brew := &domain.Brew{
		ID:   id,
		Name: name,
	}

	err = s.brewRepo.Save(ctx, brew)
	if err != nil {
		return nil, err
	}

	return brew, nil
}
