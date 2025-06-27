package dao

import (
	"MetaFarmBackend/component/db"
	"time"
)

type UserAccount struct {
	ID               uint64     `gorm:"primaryKey;column:id"`
	UserAddress      string     `gorm:"column:user_address;type:varchar(42);uniqueIndex"`
	Username         string     `gorm:"column:username;type:varchar(50);index"`
	Email            string     `gorm:"column:email;type:varchar(100)"`
	RegistrationTime time.Time  `gorm:"column:registration_time"`
	LastLoginTime    *time.Time `gorm:"column:last_login_time"`
	TotalLandCount   int        `gorm:"column:total_land_count;default:0"`
	ActiveLandCount  int        `gorm:"column:active_land_count;default:0"`
	TotalNFTCount    int        `gorm:"column:total_nft_count;default:0"`
	MFGBalance       float64    `gorm:"column:mfg_balance;type:decimal(36,18);default:0.0"`
	LastClaimTime    *time.Time `gorm:"column:last_claim_time"`
	IsBanned         bool       `gorm:"column:is_banned;default:false"`
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

func GetUserAccountByAddress(address string) (*UserAccount, error) {
	var user UserAccount
	err := db.GetDB().Where("user_address = ?", address).First(&user).Error
	return &user, err
}

func CreateUserAccount(user *UserAccount) error {
	return db.GetDB().Create(user).Error
}

func UpdateUserAccount(user *UserAccount) error {
	return db.GetDB().Save(user).Error
}
