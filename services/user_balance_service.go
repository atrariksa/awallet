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
	UserBalanceRead  repos.IUserBalanceRepoRead
}

type IUserBalanceService interface {
	GetBalanceByUsername(username string) (balance int64, err error)
	TopupBalance(user models.User, amount uint32) error
	Transfer(user models.User, amount uint32, destUsername string) error
	GetTopTransactionsPerUser(user models.User) (data []models.TopTransactionsPerUser, err error)
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

func (us *UserBalanceService) Transfer(user models.User, amount uint32, destUsername string) (err error) {
	destUser := models.User{Username: destUsername}
	err = us.UserRepoRead.GetUser(&destUser)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
		return
	}

	if destUser.ID == 0 {
		return errs.ErrDestinationUserNotFound
	}

	err = us.UserBalanceWrite.Transfer(user, amount, destUser)
	if err != nil {
		return
	}
	return
}

func (us *UserBalanceService) GetTopTransactionsPerUser(user models.User) (data []models.TopTransactionsPerUser, err error) {
	dataTopTrs, err := us.UserBalanceRead.GetTopTransactionResult(user)
	for _, v := range dataTopTrs {
		var amount int64 = int64(v.Value)
		if v.MutationType == models.OUTGOING {
			amount = amount * -1
		}
		data = append(data, models.TopTransactionsPerUser{
			Username: v.Username,
			Amount:   amount,
		})
	}
	return
}
