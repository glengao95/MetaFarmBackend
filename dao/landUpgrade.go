package dao

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// LandUpgrade 土地升级记录表结构体
type LandUpgrade struct {
	ID           uint64          `gorm:"primaryKey;column:id"`                        // 主键ID
	LandTokenID  string          `gorm:"column:land_token_id;index"`                  // 土地NFT TokenID
	OwnerAddress string          `gorm:"column:owner_address;type:varchar(42);index"` // 所有者钱包地址
	OldLevel     int8            `gorm:"column:old_level"`                            // 升级前等级
	NewLevel     int8            `gorm:"column:new_level;index"`                      // 升级后等级
	CostTokens   uint64          `gorm:"column:cost_tokens;"`                         // 消耗代币数量
	CostItems    json.RawMessage `gorm:"column:cost_items;type:json"`                 // 消耗道具列表(JSON格式)
	UpgradeTime  time.Time       `gorm:"column:upgrade_time"`                         // 升级时间
	CreateTime   time.Time       `gorm:"column:create_time"`                          // 创建时间
	UpdateTime   time.Time       `gorm:"column:update_time"`                          // 更新时间
}

func (LandUpgrade) TableName() string {
	return "land_upgrade"
}

func NewLandUpgrade(landTokenID string, ownerAddress string, oldLevel, newLevel int8, costTokens uint64, costItems map[string]int) *LandUpgrade {
	now := time.Now()
	itemsJSON, _ := json.Marshal(costItems)
	return &LandUpgrade{
		LandTokenID:  landTokenID,
		OwnerAddress: ownerAddress,
		OldLevel:     oldLevel,
		NewLevel:     newLevel,
		CostTokens:   costTokens,
		CostItems:    itemsJSON,
		UpgradeTime:  now,
		CreateTime:   now,
		UpdateTime:   now,
	}
}

func (dao *Dao) GetLandUpgradeHistory(ctx context.Context, landTokenID uint64) ([]*LandUpgrade, error) {
	var upgrades []*LandUpgrade
	err := dao.DB.WithContext(ctx).Where("land_token_id = ?", landTokenID).Order("upgrade_time DESC").Find(&upgrades).Error
	return upgrades, err
}

func (dao *Dao) CreateLandUpgrade(ctx context.Context, tx *gorm.DB, upgrade *LandUpgrade) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Create(upgrade).Error
}
