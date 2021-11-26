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

type RegisterHandler struct {
	UserService  services.IUserService
	TokenService services.ITokenService
}

func (rh *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req, err := rh.validateAndGetCreateUserPayload(r)
	if err != nil {
		rh.errInvalidPayload(w)
		return
	}
	user, err := rh.UserService.CreateUser(req.Username)
	if err != nil {
		rh.errCreateUser(w, err.Error())
		return
	}
	token, err := rh.TokenService.CreateToken(user)
	if err != nil {
		rh.errInternal(w, err.Error())
		return
	}
	resp := models.CreateUserResponse{
		Token: token,
		UserDetails: models.UserDetails{
			ID:       user.ID,
			Username: user.Username,
		},
	}
	res, _ := json.Marshal(resp)
	w.WriteHeader(201)
	w.Write([]byte(res))
}

func (rh *RegisterHandler) validateAndGetCreateUserPayload(r *http.Request) (req models.CreateUserRequest, err error) {
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

func (rh *RegisterHandler) errInvalidPayload(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Write([]byte("Bad Request"))
}

func (rh *RegisterHandler) errCreateUser(w http.ResponseWriter, message string) {
	if message == errs.ErrUserAlreadyExists.Error() {
		w.WriteHeader(409)
		w.Write([]byte(message))
	} else {
		rh.errInternal(w, message)
	}
}

func (rh *RegisterHandler) errInternal(w http.ResponseWriter, message string) {
	w.WriteHeader(500)
	w.Write([]byte(message))
}
