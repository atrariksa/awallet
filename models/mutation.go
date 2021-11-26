package models

import (
	"time"

	"github.com/atrariksa/awallet/utils"
)

type Mutation struct {
	ID           uint
	UserID       uint `gorm:"index:idx_user_id"`
	RefID        string
	MutationType MutationType
	Value        uint32
	CreatedAt    time.Time
}

type MutationType string

const (
	INCOMING MutationType = "INCOMING"
	TOPUP    MutationType = "TOPUP"
	OUTGOING MutationType = "OUTGOING"
)

func NewTopupMutation(user User, amount uint32) Mutation {
	mutation := Mutation{
		UserID:       user.ID,
		RefID:        utils.NewUUIDString(),
		MutationType: TOPUP,
		Value:        amount,
		CreatedAt:    utils.TimeNowUTC(),
	}
	return mutation
}
