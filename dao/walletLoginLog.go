package dao

import (
	"context"
	"time"
)

type WalletLoginLog struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	WalletAddress string    `gorm:"type:varchar(42);index" json:"wallet_address"`
	IPAddress     string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent     string    `gorm:"type:text" json:"user_agent"`
	LoginTime     time.Time `gorm:"index" json:"login_time"`
	Status        int       `gorm:"type:tinyint;default:0" json:"status"`
	ErrorMessage  string    `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
}

// 记录登录日志
func (dao *Dao) RecordLoginLog(ctx context.Context, walletAddress, ipAddress, userAgent string, success bool, errorMsg string) {
	log := WalletLoginLog{
		WalletAddress: walletAddress,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		LoginTime:     time.Now(),
		Status:        boolToInt(success),
		ErrorMessage:  errorMsg,
		CreatedAt:     time.Now(),
	}

	// 异步记录日志
	go func() {
		dao.DB.Create(&log)
	}()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
