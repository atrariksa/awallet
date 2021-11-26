package services

import (
	"log"

	"github.com/atrariksa/awallet/errs"
	"github.com/atrariksa/awallet/models"
	"github.com/atrariksa/awallet/repos"
	"gorm.io/gorm"
)

type UserBalanceService struct {
	UserRepoRead     repos.IUserRepoRead
	UserBalanceWrite repos.IUserBalanceRepoWrite
}

type IUserBalanceService interface {
	GetBalanceByUsername(username string) (balance int64, err error)
	TopupBalance(user models.User, amount uint32) error
}

func (us *UserBalanceService) GetBalanceByUsername(username string) (balance int64, err error) {

	user := &models.User{Username: username}
	err = us.UserRepoRead.GetUser(user)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
		return
	}

	if user.ID == 0 {
		log.Println("Alert, user not found but have valid token!")
		return 0, errs.ErrUnauthorized
	}

	balance = user.Balance
	return
}

func (us *UserBalanceService) TopupBalance(user models.User, amount uint32) (err error) {
	err = us.UserBalanceWrite.Topup(user, amount)
	if err != nil {
		err = errs.ErrInternalServer
		return
	}
	return
}
