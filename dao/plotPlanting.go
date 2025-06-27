package dao

import (
	"time"

	"MetaFarmBackend/component/db"
)

// PlotPlanting 地块种植表结构体
type PlotPlanting struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PlotID            int64     `gorm:"column:plot_id;uniqueIndex:idx_plot_id" json:"plot_id"`
	LandTokenID       int64     `gorm:"column:land_token_id;index:idx_land_token_id" json:"land_token_id"`
	CropID            int       `gorm:"column:crop_id" json:"crop_id"`
	CropName          string    `gorm:"column:crop_name;size:50" json:"crop_name"`
	PlantedAt         time.Time `gorm:"column:planted_at" json:"planted_at"`
	WaterLevel        int       `gorm:"column:water_level;default:0" json:"water_level"`
	FertilizerLevel   int       `gorm:"column:fertilizer_level;default:0" json:"fertilizer_level"`
	PestLevel         int       `gorm:"column:pest_level;default:0" json:"pest_level"`
	IsHarvestable     int8      `gorm:"column:is_harvestable;default:0;index:idx_is_harvestable" json:"is_harvestable"`
	IsPlanted         int8      `gorm:"column:is_planted;default:0;index:idx_is_planted" json:"is_planted"`
	ExpectedYield     int       `gorm:"column:expected_yield;default:0" json:"expected_yield"`
	LastWaterTime     time.Time `gorm:"column:last_water_time" json:"last_water_time"`
	LastFertilizeTime time.Time `gorm:"column:last_fertilize_time" json:"last_fertilize_time"`
	LastPesticideTime time.Time `gorm:"column:last_pesticide_time" json:"last_pesticide_time"`
	CreateTime        time.Time `gorm:"column:create_time;not null" json:"create_time"`
	UpdateTime        time.Time `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`
}

// TableName 设置表名
func (p *PlotPlanting) TableName() string {
	return "plot_planting"
}

// NewPlotPlanting 创建新的地块种植实例
func NewPlotPlanting(plotID, landTokenID int64) *PlotPlanting {
	now := time.Now()
	return &PlotPlanting{
		PlotID:      plotID,
		LandTokenID: landTokenID,
		CreateTime:  now,
		UpdateTime:  now,
	}
}

// GetPlotPlantingByPlotID 根据地块ID获取种植信息
func GetPlotPlantingByPlotID(plotID int64) (*PlotPlanting, error) {
	var plotPlanting PlotPlanting
	err := db.GetDB().Where("plot_id = ?", plotID).First(&plotPlanting).Error
	if err != nil {
		return nil, err
	}
	return &plotPlanting, nil
}

// GetPlotPlantingsByLandTokenID 根据土地NFT ID获取种植信息列表
func GetPlotPlantingsByLandTokenID(landTokenID int64) ([]*PlotPlanting, error) {
	var plotPlantings []*PlotPlanting
	err := db.GetDB().Where("land_token_id = ?", landTokenID).Find(&plotPlantings).Error
	if err != nil {
		return nil, err
	}
	return plotPlantings, nil
}

// CreatePlotPlanting 创建地块种植记录
func (p *PlotPlanting) CreatePlotPlanting() error {
	return db.GetDB().Create(p).Error
}

// UpdatePlotPlanting 更新地块种植记录
func (p *PlotPlanting) UpdatePlotPlanting() error {
	return db.GetDB().Save(p).Error
}
