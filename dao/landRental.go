package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// 土地租赁状态枚举
const (
	RentalStatusPending   int8 = iota // 0: 待确认
	RentalStatusActive                // 1: 租赁中
	RentalStatusEnded                 // 2: 已结束
	RentalStatusCancelled             // 3: 已取消
)

// LandRental 土地租赁表结构体
type LandRental struct {
	ID               uint64    `gorm:"primaryKey;column:id"`                           // 主键ID
	LandTokenID      string    `gorm:"column:land_token_id;index"`                     // 土地NFT TokenID
	OwnerAddress     string    `gorm:"column:owner_address;type:varchar(42);index"`    // 所有者钱包地址
	RenterAddress    string    `gorm:"column:renter_address;type:varchar(42);index"`   // 租客钱包地址
	RentalDuration   int       `gorm:"column:rental_duration"`                         // 租期(天，7/14/30)
	RentPerSqmPerDay float64   `gorm:"column:rent_per_sqm_per_day;type:decimal(10,6)"` // 每平方米日租金
	TotalRent        float64   `gorm:"column:total_rent;type:decimal(18,6)"`           // 总租金
	SystemFee        float64   `gorm:"column:system_fee;type:decimal(18,6)"`           // 系统手续费(5%)
	Status           int8      `gorm:"column:status;default:0;index"`                  // 状态(0-待确认,1-租赁中,2-已结束,3-已取消)
	RentalStartTime  time.Time `gorm:"column:rental_start_time"`                       // 租赁开始时间
	RentalEndTime    time.Time `gorm:"column:rental_end_time"`                         // 租赁结束时间
	CreateTime       time.Time `gorm:"column:create_time"`                             // 创建时间
	UpdateTime       time.Time `gorm:"column:update_time"`                             // 更新时间
}

func (LandRental) TableName() string {
	return "land_rental"
}

func NewLandRental(landTokenID string, ownerAddress, renterAddress string, duration int, rentPerSqm float64, area int) *LandRental {
	now := time.Now()
	totalRent := float64(area) * rentPerSqm * float64(duration)
	systemFee := totalRent * 0.05
	endTime := now.AddDate(0, 0, duration)
	return &LandRental{
		LandTokenID:      landTokenID,
		OwnerAddress:     ownerAddress,
		RenterAddress:    renterAddress,
		RentalDuration:   duration,
		RentPerSqmPerDay: rentPerSqm,
		TotalRent:        totalRent,
		SystemFee:        systemFee,
		Status:           0,
		RentalStartTime:  now,
		RentalEndTime:    endTime,
		CreateTime:       now,
		UpdateTime:       now,
	}
}

func (dao *Dao) GetLandRentalByID(ctx context.Context, rentalID uint64) (*LandRental, error) {
	var rental LandRental
	err := dao.DB.WithContext(ctx).Where("id = ?", rentalID).First(&rental).Error
	return &rental, err
}

func (dao *Dao) GetLandRentalByTokenID(ctx context.Context, tokenID string) (*LandRental, error) {
	var rental LandRental
	err := dao.DB.WithContext(ctx).Where("land_token_id = ? AND status = 1", tokenID).First(&rental).Error
	return &rental, err
}

func (dao *Dao) CreateLandRental(ctx context.Context, tx *gorm.DB, rental *LandRental) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Create(rental).Error
}

func (dao *Dao) UpdateLandRentalStatus(ctx context.Context, tx *gorm.DB, rentalID uint64, status int8) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Model(&LandRental{}).Where("id = ?", rentalID).UpdateColumns(map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}).Error
}

// 查询用户作为租客的活跃租赁订单
func (dao *Dao) GetLandRentalByRenter(ctx context.Context, renterAddress string) ([]*LandRental, error) {
	var rentals []*LandRental
	err := dao.DB.WithContext(ctx).Where("renter_address = ? AND status = 1", renterAddress).Order("rental_end_time ASC").Find(&rentals).Error
	return rentals, err
}
