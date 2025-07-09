package domain

import (
	"time"
)

type Brew struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBrew(id string, name string) *Brew {
	now := time.Now()
	return &Brew{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *Brew) UpdateName(name string) {
	b.Name = name
	b.UpdatedAt = time.Now()
}
