package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// LandLayout 土地分区信息表结构体
type LandLayout struct {
	ID               uint64    `gorm:"primaryKey;column:id"`                    // 主键ID
	LandTokenID      string    `gorm:"column:land_token_id;index"`              // 土地NFT TokenID
	Area             int       `gorm:"column:area"`                             // 分区面积(㎡)
	ZoneType         int8      `gorm:"column:zone_type;index"`                  // 分区类型(0-种植区,1-养殖区,2-装饰区)
	PositionX        int       `gorm:"column:position_x"`                       // 分区左上角X坐标
	PositionY        int       `gorm:"column:position_y"`                       // 分区左上角Y坐标
	Width            int       `gorm:"column:width"`                            // 分区宽度
	Height           int       `gorm:"column:height"`                           // 分区高度
	HasAdjacentBonus bool      `gorm:"column:has_adjacent_bonus;default:false"` // 是否激活相邻加成
	BonusType        int8      `gorm:"column:bonus_type;default:-1"`            // 加成类型(-1-无加成,0-肥力恢复,1-产量提升)
	BonusValue       float64   `gorm:"column:bonus_value;type:decimal(5,2)"`    // 加成值(百分比)
	CreateTime       time.Time `gorm:"column:create_time"`                      // 创建时间
	UpdateTime       time.Time `gorm:"column:update_time"`                      // 更新时间
}

func (LandLayout) TableName() string {
	return "land_layout"
}

func NewLandLayout(landTokenID string, area int, zoneType int8, posX, posY, width, height int) *LandLayout {
	now := time.Now()
	return &LandLayout{
		LandTokenID:      landTokenID,
		Area:             area,
		ZoneType:         zoneType,
		PositionX:        posX,
		PositionY:        posY,
		Width:            width,
		Height:           height,
		HasAdjacentBonus: false,
		BonusType:        -1,
		BonusValue:       0,
		CreateTime:       now,
		UpdateTime:       now,
	}
}

func (dao *Dao) GetLayoutByTokenID(ctx context.Context, tokenID string) (*LandLayout, error) {
	var layout LandLayout
	err := dao.DB.WithContext(ctx).Where("land_token_id = ?", tokenID).First(&layout).Error
	return &layout, err
}

func (dao *Dao) GetLayoutsByTokenID(ctx context.Context, tokenID string) ([]*LandLayout, error) {
	var layouts []*LandLayout
	err := dao.DB.WithContext(ctx).Where("land_token_id = ?", tokenID).Find(&layouts).Error
	return layouts, err
}

func (dao *Dao) CreateLandLayout(ctx context.Context, tx *gorm.DB, layout *LandLayout) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Create(layout).Error
}

func (dao *Dao) UpdateLandLayout(ctx context.Context, tx *gorm.DB, layout *LandLayout) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Save(layout).Error
}

func (dao *Dao) UpdateBonus(ctx context.Context, tx *gorm.DB, tokenID string, zoneID uint64, hasBonus bool, bonusType int8, bonusValue float64) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Model(&LandLayout{}).Where("land_token_id = ? AND id = ?", tokenID, zoneID).UpdateColumns(map[string]interface{}{
		"has_adjacent_bonus": hasBonus,
		"bonus_type":         bonusType,
		"bonus_value":        bonusValue,
		"update_time":        time.Now(),
	}).Error
}

func (dao *Dao) UpdateLandLayoutBonus(ctx context.Context, tx *gorm.DB, tokenID uint64, zoneID uint64, hasBonus bool, bonusType int8, bonusValue float64) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Model(&LandLayout{}).Where("land_token_id = ? AND id = ?", tokenID, zoneID).UpdateColumns(map[string]interface{}{
		"has_adjacent_bonus": hasBonus,
		"bonus_type":         bonusType,
		"bonus_value":        bonusValue,
		"update_time":        time.Now(),
	}).Error
}
