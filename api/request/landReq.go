package request

// GetLandDetailRequest 获取土地详情请求
type GetLandDetailRequest struct {
	LandTokenID string `json:"landTokenId" binding:"required"` // 土地NFT唯一标识
}

// UpgradeLandRequest 升级土地请求
type UpgradeLandRequest struct {
	LandTokenID string `json:"landTokenId" binding:"required"`   // 土地NFT唯一标识
	UserAddress string `json:"userAddress" binding:"required,max=42"` // 用户钱包地址(42位)
	Level       int8   `json:"level" binding:"required,min=1,max=10"` // 升级目标等级(1-10)
}

// CreateRentRequest 创建租赁请求
type CreateRentRequest struct {
	LandTokenID    string  `json:"landTokenId" binding:"required"`      // 土地NFT唯一标识
	RenterAddress  string  `json:"renterAddress" binding:"required,max=42"` // 租户钱包地址
	RentalDuration int     `json:"rentalDuration" binding:"required,oneof=7 14 30"` // 租赁时长(天)
	RentPerSqm     float64 `json:"rentPerSqm" binding:"required,min=0"`   // 每平方米租金
	UserAddress    string  `json:"userAddress" binding:"required,max=42"`   // 操作人钱包地址
}

// BuyLandRequest 购买土地请求
type BuyLandRequest struct {
	MarketID     uint64 `json:"marketId" binding:"required"`      // 市场ID
	ListingID    uint64 `json:"listingId" binding:"required"`     // 挂牌ID
	BuyerAddress string `json:"buyerAddress" binding:"required,max=42"` // 买家钱包地址
}

// LayoutZoneReq 布局区域请求
type LayoutZoneReq struct {
	Area      int  `json:"area" binding:"required,min=1"`          // 区域面积
	ZoneType  int8 `json:"zoneType" binding:"required,oneof=0 1 2"` // 区域类型(0-普通,1-特殊,2-稀有)
	PositionX int  `json:"positionX" binding:"required,min=0"`      // X坐标位置
	PositionY int  `json:"positionY" binding:"required,min=0"`      // Y坐标位置
	Width     int  `json:"width" binding:"required,min=1"`         // 区域宽度
	Height    int  `json:"height" binding:"required,min=1"`        // 区域高度
}

// PlantCropRequest 种植作物请求
type PlantCropRequest struct {
	LandTokenID  string `json:"landTokenId" binding:"required"`    // 土地NFT唯一标识
	ZoneID       uint64 `json:"zoneId" binding:"required"`         // 区域ID
	CropAnimalID uint64 `json:"cropAnimalId" binding:"required"`   // 作物/动物ID
	UserAddress  string `json:"userAddress" binding:"required,max=42"` // 用户钱包地址
	Area         int    `json:"area" binding:"required,min=1"`     // 种植面积
}

// HarvestCropRequest 收获作物请求
type HarvestCropRequest struct {
	ActivityID  uint64 `json:"activityId" binding:"required"`      // 活动ID
	UserAddress string `json:"userAddress" binding:"required,max=42"` // 用户钱包地址
}

// ListRentLandsRequest 获取租赁土地列表请求
type ListRentLandsRequest struct {
	RenterAddress string `json:"renterAddress" binding:"required,max=42"` // 租户钱包地址
}

// CreateMarketListingRequest 创建土地挂牌请求
type CreateMarketListingRequest struct {
	TokenID       string  `json:"tokenId" binding:"required"`        // 土地NFT ID
	SellerAddress string  `json:"sellerAddress" binding:"required"`  // 卖家地址
	Price         float64 `json:"price" binding:"required,min=0"`    // 售价(非负)
}

// UpdateLandLayoutRequest 更新土地布局请求
type UpdateLandLayoutRequest struct {
	TokenID     string `json:"tokenId" binding:"required"`       // 土地NFT ID
	UserAddress string `json:"userAddress" binding:"required"`   // 用户地址
	Area        int    `json:"area" binding:"required,min=1"`    // 种植面积
	ZoneType    int8   `json:"zoneType" binding:"required"`      // 区域类型
	PosX        int    `json:"posX" binding:"required,min=0"`    // 布局X坐标
	PosY        int    `json:"posY" binding:"required,min=0"`    // 布局Y坐标
	Width       int    `json:"width" binding:"required,min=1"`   // 布局宽度
	Height      int    `json:"height" binding:"required,min=1"`  // 布局高度
}

// CancelRentalRequest 取消租赁请求
type CancelRentalRequest struct {
	RentalID    uint64 `json:"rentalId" binding:"required"`      // 租赁订单ID
	UserAddress string `json:"userAddress" binding:"required,max=42"` // 用户地址
}
