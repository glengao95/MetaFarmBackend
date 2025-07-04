package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// LandInfo 土地信息表结构体
type LandInfo struct {
	ID              uint64     `gorm:"primaryKey;column:id"`                        // 主键ID
	LandTokenID     string     `gorm:"column:land_token_id;uniqueIndex"`            // 土地NFT TokenID
	OwnerAddress    string     `gorm:"column:owner_address;type:varchar(42);index"` // 所有者钱包地址
	LandType        int8       `gorm:"column:land_type"`                            // 地形类型(0-平原,1-湿地,2-山地)
	Rarity          int8       `gorm:"column:rarity;default:0;index"`               // 稀有度(0-普通,1-稀有,2-史诗,3-传说)
	Area            int        `gorm:"column:area;default:100"`                     // 土地面积(㎡)
	Level           int8       `gorm:"column:level;default:1;index"`                // 土地等级(1-10级)
	Fertility       int        `gorm:"column:fertility;default:100;index"`          // 土地肥力值(0-100)
	SpecialEffect   string     `gorm:"column:special_effect;type:varchar(100)"`     // 特殊效果(如"湿润土地"、"黄金土地")
	LastHarvestTime *time.Time `gorm:"column:last_harvest_time"`                    // 最后收获时间
	MetadataURI     string     `gorm:"column:metadata_uri;type:varchar(255)"`       // 元数据URI
	CreateTime      time.Time  `gorm:"column:create_time"`                          // 创建时间
	UpdateTime      time.Time  `gorm:"column:update_time"`                          // 更新时间
}

func (LandInfo) TableName() string {
	return "land_info"
}

func NewLandInfo(landTokenID string, ownerAddress string, landType int8, rarity int8) *LandInfo {
	now := time.Now()
	return &LandInfo{
		LandTokenID:   landTokenID,
		OwnerAddress:  ownerAddress,
		LandType:      landType,
		Rarity:        rarity,
		Area:          100,
		Level:         1,
		Fertility:     100,
		SpecialEffect: "",
		CreateTime:    now,
		UpdateTime:    now,
	}
}

func (dao *Dao) GetLandInfoByTokenID(ctx context.Context, tokenID string) (*LandInfo, error) {
	var land LandInfo
	err := dao.DB.WithContext(ctx).Where("land_token_id = ?", tokenID).First(&land).Error
	return &land, err
}

func (dao *Dao) GetLandsByOwner(ctx context.Context, ownerAddress string) ([]*LandInfo, error) {
	var lands []*LandInfo
	err := dao.DB.WithContext(ctx).Where("owner_address = ?", ownerAddress).Find(&lands).Error
	return lands, err
}

func (dao *Dao) CreateLandInfo(ctx context.Context, tx *gorm.DB, land *LandInfo) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Create(land).Error
}

func (dao *Dao) UpdateLandInfo(ctx context.Context, tx *gorm.DB, land *LandInfo) error {
	if tx == nil {
		tx = dao.DB
	}
	land.UpdateTime = time.Now()
	return tx.WithContext(ctx).Save(land).Error
}

func (dao *Dao) UpdateLevel(ctx context.Context, tx *gorm.DB, tokenID string, level int8) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Model(&LandInfo{}).Where("land_token_id = ?", tokenID).Update("level", level).Error
}

func (dao *Dao) UpdateFertility(ctx context.Context, tx *gorm.DB, tokenID string, fertility int) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Model(&LandInfo{}).Where("land_token_id = ?", tokenID).Update("fertility", fertility).Error
}

// UpdateLandOwner 更新土地所有者
func (dao *Dao) UpdateLandOwner(ctx context.Context, tx *gorm.DB, tokenID string, newOwner string) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Model(&LandInfo{}).Where("land_token_id = ?", tokenID).UpdateColumns(map[string]interface{}{
		"owner_address": newOwner,
		"update_time":   time.Now(),
	}).Error
}
