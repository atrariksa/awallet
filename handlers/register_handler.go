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

func (uh *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req, err := uh.validateAndGetCreateUserPayload(r)
	if err != nil {
		uh.errInvalidPayload(w)
		return
	}
	user, err := uh.UserService.CreateUser(req.Username)
	if err != nil {
		uh.errCreateUser(w, err.Error())
		return
	}
	token, err := uh.TokenService.CreateToken(user.ID)
	if err != nil {
		uh.errInternal(w, err.Error())
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

func (uh *RegisterHandler) validateAndGetCreateUserPayload(r *http.Request) (req models.CreateUserRequest, err error) {
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

func (uh *RegisterHandler) errInvalidPayload(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Write([]byte("Bad Request"))
}

func (uh *RegisterHandler) errCreateUser(w http.ResponseWriter, message string) {
	if message == errs.ErrUserAlreadyExists.Error() {
		w.WriteHeader(409)
		w.Write([]byte(message))
	} else {
		uh.errInternal(w, message)
	}
}

func (uh *RegisterHandler) errInternal(w http.ResponseWriter, message string) {
	w.WriteHeader(500)
	w.Write([]byte(message))
}
