package dao

import (
	"context"
	"time"
)

// UserAccount 用户账户信息表结构体
type UserAccount struct {
	ID               uint64     `gorm:"primaryKey;column:id"`                               // 主键ID
	UserAddress      string     `gorm:"column:user_address;type:varchar(42);uniqueIndex"`   // 用户钱包地址
	Username         string     `gorm:"column:username;type:varchar(50);index"`             // 用户名
	Email            string     `gorm:"column:email;type:varchar(100)"`                     // 电子邮箱
	RegistrationTime time.Time  `gorm:"column:registration_time"`                           // 注册时间
	LastLoginTime    *time.Time `gorm:"column:last_login_time"`                             // 最后登录时间
	TotalLandCount   int        `gorm:"column:total_land_count;default:0"`                  // 土地总数
	ActiveLandCount  int        `gorm:"column:active_land_count;default:0"`                 // 活跃土地数
	TotalNFTCount    int        `gorm:"column:total_nft_count;default:0"`                   // NFT总数
	MFGBalance       float64    `gorm:"column:mfg_balance;type:decimal(36,18);default:0.0"` // MFG代币余额
	LastClaimTime    *time.Time `gorm:"column:last_claim_time"`                             // 最后领取时间
	IsBanned         bool       `gorm:"column:is_banned;default:false"`                     // 是否封禁
}

func (u *UserAccount) TableName() string {
	return "user_account"
}

func NewUserAccount(userAddress string) *UserAccount {
	return &UserAccount{
		UserAddress:      userAddress,
		RegistrationTime: time.Now(),
	}
}

func (dao *Dao) GetUserAccountByAddress(ctx context.Context, address string) (*UserAccount, error) {
	var user UserAccount
	err := dao.DB.WithContext(ctx).Where("user_address = ?", address).First(&user).Error
	return &user, err
}

func (dao *Dao) CreateUserAccount(ctx context.Context, user *UserAccount) error {
	return dao.DB.WithContext(ctx).Create(user).Error
}

func (dao *Dao) UpdateUserAccount(ctx context.Context, user *UserAccount) error {
	return dao.DB.WithContext(ctx).Save(user).Error
}
