package repos

import (
	"log"

	"github.com/atrariksa/awallet/errs"
	"github.com/atrariksa/awallet/models"
	"gorm.io/gorm"
)

type UserBalanceRepoWrite struct {
	DBWrite *gorm.DB
	Cache   ICache
}

type IUserBalanceRepoWrite interface {
	Topup(user models.User, amount uint32) error
	Transfer(user models.User, amount uint32, destUser models.User) error
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

func (ur *UserBalanceRepoWrite) Transfer(user models.User, amount uint32, destUser models.User) (err error) {

	mutationOutgoing, mutationIncoming := models.NewTransferBalanceMutation(user, amount, destUser)
	tx := ur.DBWrite.Begin()

	err = tx.Debug().Create(&mutationOutgoing).Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	deductBalance := tx.Debug().Model(&user).Where("balance >= ?", amount).UpdateColumn("balance", gorm.Expr("balance - ?", amount))
	if deductBalance.Error != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	if deductBalance.RowsAffected == 0 {
		err = errs.ErrInsufficientBalance
		tx.Rollback()
		return
	}

	err = tx.Debug().Create(&mutationIncoming).Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	err = tx.Debug().Model(&destUser).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
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
