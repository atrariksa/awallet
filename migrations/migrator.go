package migrations

import (
	"github.com/atrariksa/awallet/models"
	"gorm.io/gorm"
)

type Migrator struct {
	DB *gorm.DB
}

func (m *Migrator) MigrateUp() {
	m.DB.AutoMigrate(
		&models.User{},
		&models.Mutation{},
		&models.UserTotalOutgoing{},
	)
}
