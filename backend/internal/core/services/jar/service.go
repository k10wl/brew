package jar

import (
	"context"
	"fmt"

	"brew/internal/core/domain"
	"brew/internal/core/ports"
)

type Service struct {
	jarRepo       ports.JarRepository
	sessionRepo   ports.SessionRepository
	identifierGen ports.IdentifierGenerator
}

func NewService(
	jarRepo ports.JarRepository,
	sessionRepo ports.SessionRepository,
	identifierGen ports.IdentifierGenerator,
) *Service {
	return &Service{
		jarRepo:       jarRepo,
		sessionRepo:   sessionRepo,
		identifierGen: identifierGen,
	}
}

func (s *Service) CreateJar(
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

	jar := domain.NewJar(id, name)
	err = s.jarRepo.Save(ctx, jar)
	if err != nil {
		return nil, err
	}

	return jar, nil
}

func (s *Service) GetJarByID(
	ctx context.Context,
	id string,
) (*domain.Jar, error) {
	return s.jarRepo.GetByID(ctx, id)
}

func (s *Service) UpdateJarName(
	ctx context.Context,
	id string,
	name string,
) error {
	jar, err := s.jarRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	jar.UpdateName(name)
	return s.jarRepo.Update(ctx, jar)
}

func (s *Service) GetJarsBySessionID(
	ctx context.Context,
	sessionID string,
	pointer *string,
	limit int,
) (*ports.PaginatedResult[*domain.Jar], error) {
	return s.jarRepo.GetBySessionID(ctx, sessionID, pointer, limit)
}
