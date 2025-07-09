package ports

import (
	"context"

	"brew/internal/core/domain"
)

type PaginatedResult[T any] struct {
	Items       []T
	TotalCount  int
	NextPointer *string
	HasMore     bool
}

type BrewRepository interface {
	Save(ctx context.Context, brew *domain.Brew) error
	GetByID(ctx context.Context, id string) (*domain.Brew, error)
	GetBySessionID(
		ctx context.Context,
		sessionID string,
		pointer *string,
		limit int,
	) (*PaginatedResult[*domain.Brew], error)
	Update(ctx context.Context, brew *domain.Brew) error
	Exists(ctx context.Context, id string) (bool, error)
}

type SessionRepository interface {
	Save(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id string) (*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id string) error
}
