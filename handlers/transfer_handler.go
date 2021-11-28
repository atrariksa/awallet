package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/atrariksa/awallet/errs"
	"github.com/atrariksa/awallet/models"
	"github.com/atrariksa/awallet/services"
)

type TransferHandler struct {
	UserBalanceService services.IUserBalanceService
}

func (tbh *TransferHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := ctx.Value("token").(*models.JwtClaims)

	user := models.User{
		ID:       claims.UserID,
		Username: claims.Username,
	}

	req, err := tbh.validateAndGetTransferPayload(r)
	if err != nil {
		tbh.errBadRequest(w, err.Error())
		return
	}

	err = tbh.UserBalanceService.Transfer(user, req.Amount, req.ToUsername)
	if err != nil {
		if err.Error() == errs.ErrInsufficientBalance.Error() {
			tbh.errBadRequest(w, err.Error())
			return
		}
		if err.Error() == errs.ErrDestinationUserNotFound.Error() {
			tbh.errUserNotFound(w, err.Error())
			return
		}
		tbh.errInternal(w, err.Error())
		return
	}
	w.WriteHeader(204)
}

func (tbh *TransferHandler) validateAndGetTransferPayload(r *http.Request) (req models.TransferRequest, err error) {
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

func (tbh *TransferHandler) errBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(400)
	w.Write([]byte(message))
}

func (tbh *TransferHandler) errUserNotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(404)
	w.Write([]byte(message))
}

func (tbh *TransferHandler) errInternal(w http.ResponseWriter, message string) {
	w.WriteHeader(500)
	w.Write([]byte(message))
}
