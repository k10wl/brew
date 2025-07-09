package services

import (
	"context"
	"fmt"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
)

type JarService struct {
	jarRepo       ports.JarRepository
	sessionRepo   ports.SessionRepository
	identifierGen ports.IdentifierGenerator
}

func NewJarService(
	jarRepo ports.JarRepository,
	sessionRepo ports.SessionRepository,
	identifierGen ports.IdentifierGenerator,
) *JarService {
	return &JarService{
		jarRepo:       jarRepo,
		sessionRepo:   sessionRepo,
		identifierGen: identifierGen,
	}
}

func (s *JarService) CreateJar(
	ctx context.Context,
	name string,
	sessionID string,
) (*domain.Jar, error) {
	id, err := s.identifierGen.Generate(ctx, name)
	if err != nil {
		return nil, err
	}

	exists, err := s.jarRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("jar with id %s already exists", id)
	}

	jar := &domain.Jar{
		ID:   id,
		Name: name,
	}

	err = s.jarRepo.Save(ctx, jar)
	if err != nil {
		return nil, err
	}

	return jar, nil
}
