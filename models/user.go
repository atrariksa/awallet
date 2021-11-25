package models

type User struct {
	ID       uint
	Username string `gorm:"index:idx_username,unique"`
}
