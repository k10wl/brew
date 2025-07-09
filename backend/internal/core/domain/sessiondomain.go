package domain

import (
	"time"
)

type Session struct {
	ID           string
	CreatedAt    time.Time
	LastAccessed time.Time
	ExpiresAt    *time.Time
	IsActive     bool
	ShareTokens  []ShareToken
}

type ShareToken struct {
	Token     string
	Scope     ShareScope
	CreatedAt time.Time
	ExpiresAt *time.Time
	IsActive  bool
}

type ShareScope string

const (
	ReadOnlyScope  ShareScope = "read-only"
	ReadWriteScope ShareScope = "read-write"
)
