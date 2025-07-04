package dao

import (
	"context"
	"time"
)

// User 用户基本信息表结构体
type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`      // 用户ID
	Username  string    `gorm:"type:varchar(50);unique" json:"username"` // 用户名
	CreatedAt time.Time `gorm:"index" json:"created_at"`                 // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                              // 更新时间
}

func (u *User) TableName() string {
	return "user"
}

// CreateUser 创建用户
func (dao *Dao) CreateUser(ctx context.Context, user *User) error {
	return dao.DB.WithContext(ctx).Create(user).Error
}
