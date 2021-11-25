package models

type CreateUserResponse struct {
	Token       string      `json:"token"`
	UserDetails UserDetails `json:"user_details"`
}

type UserDetails struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
