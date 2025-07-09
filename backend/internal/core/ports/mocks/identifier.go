package mocks

import (
	"context"

	"brew/internal/core/ports"
)

var _ ports.IdentifierGenerator = (*IdentifierGenerator)(nil)

type IdentifierGenerator struct {
	GenerateFunc func(ctx context.Context, name string) (string, error)
	ValidateFunc func(ctx context.Context, identifier string) error
}

func (m *IdentifierGenerator) Generate(ctx context.Context, name string) (string, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, name)
	}
	return "", nil
}

func (m *IdentifierGenerator) Validate(ctx context.Context, identifier string) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(ctx, identifier)
	}
	return nil
}
