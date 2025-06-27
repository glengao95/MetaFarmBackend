package dao

import (
	"time"

	"MetaFarmBackend/component/db"

	"github.com/pkg/errors"
)

// UserItems 用户道具表结构体
type UserItems struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserAddress   string    `gorm:"column:user_address;size:42;not null" json:"user_address"`
	ItemTokenID   int64     `gorm:"column:item_token_id;not null" json:"item_token_id"`
	ItemType      int8      `gorm:"column:item_type;not null;index:idx_item_type" json:"item_type"`
	ItemName      string    `gorm:"column:item_name;size:50;not null" json:"item_name"`
	Rarity        int8      `gorm:"column:rarity;not null;default:1;index:idx_rarity" json:"rarity"`
	Power         int       `gorm:"column:power;default:0" json:"power"`
	MaxUses       int       `gorm:"column:max_uses;default:0" json:"max_uses"`
	RemainingUses int       `gorm:"column:remaining_uses;default:0;index:idx_remaining_uses" json:"remaining_uses"`
	MetadataURI   string    `gorm:"column:metadata_uri;size:255" json:"metadata_uri"`
	IsActive      int8      `gorm:"column:is_active;default:1" json:"is_active"`
	CreateTime    time.Time `gorm:"column:create_time;not null" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`
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
func GetUserItemByUserAndToken(userAddress string, itemTokenID int64) (*UserItems, error) {
	var userItem UserItems
	err := db.GetDB().Where("user_address = ? AND item_token_id = ?", userAddress, itemTokenID).First(&userItem).Error
	if err != nil {
		return nil, err
	}
	return &userItem, nil
}

// GetUserItemsByType 根据道具类型获取用户道具列表
func GetUserItemsByType(userAddress string, itemType int8) ([]*UserItems, error) {
	var userItems []*UserItems
	err := db.GetDB().Where("user_address = ? AND item_type = ?", userAddress, itemType).Find(&userItems).Error
	if err != nil {
		return nil, err
	}
	return userItems, nil
}

// CreateUserItems 创建用户道具记录
func (u *UserItems) CreateUserItems() error {
	return db.GetDB().Create(u).Error
}

// UpdateUserItems 更新用户道具记录
func (u *UserItems) UpdateUserItems() error {
	return db.GetDB().Save(u).Error
}

// DecreaseRemainingUses 减少道具剩余使用次数
func (u *UserItems) DecreaseRemainingUses(amount int) error {
	if u.RemainingUses < amount {
		return errors.New("insufficient remaining uses")
	}
	u.RemainingUses -= amount
	return u.UpdateUserItems()
}
