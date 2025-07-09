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

type JarRepository interface {
	Save(ctx context.Context, jar *domain.Jar) error
	GetByID(ctx context.Context, id string) (*domain.Jar, error)
	GetBySessionID(
		ctx context.Context,
		sessionID string,
		pointer *string,
		limit int,
	) (*PaginatedResult[*domain.Jar], error)
	Update(ctx context.Context, jar *domain.Jar) error
	Exists(ctx context.Context, id string) (bool, error)
}

type SessionRepository interface {
	Save(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id string) (*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id string) error
}
