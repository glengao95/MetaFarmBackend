package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// UserWallet 用户钱包信息表结构体
type UserWallet struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`            // 主键ID
	UserID        uint64    `gorm:"index" json:"user_id"`                          // 用户ID
	WalletAddress string    `gorm:"type:varchar(42);unique" json:"wallet_address"` // 钱包地址
	WalletType    int       `gorm:"type:tinyint;default:1" json:"wallet_type"`     // 钱包类型(1:以太坊, 2:Polygon等)
	PublicKey     string    `gorm:"type:text" json:"public_key"`                   // 公钥
	Nonce         string    `gorm:"type:varchar(36)" json:"nonce"`                 // 随机数
	LastLoginAt   time.Time `gorm:"index" json:"last_login_at"`                    // 最后登录时间
	IsPrimary     bool      `gorm:"type:tinyint;default:0" json:"is_primary"`      // 是否主钱包
	CreatedAt     time.Time `json:"created_at"`                                    // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                                    // 更新时间
}

func NewUserWallet() *UserWallet {
	return &UserWallet{}
}

// TableName 设置表名
func (u *UserWallet) TableName() string {
	return "user_wallet"
}

// GetUserWalletByAddressAndNonce 根据钱包地址和nonce查询钱包记录
func (dao *Dao) GetUserWalletByAddressAndNonce(ctx context.Context, walletAddress, nonce string) (*UserWallet, error) {
	var wallet UserWallet
	err := dao.DB.WithContext(ctx).Where("wallet_address = ? AND nonce = ?", walletAddress, nonce).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("无效的钱包地址或随机数")
		}
		return nil, errors.Wrap(err, "查询用户钱包失败")
	}
	return &wallet, nil
}

// CreateUserWallet 创建用户钱包记录
func (dao *Dao) CreateUserWallet(ctx context.Context, wallet *UserWallet) error {
	err := dao.DB.WithContext(ctx).Create(wallet).Error

	if err != nil {
		return errors.Wrap(err, "创建用户钱包失败")
	}
	return nil
}

// 更新最后登录时间
func (dao *Dao) UpdateLastLoginAt(ctx context.Context, walletAddress string) error {
	return dao.DB.WithContext(ctx).Model(&UserWallet{}).
		Where("wallet_address = ?", walletAddress).
		Update("last_login_at", time.Now()).Error
}

// GetUserWalletByAddress 根据钱包地址查询钱包记录
func (dao *Dao) GetUserWalletByAddress(ctx context.Context, walletAddress string) (*UserWallet, error) {
	var wallet UserWallet
	err := dao.DB.WithContext(ctx).Where("wallet_address = ?", walletAddress).First(&wallet).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("无效的钱包地址")
		}
		return nil, errors.Wrap(err, "查询用户钱包失败")
	}
	return &wallet, nil
}

// 更新nonce
func (dao *Dao) UpdateNonce(ctx context.Context, walletAddress, nonce string) error {
	return dao.DB.WithContext(ctx).Model(&UserWallet{}).
		Where("wallet_address = ?", walletAddress).
		Update("nonce", nonce).Error
}
