package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// 土地市场状态枚举
const (
	MarketStatusPending   int8 = iota // 0: 待出售
	MarketStatusSold                  // 1: 已售出
	MarketStatusCancelled             // 2: 已取消
)

// LandMarket 土地交易市场表结构体
type LandMarket struct {
	ID              uint64     `gorm:"primaryKey;column:id"`                         // 主键ID
	LandTokenID     string     `gorm:"column:land_token_id;uniqueIndex"`             // 土地NFT TokenID
	SellerAddress   string     `gorm:"column:seller_address;type:varchar(42);index"` // 卖家钱包地址
	BuyerAddress    string     `gorm:"column:buyer_address;type:varchar(42);index"`  // 买家钱包地址
	Area            int        `gorm:"column:area"`                                  // 土地面积(㎡)
	Price           float64    `gorm:"column:price;type:decimal(18,6)"`              // 售价
	Status          int8       `gorm:"column:status;default:0;index"`                // 状态(0-待出售,1-已售出,2-已取消)
	ListingTime     time.Time  `gorm:"column:listing_time"`                          // 挂牌时间
	TransactionTime *time.Time `gorm:"column:transaction_time"`                      // 交易完成时间
	CreateTime      time.Time  `gorm:"column:create_time"`                           // 创建时间
	UpdateTime      time.Time  `gorm:"column:update_time"`                           // 更新时间
}

func (LandMarket) TableName() string {
	return "land_market"
}

func NewLandMarket(landTokenID string, sellerAddress string, area int, price float64) *LandMarket {
	now := time.Now()
	return &LandMarket{
		LandTokenID:   landTokenID,
		SellerAddress: sellerAddress,
		Area:          area,
		Price:         price,
		Status:        0,
		ListingTime:   now,
		CreateTime:    now,
		UpdateTime:    now,
	}
}

func (dao *Dao) GetActiveLandListings(ctx context.Context) ([]*LandMarket, error) {
	var listings []*LandMarket
	err := dao.DB.WithContext(ctx).Where("status = 0").Order("price ASC").Find(&listings).Error
	return listings, err
}

func (dao *Dao) UpdateMarketStatusByID(ctx context.Context, tx *gorm.DB, id uint64, status int8, buyerAddress string) error {
	updateData := map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}
	if status == MarketStatusSold {
		updateData["buyer_address"] = buyerAddress
		updateData["transaction_time"] = time.Now()
	}

	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Model(&LandMarket{}).Where("id = ?", id).UpdateColumns(updateData).Error
}

func (dao *Dao) CreateLandMarketListing(ctx context.Context, tx *gorm.DB, listing *LandMarket) error {
	if tx == nil {
		tx = dao.DB
	}
	return tx.WithContext(ctx).Create(listing).Error
}

// 检查是否已有活跃挂牌
func (dao *Dao) GetLandMarketByID(ctx context.Context, id uint64) (*LandMarket, error) {
	var listing LandMarket
	err := dao.DB.WithContext(ctx).Where("id = ?", id).First(&listing).Error
	return &listing, err
}

func (dao *Dao) GetLandMarketByTokenID(ctx context.Context, tokenID string) (*LandMarket, error) {
	var listing LandMarket
	err := dao.DB.WithContext(ctx).Where("land_token_id = ? AND status = 0", tokenID).First(&listing).Error
	return &listing, err
}
