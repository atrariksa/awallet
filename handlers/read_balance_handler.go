package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/atrariksa/awallet/models"
	"github.com/atrariksa/awallet/services"
)

type ReadBalanceHandler struct {
	UserBalanceService services.IUserBalanceService
}

func (rbh *ReadBalanceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := ctx.Value("token").(*models.JwtClaims)
	balance, err := rbh.UserBalanceService.GetBalanceByUsername(claims.Username)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		resp := models.ReadBalanceResponse{Balance: balance}
		bResp, _ := json.Marshal(&resp)
		w.WriteHeader(200)
		w.Write(bResp)
	}
}
