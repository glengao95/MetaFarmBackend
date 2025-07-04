package response

import (
	"MetaFarmBackend/dao"
	"time"
)

// LandDetailResponse 土地详情响应
type LandDetailResponse struct {
	LandTokenID     string     `json:"landTokenId"`     // 土地NFT唯一标识
	OwnerAddress    string     `json:"ownerAddress"`    // 所有者钱包地址
	LandType        int8       `json:"landType"`        // 土地类型(0-普通,1-特殊,2-稀有)
	Rarity          int8       `json:"rarity"`          // 稀有度(1-5星)
	Area            int        `json:"area"`            // 土地面积(平方米)
	Level           int8       `json:"level"`           // 土地等级(1-10)
	Fertility       int        `json:"fertility"`       // 土地肥力(0-100)
	SpecialEffect   string     `json:"specialEffect"`   // 特殊效果描述
	LastHarvestTime *time.Time `json:"lastHarvestTime"` // 最后收获时间
	MetadataURI     string     `json:"metadataUri"`     // 元数据URI
}

// UpgradeLandResponse 土地升级响应
type UpgradeLandResponse struct {
	Success  bool    `json:"success"`  // 升级是否成功
	NewLevel int8    `json:"newLevel"` // 升级后的等级
	Cost     float64 `json:"cost"`     // 升级消耗资源数量
	Message  string  `json:"message"`  // 操作结果描述
}

// RentLandResponse 租赁土地响应
type RentLandResponse struct {
	RentalID        uint64    `json:"rentalId"`        // 租赁订单ID
	LandTokenID     string    `json:"landTokenId"`     // 土地NFT唯一标识
	TotalRent       float64   `json:"totalRent"`       // 总租金
	RentalStartTime time.Time `json:"rentalStartTime"` // 租赁开始时间
	RentalEndTime   time.Time `json:"rentalEndTime"`   // 租赁结束时间
}

// BuyLandResponse 购买土地响应
type BuyLandResponse struct {
	Success         bool   `json:"success"`         // 购买是否成功
	LandTokenID     string `json:"landTokenId"`     // 土地NFT唯一标识
	TransactionHash string `json:"transactionHash,omitempty"` // 交易哈希(可选)
}

// LandLayoutResponse 土地布局响应
type LandLayoutResponse struct {
	LandTokenID string          `json:"landTokenId"` // 土地NFT唯一标识
	Layouts     []LayoutZoneRes `json:"layouts"`     // 布局区域列表
}

// LayoutZoneRes 布局区域详情
type LayoutZoneRes struct {
	ID               uint64  `json:"id"`                // 区域ID
	Area             int     `json:"area"`              // 区域面积
	ZoneType         int8    `json:"zoneType"`          // 区域类型(0-普通,1-特殊,2-稀有)
	PositionX        int     `json:"positionX"`         // X坐标位置
	PositionY        int     `json:"positionY"`         // Y坐标位置
	Width            int     `json:"width"`             // 区域宽度
	Height           int     `json:"height"`            // 区域高度
	HasAdjacentBonus bool    `json:"hasAdjacentBonus"`  // 是否有相邻奖励
	BonusType        int8    `json:"bonusType"`         // 奖励类型
	BonusValue       float64 `json:"bonusValue"`        // 奖励值
}

// PlantCropResponse 种植作物响应
type PlantCropResponse struct {
	ActivityID      uint64    `json:"activityId"`      // 活动ID
	ExpectedEndTime time.Time `json:"expectedEndTime"` // 预计完成时间
}

// HarvestCropResponse 收获作物响应
type HarvestCropResponse struct {
	Success    bool   `json:"success"`    // 收获是否成功
	Yield      int    `json:"yield"`      // 产量数量
	CropName   string `json:"cropName"`   // 作物名称
	Experience int    `json:"experience,omitempty"` // 获得经验值(可选)
}

// RentLandsListResponse 租赁土地列表响应
type RentLandsListResponse struct {
	Total   int                `json:"total"`   // 总记录数
	Rentals []RentLandResponse `json:"rentals"` // 租赁列表
}

// ToLandDetailResponse 将dao.LandInfo转换为API响应结构体
func ToLandDetailResponse(land *dao.LandInfo) *LandDetailResponse {
	return &LandDetailResponse{
		LandTokenID:     land.LandTokenID,
		OwnerAddress:    land.OwnerAddress,
		LandType:        land.LandType,
		Rarity:          land.Rarity,
		Area:            land.Area,
		Level:           land.Level,
		Fertility:       land.Fertility,
		SpecialEffect:   land.SpecialEffect,
		LastHarvestTime: land.LastHarvestTime,
		MetadataURI:     land.MetadataURI,
	}
}
