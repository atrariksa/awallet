package repos

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/atrariksa/awallet/errs"
	"github.com/atrariksa/awallet/models"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

const (
	TopTrsPrefix = "top_trs_%v"
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

	updateTotalOutgoing := tx.Debug().Model(&models.UserTotalOutgoing{}).Where("user_id = ?", user.ID).UpdateColumn("value", gorm.Expr("value + ?", amount))
	if updateTotalOutgoing.Error != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	if updateTotalOutgoing.RowsAffected == 0 {
		err = errs.ErrInternalServer
		tx.Rollback()
		return
	}

	err = tx.Debug().Create(&mutationIncoming).Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	addBalance := tx.Debug().Model(&destUser).UpdateColumn("balance", gorm.Expr("balance + ?", amount))
	err = addBalance.Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	if addBalance.RowsAffected == 0 {
		err = errs.ErrDestinationUserNotFound
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
	ur.Cache.Del(destUser.Username)
	ur.Cache.Del(TopUser)
	ur.Cache.Del(fmt.Sprintf(TopTrsPrefix, user.Username))

	return nil
}

/*
 */

type UserBalanceRepoRead struct {
	DBRead *gorm.DB
	Cache  ICache
}

type IUserBalanceRepoRead interface {
	GetTopTransactionResult(user models.User) (data []models.TopTransactionResult, err error)
}

func (ubr *UserBalanceRepoRead) GetTopTransactionResult(user models.User) (data []models.TopTransactionResult, err error) {

	key := fmt.Sprintf(TopTrsPrefix, user.Username)
	bData, err := ubr.Cache.Get(key)
	if err == redis.Nil {
		userMutations := ubr.DBRead.Debug().
			Table("mutations as m2").
			Limit(10).
			Order("m2.value desc").
			Where("m2.user_id = ?", user.ID).
			Where("m2.mutation_type IN ?",
				[]string{string(models.INCOMING), string(models.OUTGOING)})

		userJoinMutation := ubr.DBRead.Debug().
			Table("users as u").
			Order("m2.value desc").
			Select("u.username, m2.ref_id , m2.mutation_type ,m2.value").
			Joins("join mutations m on m.user_id = u.id").
			Joins(
				"join (?) as m2 on m2.ref_id = m.ref_id "+
					"and m.user_id != m2.user_id",
				userMutations,
			).
			Scan(&data)

		err = userJoinMutation.Error
		if err != nil {
			log.Println(err)
			return
		}

		bData, err = json.Marshal(&data)
		if err != nil {
			log.Println(err)
			return
		}

		ubr.Cache.Set(key, bData)
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
