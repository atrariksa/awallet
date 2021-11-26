package repos

import (
	"log"

	"github.com/atrariksa/awallet/models"
	"gorm.io/gorm"
)

type UserBalanceRepoWrite struct {
	DBWrite *gorm.DB
	Cache   ICache
}

type IUserBalanceRepoWrite interface {
	Topup(user models.User, amount uint32) error
}

func (ur *UserBalanceRepoWrite) Topup(user models.User, amount uint32) (err error) {

	mutation := models.NewTopupMutation(user, amount)
	tx := ur.DBWrite.Begin()

	err = tx.Debug().Create(&mutation).Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	err = tx.Debug().Model(&user).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	err = tx.Debug().Commit().Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	ur.Cache.Del(user.Username)

	return nil
}
