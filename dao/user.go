package dao

import (
	"time"
)

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);unique" json:"username"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) TableName() string {
	return "user"
}

// CreateUser 创建用户
func (dao *Dao) CreateUser(user *User) error {
	return dao.DB.Create(user).Error
}
