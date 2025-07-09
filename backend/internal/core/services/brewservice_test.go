package services

import (
	"context"
	"errors"
	"testing"

	"brew/internal/core/domain"
	"brew/internal/core/ports/mocks"
)

func TestBrewService_CreateBrew_Success(t *testing.T) {
	var receivedName string
	var receivedExistsID string
	var receivedSaveBrew *domain.Brew

	brewRepo := &mocks.BrewRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			receivedExistsID = id
			return false, nil
		},
		SaveFunc: func(ctx context.Context, brew *domain.Brew) error {
			receivedSaveBrew = brew
			return nil
		},
	}
	sessionRepo := &mocks.SessionRepository{}
	identifierGen := &mocks.IdentifierGenerator{
		GenerateFunc: func(
			ctx context.Context,
			name string,
		) (string, error) {
			receivedName = name
			return "brew-123", nil
		},
	}

	service := NewBrewService(brewRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-brew"
	sessionID := "session-123"

	brew, err := service.CreateBrew(ctx, name, sessionID)

	if err != nil {
		t.Fatalf("CreateBrew() error = %v, want nil", err)
	}
	if brew == nil {
		t.Fatal("CreateBrew() returned nil brew")
	}
	if brew.ID != "brew-123" {
		t.Fatalf("CreateBrew() brew.ID = %v, want brew-123", brew.ID)
	}
	if brew.Name != name {
		t.Fatalf("CreateBrew() brew.Name = %v, want %v", brew.Name, name)
	}

	if receivedName != name {
		t.Fatalf("Generate called with name = %v, want %v", receivedName, name)
	}
	if receivedExistsID != "brew-123" {
		t.Fatalf("Exists called with id = %v, want brew-123", receivedExistsID)
	}
	if receivedSaveBrew == nil {
		t.Fatal("Save was not called")
	}
	if receivedSaveBrew.ID != "brew-123" {
		t.Fatalf(
			"Save called with brew.ID = %v, want brew-123",
			receivedSaveBrew.ID,
		)
	}
	if receivedSaveBrew.Name != name {
		t.Fatalf(
			"Save called with brew.Name = %v, want %v",
			receivedSaveBrew.Name,
			name,
		)
	}
}

func TestBrewService_CreateBrew_IdentifierGenerationError(t *testing.T) {
	brewRepo := &mocks.BrewRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		SaveFunc: func(ctx context.Context, brew *domain.Brew) error {
			return nil
		},
	}
	sessionRepo := &mocks.SessionRepository{}
	identifierGen := &mocks.IdentifierGenerator{
		GenerateFunc: func(
			ctx context.Context,
			name string,
		) (string, error) {
			return "", errors.New("identifier generation failed")
		},
	}

	service := NewBrewService(brewRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-brew"
	sessionID := "session-123"

	brew, err := service.CreateBrew(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateBrew() error = nil, want error")
	}
	if err.Error() != "identifier generation failed" {
		t.Fatalf("CreateBrew() error = %v, want identifier generation failed", err.Error())
	}
	if brew != nil {
		t.Fatal("CreateBrew() returned brew, want nil")
	}
}

func TestBrewService_CreateBrew_BrewAlreadyExists(t *testing.T) {
	brewRepo := &mocks.BrewRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
	}
	sessionRepo := &mocks.SessionRepository{}
	identifierGen := &mocks.IdentifierGenerator{
		GenerateFunc: func(
			ctx context.Context,
			name string,
		) (string, error) {
			return "brew-123", nil
		},
	}

	service := NewBrewService(brewRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-brew"
	sessionID := "session-123"

	brew, err := service.CreateBrew(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateBrew() error = nil, want error")
	}
	if brew != nil {
		t.Fatal("CreateBrew() returned brew, want nil")
	}
	expectedError := "brew with id brew-123 already exists"
	if err.Error() != expectedError {
		t.Fatalf("CreateBrew() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestBrewService_CreateBrew_ExistsCheckError(t *testing.T) {
	brewRepo := &mocks.BrewRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			return false, errors.New("exists check failed")
		},
	}
	sessionRepo := &mocks.SessionRepository{}
	identifierGen := &mocks.IdentifierGenerator{
		GenerateFunc: func(
			ctx context.Context,
			name string,
		) (string, error) {
			return "brew-123", nil
		},
	}

	service := NewBrewService(brewRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-brew"
	sessionID := "session-123"

	brew, err := service.CreateBrew(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateBrew() error = nil, want error")
	}
	if err.Error() != "exists check failed" {
		t.Fatalf("CreateBrew() error = %v, want exists check failed", err.Error())
	}
	if brew != nil {
		t.Fatal("CreateBrew() returned brew, want nil")
	}
}

func TestBrewService_CreateBrew_SaveError(t *testing.T) {
	brewRepo := &mocks.BrewRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		SaveFunc: func(ctx context.Context, brew *domain.Brew) error {
			return errors.New("save failed")
		},
	}
	sessionRepo := &mocks.SessionRepository{}
	identifierGen := &mocks.IdentifierGenerator{
		GenerateFunc: func(
			ctx context.Context,
			name string,
		) (string, error) {
			return "brew-123", nil
		},
	}

	service := NewBrewService(brewRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-brew"
	sessionID := "session-123"

	brew, err := service.CreateBrew(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateBrew() error = nil, want error")
	}
	if err.Error() != "save failed" {
		t.Fatalf("CreateBrew() error = %v, want save failed", err.Error())
	}
	if brew != nil {
		t.Fatal("CreateBrew() returned brew, want nil")
	}
}
