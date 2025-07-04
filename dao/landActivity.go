package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// 土地活动类型枚举
const (
	ActivityTypePlanting int8 = iota // 0: 种植
	ActivityTypeBreeding             // 1: 养殖
)

// 土地活动状态枚举
const (
	ActivityStatusGrowing   int8 = iota // 0: 生长中
	ActivityStatusHarvested             // 1: 已收获
	ActivityStatusDead                  // 2: 枯萎
)

// LandActivity 土地活动记录表结构体
type LandActivity struct {
	ID              uint64     `gorm:"primaryKey;column:id"`                        // 主键ID
	LandTokenID     string     `gorm:"column:land_token_id;index"`                  // 土地NFT TokenID
	OwnerAddress    string     `gorm:"column:owner_address;type:varchar(42);index"` // 所有者钱包地址
	ActivityType    int8       `gorm:"column:activity_type;index"`                  // 活动类型(0-种植,1-养殖)
	CropAnimalID    uint64     `gorm:"column:crop_animal_id"`                       // 作物/动物ID
	CropAnimalName  string     `gorm:"column:crop_animal_name;type:varchar(50)"`    // 作物/动物名称
	Area            int        `gorm:"column:area"`                                 // 占用面积(㎡)
	FertilityCost   int        `gorm:"column:fertility_cost"`                       // 消耗肥力值
	Start_time      time.Time  `gorm:"column:start_time"`                           // 开始时间
	ExpectedEndTime time.Time  `gorm:"column:expected_end_time"`                    // 预计结束时间
	Status          int8       `gorm:"column:status;default:0;index"`               // 状态(0-进行中,1-已完成,2-已中断)
	ActualEndTime   *time.Time `gorm:"column:actual_end_time"`                      // 实际结束时间
	CreateTime      time.Time  `gorm:"column:create_time"`                          // 创建时间
	UpdateTime      time.Time  `gorm:"column:update_time"`                          // 更新时间
}

func (LandActivity) TableName() string {
	return "land_activity"
}

func NewLandActivity(landTokenID string, ownerAddress string, activityType int8, cropAnimalID uint64, name string, area, fertilityCost int, durationHours int) *LandActivity {
	now := time.Now()
	endTime := now.Add(time.Duration(durationHours) * time.Hour)
	return &LandActivity{
		LandTokenID:     landTokenID,
		OwnerAddress:    ownerAddress,
		ActivityType:    activityType,
		CropAnimalID:    cropAnimalID,
		CropAnimalName:  name,
		Area:            area,
		FertilityCost:   fertilityCost,
		Start_time:      now,
		ExpectedEndTime: endTime,
		Status:          0,
		CreateTime:      now,
		UpdateTime:      now,
	}
}

func (dao *Dao) GetLandActivityByID(ctx context.Context, activityID uint64) (*LandActivity, error) {
	var activity LandActivity
	err := dao.DB.WithContext(ctx).Where("id = ?", activityID).First(&activity).Error
	return &activity, err
}

func (dao *Dao) GetActiveByTokenID(ctx context.Context, landTokenID string) ([]*LandActivity, error) {
	var activities []*LandActivity
	err := dao.DB.WithContext(ctx).Where("land_token_id = ? AND status = 0", landTokenID).Find(&activities).Error
	return activities, err
}

func (dao *Dao) UpdateLandActivityStatus(ctx context.Context, tx *gorm.DB, activityID uint64, status int8) error {
	if tx == nil {
		tx = dao.DB
	}
	updateData := map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}
	if status == 1 {
		updateData["actual_end_time"] = time.Now()
	}

	return tx.WithContext(ctx).Model(&LandActivity{}).Where("id = ?", activityID).Updates(updateData).Error
}

func (dao *Dao) CreateLandActivity(ctx context.Context, tx *gorm.DB, activity *LandActivity) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Create(activity).Error
}
