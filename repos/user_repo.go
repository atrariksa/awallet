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
	TopUser      string = "top_user"
)

type UserRepoRead struct {
	DBRead *gorm.DB
	Cache  ICache
}

type IUserRepoRead interface {
	GetUser(user *models.User) error
	GetListTopUser() (data []models.UserTotalOutgoingResult, err error)
}

func (ur *UserRepoRead) GetUser(user *models.User) error {
	bUser, err := ur.Cache.Get(user.Username)
	if err == redis.Nil {
		err = ur.DBRead.Debug().Where(user).First(user).Error
		if err != nil {
			log.Println(err)
			return err
		}

		bUser, err = json.Marshal(&user)
		if err != nil {
			log.Println(err)
			return err
		}

		ur.Cache.Set(user.Username, bUser)
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

func (ur *UserRepoRead) GetListTopUser() (data []models.UserTotalOutgoingResult, err error) {

	bData, err := ur.Cache.Get(TopUser)
	if err == redis.Nil {

		userOutgoing := ur.DBRead.Debug().
			Table("user_total_outgoings as utg").
			Limit(10).
			Order("value desc").
			Select("utg.*, u.*").
			Joins("join users u on utg.user_id = u.id").
			Scan(&data)

		err = userOutgoing.Error
		if err != nil {
			log.Println(err)
			return
		}

		bData, err = json.Marshal(&data)
		if err != nil {
			log.Println(err)
			return
		}

		ur.Cache.Set(TopUser, bData)
		return
	}

	if bData != nil {
		err = json.Unmarshal(bData, &data)
		if err != nil {
			return
		}
	}

	return
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

	tx := ur.DBWrite.Debug().Begin()

	err := tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Create(&models.UserTotalOutgoing{UserID: user.ID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	bUser, err := json.Marshal(&user)
	if err != nil {
		tx.Rollback()
		return err
	}

	ur.Cache.Set(user.Username, bUser)

	return nil
}
