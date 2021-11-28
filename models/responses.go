package models

type CreateUserResponse struct {
	Token       string      `json:"token"`
	UserDetails UserDetails `json:"user_details"`
}

type UserDetails struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
}

type ReadBalanceResponse struct {
	Balance int64 `json:"balance"`
}

type TopTransactionsPerUser struct {
	Username string `json:"username"`
	Amount   int64  `json:"amount"`
}

type TopUserResponse struct {
	Username        string `json:"username"`
	TransactedValue uint64 `json:"transacted_value"`
}
