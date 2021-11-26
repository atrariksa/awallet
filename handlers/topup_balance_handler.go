package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/atrariksa/awallet/models"
	"github.com/atrariksa/awallet/services"
)

type TopupBalanceHandler struct {
	UserBalanceService services.IUserBalanceService
}

func (tbh *TopupBalanceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := ctx.Value("token").(*models.JwtClaims)

	user := models.User{
		ID:       claims.UserID,
		Username: claims.Username,
	}

	req, err := tbh.validateAndGetTopupPayload(r)
	if err != nil {
		tbh.errInvalidAmount(w)
		return
	}

	err = tbh.UserBalanceService.TopupBalance(user, req.Amount)
	if err != nil {
		tbh.errInternal(w, err.Error())
		return
	}
	w.WriteHeader(204)
	w.Write([]byte("Topup successful"))
}

func (tbh *TopupBalanceHandler) validateAndGetTopupPayload(r *http.Request) (req models.TopupBalanceRequest, err error) {
	bodyByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyByte, &req)
	if err != nil {
		return
	}
	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		return
	}
	return
}

func (tbh *TopupBalanceHandler) errInvalidAmount(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Write([]byte("Invalid topup amount"))
}

func (tbh *TopupBalanceHandler) errInternal(w http.ResponseWriter, message string) {
	w.WriteHeader(500)
	w.Write([]byte(message))
}
