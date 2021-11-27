package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/atrariksa/awallet/models"
	"github.com/atrariksa/awallet/services"
)

type TopTransactionsPerUserHandler struct {
	UserBalanceService services.IUserBalanceService
}

func (ttpu *TopTransactionsPerUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := ctx.Value("token").(*models.JwtClaims)
	user := models.User{
		ID:       claims.UserID,
		Username: claims.Username,
	}

	resp, err := ttpu.UserBalanceService.GetTopTransactionsPerUser(user)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		bResp, _ := json.Marshal(&resp)
		w.WriteHeader(200)
		w.Write(bResp)
	}
}
