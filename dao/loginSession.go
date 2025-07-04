package dao

import (
	"context"
	"time"
)

// LoginSession 用户登录会话表结构体
type LoginSession struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)" json:"id"`  // 会话ID
	UserID        uint64    `gorm:"index" json:"user_id"`                   // 用户ID
	WalletAddress string    `gorm:"type:varchar(42)" json:"wallet_address"` // 钱包地址
	Token         string    `gorm:"type:varchar(255);unique" json:"token"`  // 会话令牌
	IPAddress     string    `gorm:"type:varchar(45)" json:"ip_address"`     // IP地址
	UserAgent     string    `gorm:"type:text" json:"user_agent"`            // 用户代理
	ExpiresAt     time.Time `gorm:"index" json:"expires_at"`                // 过期时间
	RevokedAt     time.Time `gorm:"index;null" json:"revoked_at"`           // 吊销时间
	CreatedAt     time.Time `json:"created_at"`                             // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                             // 更新时间
}

func NewLoginSession() *LoginSession {
	return &LoginSession{}
}

func (dao *Dao) GetValidSessionByToken(ctx context.Context, token string) (*LoginSession, error) {
	var session LoginSession
	err := dao.DB.WithContext(ctx).Where("token = ? AND expires_at > ? AND revoked_at IS NULL", token, time.Now()).First(&session).Error
	return &session, err
}

func (dao *Dao) CreateLoginSession(ctx context.Context, session *LoginSession) error {
	return dao.DB.WithContext(ctx).Create(session).Error
}

func (dao *Dao) RevokeSessionByToken(ctx context.Context, token string) error {
	return dao.DB.WithContext(ctx).Model(&LoginSession{}).Where("token = ?", token).Update("revoked_at", time.Now()).Error
}
