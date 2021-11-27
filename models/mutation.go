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

type TopTransactionResult struct {
	Username     string
	RefID        string
	MutationType MutationType
	Value        uint32
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

func NewOutgoingMutation(user User, amount uint32) Mutation {
	mutation := Mutation{
		UserID:       user.ID,
		RefID:        utils.NewUUIDString(),
		MutationType: OUTGOING,
		Value:        amount,
		CreatedAt:    utils.TimeNowUTC(),
	}
	return mutation
}

func NewIncomingMutation(user User, amount uint32) Mutation {
	mutation := Mutation{
		UserID:       user.ID,
		RefID:        utils.NewUUIDString(),
		MutationType: INCOMING,
		Value:        amount,
		CreatedAt:    utils.TimeNowUTC(),
	}
	return mutation
}

func NewTransferBalanceMutation(user User, amount uint32, destUser User) (outgoing Mutation, incoming Mutation) {
	outgoing = NewOutgoingMutation(user, amount)
	incoming = NewIncomingMutation(destUser, amount)
	incoming.RefID = outgoing.RefID
	incoming.CreatedAt = outgoing.CreatedAt
	return
}
