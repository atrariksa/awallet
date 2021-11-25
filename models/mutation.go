package models

import "time"

type Mutation struct {
	ID        uint
	UserID    uint
	RefID     string
	Type      string
	Value     uint32
	CreatedAt time.Time
}
