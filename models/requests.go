package models

type CreateUserRequest struct {
	Username string `json:"username" valid:"required"`
}
