package model

import "time"

type Workspace struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}