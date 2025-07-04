package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// UserItems 用户道具表结构体
type UserItems struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                   // 主键ID
	UserAddress   string    `gorm:"column:user_address;size:42;not null" json:"user_address"`                       // 用户钱包地址
	ItemTokenID   int64     `gorm:"column:item_token_id;not null" json:"item_token_id"`                             // 道具TokenID
	ItemType      int8      `gorm:"column:item_type;not null;index:idx_item_type" json:"item_type"`                 // 道具类型(1:肥料, 2:杀虫剂等)
	ItemName      string    `gorm:"column:item_name;size:50;not null" json:"item_name"`                             // 道具名称
	Rarity        int8      `gorm:"column:rarity;not null;default:1;index:idx_rarity" json:"rarity"`                // 稀有度(1:普通, 2:稀有, 3:史诗)
	Power         int       `gorm:"column:power;default:0" json:"power"`                                            // 道具效果值
	MaxUses       int       `gorm:"column:max_uses;default:0" json:"max_uses"`                                      // 最大使用次数
	RemainingUses int       `gorm:"column:remaining_uses;default:0;index:idx_remaining_uses" json:"remaining_uses"` // 剩余使用次数
	MetadataURI   string    `gorm:"column:metadata_uri;size:255" json:"metadata_uri"`                               // 元数据URI
	IsActive      int8      `gorm:"column:is_active;default:1" json:"is_active"`                                    // 是否激活(0:否, 1:是)
	CreateTime    time.Time `gorm:"column:create_time;not null" json:"create_time"`                                 // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`                  // 更新时间
}

// TableName 设置表名
func (u *UserItems) TableName() string {
	return "user_items"
}

// NewUserItems 创建新的用户道具实例
func NewUserItems(userAddress string, itemTokenID int64, itemType int8, itemName string) *UserItems {
	now := time.Now()
	return &UserItems{
		UserAddress: userAddress,
		ItemTokenID: itemTokenID,
		ItemType:    itemType,
		ItemName:    itemName,
		Rarity:      1, // 默认稀有度为普通
		IsActive:    1, // 默认激活
		CreateTime:  now,
		UpdateTime:  now,
	}
}

// GetUserItemByUserAndToken 根据用户地址和道具TokenID获取道具信息
func (dao *Dao) GetUserItemByUserAndToken(ctx context.Context, userAddress string, itemTokenID int64) (*UserItems, error) {
	var userItem UserItems
	err := dao.DB.WithContext(ctx).Where("user_address = ? AND item_token_id = ?", userAddress, itemTokenID).First(&userItem).Error
	if err != nil {
		return nil, err
	}
	return &userItem, nil
}

// GetUserItemsByType 根据道具类型获取用户道具列表
func (dao *Dao) GetUserItemsByType(ctx context.Context, userAddress string, itemType int8) ([]*UserItems, error) {
	var userItems []*UserItems
	err := dao.DB.WithContext(ctx).Where("user_address = ? AND item_type = ?", userAddress, itemType).Find(&userItems).Error
	if err != nil {
		return nil, err
	}
	return userItems, nil
}

// CreateUserItems 创建用户道具记录
func (dao *Dao) CreateUserItems(ctx context.Context, userItems *UserItems) error {
	return dao.DB.WithContext(ctx).Create(userItems).Error
}

// UpdateUserItems 更新用户道具记录
func (dao *Dao) UpdateUserItems(ctx context.Context, userItems *UserItems) error {
	return dao.DB.WithContext(ctx).Save(userItems).Error
}

// DecreaseRemainingUses 减少道具剩余使用次数
func (dao *Dao) DecreaseRemainingUses(ctx context.Context, userItems *UserItems, amount int) error {
	if userItems.RemainingUses < amount {
		return errors.New("insufficient remaining uses")
	}
	userItems.RemainingUses -= amount
	return dao.UpdateUserItems(ctx, userItems)
}
