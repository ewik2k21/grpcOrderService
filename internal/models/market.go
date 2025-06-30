package models

import (
	"github.com/google/uuid"
	"time"
)

type Market struct {
	ID        uuid.UUID
	Name      string
	Enabled   bool
	DeletedAt *time.Time
}
