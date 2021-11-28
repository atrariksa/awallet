package models

type User struct {
	ID       uint
	Username string `gorm:"index:idx_username,unique"`
	Balance  int64
}

type UserTotalOutgoing struct {
	ID     uint
	User   User
	UserID uint   `gorm:"index:idx_user_id"`
	Value  uint64 `gorm:"index:idx_value,sort:desc,type:btree"`
}

type UserTotalOutgoingResult struct {
	UserID   uint
	Username string
	Value    uint64
}
