package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/atrariksa/awallet/services"
)

type ListTopUserHandler struct {
	UserService services.IUserService
}

func (ltu *ListTopUserHandler) Handle(w http.ResponseWriter, r *http.Request) {

	resp, err := ltu.UserService.GetListTopUser()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		bResp, _ := json.Marshal(&resp)
		w.WriteHeader(200)
		w.Write(bResp)
	}
}
