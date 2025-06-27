package dao

import (
	"database/sql"
	"time"

	"MetaFarmBackend/component/db"

	"gorm.io/gorm"
)

// UserAssetsSummary 用户资产汇总表结构体
type UserAssetsSummary struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserAddress    string          `gorm:"column:user_address;size:42;not null;uniqueIndex:idx_user_address" json:"user_address"`
	TotalLandValue sql.NullFloat64 `gorm:"column:total_land_value;type:decimal(36,18);default:0.0;index:idx_total_land_value" json:"total_land_value"`
	TotalItemValue sql.NullFloat64 `gorm:"column:total_item_value;type:decimal(36,18);default:0.0;index:idx_total_item_value" json:"total_item_value"`
	TotalNFTCount  int             `gorm:"column:total_nft_count;default:0" json:"total_nft_count"`
	MFGBalance     sql.NullFloat64 `gorm:"column:mfg_balance;type:decimal(36,18);default:0.0" json:"mfg_balance"`
	LastUpdated    time.Time       `gorm:"column:last_updated;not null;autoUpdateTime" json:"last_updated"`
}

// TableName 设置表名
func (u *UserAssetsSummary) TableName() string {
	return "user_assets_summary"
}

// NewUserAssetsSummary 创建新的用户资产汇总实例
func NewUserAssetsSummary(userAddress string) *UserAssetsSummary {
	now := time.Now()
	return &UserAssetsSummary{
		UserAddress: userAddress,
		TotalLandValue: sql.NullFloat64{
			Float64: 0.0,
			Valid:   true,
		},
		TotalItemValue: sql.NullFloat64{
			Float64: 0.0,
			Valid:   true,
		},
		TotalNFTCount: 0,
		MFGBalance: sql.NullFloat64{
			Float64: 0.0,
			Valid:   true,
		},
		LastUpdated: now,
	}
}

// GetUserAssetsSummaryByAddress 根据用户地址获取资产汇总信息
func GetUserAssetsSummaryByAddress(userAddress string) (*UserAssetsSummary, error) {
	var summary UserAssetsSummary
	err := db.GetDB().Where("user_address = ?", userAddress).First(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// CreateUserAssetsSummary 创建用户资产汇总记录
func (u *UserAssetsSummary) CreateUserAssetsSummary() error {
	return db.GetDB().Create(u).Error
}

// UpdateUserAssetsSummary 更新用户资产汇总记录
func (u *UserAssetsSummary) UpdateUserAssetsSummary() error {
	return db.GetDB().Save(u).Error
}

// IncrementTotalNFTCount 增加NFT总数
func (u *UserAssetsSummary) IncrementTotalNFTCount(count int) error {
	return db.GetDB().Model(u).Update("total_nft_count", gorm.Expr("total_nft_count + ?", count)).Error
}

// UpdateLandValue 更新土地总价值
func (u *UserAssetsSummary) UpdateLandValue(value float64) error {
	return db.GetDB().Model(u).Update("total_land_value", value).Error
}

// UpdateItemValue 更新道具总价值
func (u *UserAssetsSummary) UpdateItemValue(value float64) error {
	return db.GetDB().Model(u).Update("total_item_value", value).Error
}
