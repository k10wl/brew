package services

import (
	"context"
	"errors"
	"testing"

	"brew/internal/core/domain"
	"brew/internal/core/ports/mocks"
)

func TestJarService_CreateJar_Success(t *testing.T) {
	var receivedName string
	var receivedExistsID string
	var receivedSaveJar *domain.Jar

	jarRepo := &mocks.JarRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			receivedExistsID = id
			return false, nil
		},
		SaveFunc: func(ctx context.Context, jar *domain.Jar) error {
			receivedSaveJar = jar
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
			return "jar-123", nil
		},
	}

	service := NewJarService(jarRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-jar"
	sessionID := "session-123"

	jar, err := service.CreateJar(ctx, name, sessionID)

	if err != nil {
		t.Fatalf("CreateJar() error = %v, want nil", err)
	}
	if jar == nil {
		t.Fatal("CreateJar() returned nil jar")
	}
	if jar.ID != "jar-123" {
		t.Fatalf("CreateJar() jar.ID = %v, want jar-123", jar.ID)
	}
	if jar.Name != name {
		t.Fatalf("CreateJar() jar.Name = %v, want %v", jar.Name, name)
	}

	if receivedName != name {
		t.Fatalf("Generate called with name = %v, want %v", receivedName, name)
	}
	if receivedExistsID != "jar-123" {
		t.Fatalf("Exists called with id = %v, want jar-123", receivedExistsID)
	}
	if receivedSaveJar == nil {
		t.Fatal("Save was not called")
	}
	if receivedSaveJar.ID != "jar-123" {
		t.Fatalf(
			"Save called with jar.ID = %v, want jar-123",
			receivedSaveJar.ID,
		)
	}
	if receivedSaveJar.Name != name {
		t.Fatalf(
			"Save called with jar.Name = %v, want %v",
			receivedSaveJar.Name,
			name,
		)
	}
}

func TestJarService_CreateJar_IdentifierGenerationError(t *testing.T) {
	jarRepo := &mocks.JarRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		SaveFunc: func(ctx context.Context, jar *domain.Jar) error {
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

	service := NewJarService(jarRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-jar"
	sessionID := "session-123"

	jar, err := service.CreateJar(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateJar() error = nil, want error")
	}
	if err.Error() != "identifier generation failed" {
		t.Fatalf("CreateJar() error = %v, want identifier generation failed", err.Error())
	}
	if jar != nil {
		t.Fatal("CreateJar() returned jar, want nil")
	}
}

func TestJarService_CreateJar_JarAlreadyExists(t *testing.T) {
	jarRepo := &mocks.JarRepository{
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
			return "jar-123", nil
		},
	}

	service := NewJarService(jarRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-jar"
	sessionID := "session-123"

	jar, err := service.CreateJar(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateJar() error = nil, want error")
	}
	if jar != nil {
		t.Fatal("CreateJar() returned jar, want nil")
	}
	expectedError := "jar with id jar-123 already exists"
	if err.Error() != expectedError {
		t.Fatalf("CreateJar() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestJarService_CreateJar_ExistsCheckError(t *testing.T) {
	jarRepo := &mocks.JarRepository{
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
			return "jar-123", nil
		},
	}

	service := NewJarService(jarRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-jar"
	sessionID := "session-123"

	jar, err := service.CreateJar(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateJar() error = nil, want error")
	}
	if err.Error() != "exists check failed" {
		t.Fatalf("CreateJar() error = %v, want exists check failed", err.Error())
	}
	if jar != nil {
		t.Fatal("CreateJar() returned jar, want nil")
	}
}

func TestJarService_CreateJar_SaveError(t *testing.T) {
	jarRepo := &mocks.JarRepository{
		ExistsFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		SaveFunc: func(ctx context.Context, jar *domain.Jar) error {
			return errors.New("save failed")
		},
	}
	sessionRepo := &mocks.SessionRepository{}
	identifierGen := &mocks.IdentifierGenerator{
		GenerateFunc: func(
			ctx context.Context,
			name string,
		) (string, error) {
			return "jar-123", nil
		},
	}

	service := NewJarService(jarRepo, sessionRepo, identifierGen)

	ctx := context.Background()
	name := "test-jar"
	sessionID := "session-123"

	jar, err := service.CreateJar(ctx, name, sessionID)

	if err == nil {
		t.Fatal("CreateJar() error = nil, want error")
	}
	if err.Error() != "save failed" {
		t.Fatalf("CreateJar() error = %v, want save failed", err.Error())
	}
	if jar != nil {
		t.Fatal("CreateJar() returned jar, want nil")
	}
}
