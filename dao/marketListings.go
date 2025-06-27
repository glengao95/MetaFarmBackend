package dao

import (
	"database/sql"
	"time"

	"MetaFarmBackend/component/db"
)

// MarketListings 市场挂单表结构体
type MarketListings struct {
	ID                 int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ListingID          int64           `gorm:"column:listing_id;not null;uniqueIndex:idx_listing_id" json:"listing_id"`
	SellerAddress      string          `gorm:"column:seller_address;size:42;not null;index:idx_seller_address" json:"seller_address"`
	NFTContractAddress string          `gorm:"column:nft_contract_address;size:42;not null;index:idx_nft_contract" json:"nft_contract_address"`
	TokenID            int64           `gorm:"column:token_id;not null" json:"token_id"`
	NFTType            int8            `gorm:"column:nft_type;not null" json:"nft_type"`
	SaleType           int8            `gorm:"column:sale_type;not null" json:"sale_type"`
	Price              sql.NullFloat64 `gorm:"column:price;type:decimal(36,18);not null" json:"price"`
	StartTime          time.Time       `gorm:"column:start_time;not null" json:"start_time"`
	EndTime            sql.NullTime    `gorm:"column:end_time" json:"end_time"`
	MinBid             sql.NullFloat64 `gorm:"column:min_bid;type:decimal(36,18)" json:"min_bid"`
	HighestBidder      string          `gorm:"column:highest_bidder;size:42" json:"highest_bidder"`
	HighestBid         sql.NullFloat64 `gorm:"column:highest_bid;type:decimal(36,18)" json:"highest_bid"`
	Status             int8            `gorm:"column:status;not null;default:1;index:idx_status" json:"status"`
	CreateTime         time.Time       `gorm:"column:create_time;not null;index:idx_create_time" json:"create_time"`
	UpdateTime         time.Time       `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`
}

// TableName 设置表名
func (m *MarketListings) TableName() string {
	return "market_listings"
}

// NewMarketListing 创建新的市场挂单实例
func NewMarketListing(listingID int64, sellerAddress string, nftContractAddress string, tokenID int64, nftType int8, saleType int8, price float64, startTime time.Time) *MarketListings {
	now := time.Now()
	return &MarketListings{
		ListingID:          listingID,
		SellerAddress:      sellerAddress,
		NFTContractAddress: nftContractAddress,
		TokenID:            tokenID,
		NFTType:            nftType,
		SaleType:           saleType,
		Price: sql.NullFloat64{
			Float64: price,
			Valid:   true,
		},
		StartTime:  startTime,
		Status:     1, // 默认状态为活跃
		CreateTime: now,
		UpdateTime: now,
	}
}

// GetMarketListingByID 根据挂单ID获取挂单信息
func GetMarketListingByID(listingID int64) (*MarketListings, error) {
	var listing MarketListings
	err := db.GetDB().Where("listing_id = ?", listingID).First(&listing).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

// GetActiveListings 获取活跃挂单列表
func GetActiveListings(page, pageSize int) ([]*MarketListings, int64, error) {
	var listings []*MarketListings
	var total int64

	// 获取总数
	err := db.GetDB().Model(&MarketListings{}).Where("status = 1 AND start_time <= NOW() AND (end_time IS NULL OR end_time >= NOW())").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = db.GetDB().Where("status = 1 AND start_time <= NOW() AND (end_time IS NULL OR end_time >= NOW())").Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&listings).Error
	if err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}

// CreateMarketListing 创建挂单记录
func (m *MarketListings) CreateMarketListing() error {
	return db.GetDB().Create(m).Error
}

// UpdateMarketListingStatus 更新挂单状态
func (m *MarketListings) UpdateMarketListingStatus(status int8) error {
	return db.GetDB().Model(m).Update("status", status).Error
}

// UpdateBidInfo 更新竞价信息
func (m *MarketListings) UpdateBidInfo(bidder string, bidAmount float64) error {
	return db.GetDB().Model(m).Updates(map[string]interface{}{
		"highest_bidder": bidder,
		"highest_bid":    bidAmount,
	}).Error
}
