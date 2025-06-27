package dao

import (
	"MetaFarmBackend/component/db"
	"time"
)

type LandInfo struct {
	ID              uint64     `gorm:"primaryKey;column:id"`
	LandTokenID     uint64     `gorm:"column:land_token_id;uniqueIndex"`
	OwnerAddress    string     `gorm:"column:owner_address;type:varchar(42);index"`
	LandType        int8       `gorm:"column:land_type"`
	Area            int        `gorm:"column:area;default:100"`
	Level           int8       `gorm:"column:level;default:1;index"`
	Fertility       int        `gorm:"column:fertility;default:100;index"`
	IsLocked        bool       `gorm:"column:is_locked;default:false"`
	RenterAddress   string     `gorm:"column:renter_address;type:varchar(42)"`
	RentalStartTime *time.Time `gorm:"column:rental_start_time"`
	RentalEndTime   *time.Time `gorm:"column:rental_end_time"`
	LastHarvestTime *time.Time `gorm:"column:last_harvest_time"`
	MetadataURI     string     `gorm:"column:metadata_uri;type:varchar(255)"`
	CreateTime      time.Time  `gorm:"column:create_time"`
	UpdateTime      time.Time  `gorm:"column:update_time"`
}

func (LandInfo) TableName() string {
	return "land_info"
}

func NewLandInfo(landTokenID uint64, ownerAddress string, landType int8) *LandInfo {
	now := time.Now()
	return &LandInfo{
		LandTokenID:  landTokenID,
		OwnerAddress: ownerAddress,
		LandType:     landType,
		Area:         100,
		Level:        1,
		Fertility:    100,
		IsLocked:     false,
		CreateTime:   now,
		UpdateTime:   now,
	}
}

func GetLandInfoByTokenID(tokenID uint64) (*LandInfo, error) {
	var land LandInfo
	err := db.GetDB().Where("land_token_id = ?", tokenID).First(&land).Error
	return &land, err
}

func GetLandInfosByOwner(ownerAddress string) ([]*LandInfo, error) {
	var lands []*LandInfo
	err := db.GetDB().Where("owner_address = ?", ownerAddress).Find(&lands).Error
	return lands, err
}

func CreateLandInfo(land *LandInfo) error {
	return db.GetDB().Create(land).Error
}

func UpdateLandInfo(land *LandInfo) error {
	land.UpdateTime = time.Now()
	return db.GetDB().Save(land).Error
}
