package dao

import (
	"context"
	"time"
)

// PlotPlanting 地块种植表结构体
type PlotPlanting struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                   // 主键ID
	PlotID            int64     `gorm:"column:plot_id;uniqueIndex:idx_plot_id" json:"plot_id"`                          // 地块ID
	LandTokenID       int64     `gorm:"column:land_token_id;index:idx_land_token_id" json:"land_token_id"`              // 土地NFT TokenID
	CropID            int       `gorm:"column:crop_id" json:"crop_id"`                                                  // 作物ID
	CropName          string    `gorm:"column:crop_name;size:50" json:"crop_name"`                                      // 作物名称
	PlantedAt         time.Time `gorm:"column:planted_at" json:"planted_at"`                                            // 种植时间
	WaterLevel        int       `gorm:"column:water_level;default:0" json:"water_level"`                                // 水分等级
	FertilizerLevel   int       `gorm:"column:fertilizer_level;default:0" json:"fertilizer_level"`                      // 肥料等级
	PestLevel         int       `gorm:"column:pest_level;default:0" json:"pest_level"`                                  // 虫害等级
	IsHarvestable     int8      `gorm:"column:is_harvestable;default:0;index:idx_is_harvestable" json:"is_harvestable"` // 是否可收获(0:否, 1:是)
	IsPlanted         int8      `gorm:"column:is_planted;default:0;index:idx_is_planted" json:"is_planted"`             // 是否已种植(0:否, 1:是)
	ExpectedYield     int       `gorm:"column:expected_yield;default:0" json:"expected_yield"`                          // 预计产量
	LastWaterTime     time.Time `gorm:"column:last_water_time" json:"last_water_time"`                                  // 最后浇水时间
	LastFertilizeTime time.Time `gorm:"column:last_fertilize_time" json:"last_fertilize_time"`                          // 最后施肥时间
	LastPesticideTime time.Time `gorm:"column:last_pesticide_time" json:"last_pesticide_time"`                          // 最后除虫时间
	CreateTime        time.Time `gorm:"column:create_time;not null" json:"create_time"`                                 // 创建时间
	UpdateTime        time.Time `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`                  // 更新时间
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
func (dao *Dao) GetPlotPlantingByPlotID(ctx context.Context, plotID int64) (*PlotPlanting, error) {
	var plotPlanting PlotPlanting
	err := dao.DB.WithContext(ctx).Where("plot_id = ?", plotID).First(&plotPlanting).Error
	if err != nil {
		return nil, err
	}
	return &plotPlanting, nil
}

// GetPlotPlantingsByLandTokenID 根据土地NFT ID获取种植信息列表
func (dao *Dao) GetPlotPlantingsByLandTokenID(ctx context.Context, landTokenID int64) ([]*PlotPlanting, error) {
	var plotPlantings []*PlotPlanting
	err := dao.DB.WithContext(ctx).Where("land_token_id = ?", landTokenID).Find(&plotPlantings).Error
	if err != nil {
		return nil, err
	}
	return plotPlantings, nil
}

// CreatePlotPlanting 创建地块种植记录
func (dao *Dao) CreatePlotPlanting(ctx context.Context, p *PlotPlanting) error {
	return dao.DB.WithContext(ctx).Create(p).Error
}

// UpdatePlotPlanting 更新地块种植记录
func (dao *Dao) UpdatePlotPlanting(ctx context.Context, p *PlotPlanting) error {
	return dao.DB.WithContext(ctx).Save(p).Error
}
