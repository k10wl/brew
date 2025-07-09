package domain

import (
	"time"
)

type Jar struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewJar(id string, name string) *Jar {
	now := time.Now()
	return &Jar{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (j *Jar) UpdateName(name string) {
	j.Name = name
	j.UpdatedAt = time.Now()
}
