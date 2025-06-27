package dao

import (
	"time"
)

type LoginSession struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID        uint64    `gorm:"index" json:"user_id"`
	WalletAddress string    `gorm:"type:varchar(42)" json:"wallet_address"`
	Token         string    `gorm:"type:varchar(255);unique" json:"token"`
	IPAddress     string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent     string    `gorm:"type:text" json:"user_agent"`
	ExpiresAt     time.Time `gorm:"index" json:"expires_at"`
	RevokedAt     time.Time `gorm:"index;null" json:"revoked_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewLoginSession() *LoginSession {
	return &LoginSession{}
}

func (dao *Dao) GetValidSessionByToken(token string) (*LoginSession, error) {
	var session LoginSession
	err := dao.DB.Where("token = ? AND expires_at > ? AND revoked_at IS NULL", token, time.Now()).First(&session).Error
	return &session, err
}

func (dao *Dao) CreateLoginSession(session *LoginSession) error {
	return dao.DB.Create(session).Error
}

func (dao *Dao) RevokeSessionByToken(token string) error {
	return dao.DB.Model(&LoginSession{}).Where("token = ?", token).Update("revoked_at", time.Now()).Error
}
