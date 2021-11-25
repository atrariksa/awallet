package drivers

import (
	"fmt"

	"github.com/atrariksa/awallet/configs"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBClientRead(cfg *configs.Config) *gorm.DB {
	dbCfg := cfg.DB.MySQL.Read
	return buildDBclient(
		dbCfg.Name,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.AdditionalParameters,
	)
}

func NewDBClientWrite(cfg *configs.Config) *gorm.DB {
	dbCfg := cfg.DB.MySQL.Write
	return buildDBclient(
		dbCfg.Name,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.AdditionalParameters,
	)
}

func buildDBclient(name, username, password, host, port, additionalParams string) *gorm.DB {
	prepStr := "%v:%v@tcp(%v:%v)/%v%v"
	dbPar := fmt.Sprintf(
		prepStr,
		username,
		password,
		host,
		port,
		name,
		additionalParams,
	)
	db, err := gorm.Open(mysql.Open(dbPar), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
