package models

type CreateUserRequest struct {
	Username string `json:"username" valid:"required"`
}

type TopupBalanceRequest struct {
	Amount uint32 `json:"amount" valid:"required,range(0|9999999)~Invalid topup amount"`
}
