package dao

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserWallet struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"index" json:"user_id"`
	WalletAddress string    `gorm:"type:varchar(42);unique" json:"wallet_address"`
	WalletType    int       `gorm:"type:tinyint;default:1" json:"wallet_type"`
	PublicKey     string    `gorm:"type:text" json:"public_key"`
	Nonce         string    `gorm:"type:varchar(36)" json:"nonce"`
	LastLoginAt   time.Time `gorm:"index" json:"last_login_at"`
	IsPrimary     bool      `gorm:"type:tinyint;default:0" json:"is_primary"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewUserWallet() *UserWallet {
	return &UserWallet{}
}

// TableName 设置表名
func (u *UserWallet) TableName() string {
	return "user_wallet"
}

// GetUserWalletByAddressAndNonce 根据钱包地址和nonce查询钱包记录
func (dao *Dao) GetUserWalletByAddressAndNonce(walletAddress, nonce string) (*UserWallet, error) {

	var wallet UserWallet
	err := dao.DB.Where("wallet_address = ? AND nonce = ?", walletAddress, nonce).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("无效的钱包地址或随机数")
		}
		return nil, errors.Wrap(err, "查询用户钱包失败")
	}
	return &wallet, nil
}

// CreateUserWallet 创建用户钱包记录
func (dao *Dao) CreateUserWallet(wallet *UserWallet) error {
	err := dao.DB.Create(wallet).Error

	if err != nil {
		return errors.Wrap(err, "创建用户钱包失败")
	}
	return nil
}

// 更新最后登录时间
func (dao *Dao) UpdateLastLoginAt(walletAddress string) error {
	err := dao.DB.Model(&UserWallet{}).
		Where("wallet_address = ?", walletAddress).
		Update("last_login_at", time.Now()).Error
	if err != nil {
		return errors.Wrap(err, "更新最后登录时间失败")
	}
	return nil
}

// GetUserWalletByAddress 根据钱包地址查询钱包记录
func (dao *Dao) GetUserWalletByAddress(walletAddress string) (*UserWallet, error) {
	var wallet UserWallet
	err := dao.DB.Where("wallet_address = ?", walletAddress).First(&wallet).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("无效的钱包地址")
		}
		return nil, errors.Wrap(err, "查询用户钱包失败")
	}
	return &wallet, nil
}

// 更新nonce
func (dao *Dao) UpdateNonce(walletAddress, nonce string) error {
	err := dao.DB.Model(&UserWallet{}).
		Where("wallet_address = ?", walletAddress).
		Update("nonce", nonce).Error
	if err != nil {
		return errors.Wrap(err, "更新nonce失败")
	}
	return nil
}
