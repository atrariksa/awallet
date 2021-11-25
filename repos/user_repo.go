package repos

import (
	"encoding/json"
	"log"

	"github.com/atrariksa/awallet/models"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

const (
	DuplicateKey string = "1062"
)

type UserRepoRead struct {
	DBRead *gorm.DB
	Cache  ICache
}

type IUserRepoRead interface {
	Read(username string, user *models.User) error
}

func (ur *UserRepoRead) Read(username string, user *models.User) error {
	bUser, err := ur.Cache.Get(username)
	if err == redis.Nil {
		err = ur.DBRead.Debug().Where(&models.User{Username: username}).First(user).Error
		if err != nil {
			log.Println(err)
			return err
		}

		bUser, err = json.Marshal(&user)
		if err != nil {
			log.Println(err)
			return err
		}

		ur.Cache.Set(username, bUser)
		return nil
	}

	if bUser != nil {
		err = json.Unmarshal(bUser, &user)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
 */

type UserRepoWrite struct {
	DBWrite *gorm.DB
	Cache   ICache
}

type IUserRepoWrite interface {
	Create(user *models.User) error
}

func (ur *UserRepoWrite) Create(user *models.User) error {

	err := ur.DBWrite.Debug().Create(user).Error
	if err != nil {
		return err
	}

	bUser, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	ur.Cache.Set(user.Username, bUser)

	return nil
}
