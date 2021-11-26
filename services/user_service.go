package services

import (
	"log"
	"strings"

	"github.com/atrariksa/awallet/errs"
	"github.com/atrariksa/awallet/models"
	"github.com/atrariksa/awallet/repos"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepoRead  repos.IUserRepoRead
	UserRepoWrite repos.IUserRepoWrite
}

type IUserService interface {
	CreateUser(username string) (user models.User, err error)
}

func (us *UserService) CreateUser(username string) (user models.User, err error) {

	user.Username = username
	err = us.UserRepoRead.GetUser(&user)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
		return
	}

	if user.ID != 0 {
		return models.User{}, errs.ErrUserAlreadyExists
	}

	// user = models.User{Username: username}
	err = us.UserRepoWrite.Create(&user)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), repos.DuplicateKey) {
			return models.User{}, errs.ErrUserAlreadyExists
		}
		return models.User{}, errs.ErrInternalServer
	}
	return
}
