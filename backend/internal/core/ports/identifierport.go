package ports

import "context"

type IdentifierGenerator interface {
	Generate(ctx context.Context, name string) (string, error)
	Validate(ctx context.Context, identifier string) error
}
