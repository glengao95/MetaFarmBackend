package service

import (
	"MetaFarmBackend/api/request"
	"MetaFarmBackend/component/logger"
	"MetaFarmBackend/dao"
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// LandService 土地系统业务逻辑接口
type LandService interface {
	// 获取用户拥有的土地列表
	GetUserLands(ctx context.Context, userAddress string) ([]*dao.LandInfo, error)
	// 获取土地详细信息
	GetLandDetail(ctx context.Context, tokenID string) (*dao.LandInfo, error)
	// 升级土地
	UpgradeLand(ctx context.Context, req request.UpgradeLandRequest) error
	// 创建土地租赁订单
	CreateRental(ctx context.Context, req request.CreateRentRequest) (*dao.LandRental, error)
	// 获取活跃租赁订单
	GetActiveRentals(ctx context.Context, userAddress string) ([]*dao.LandRental, error)
	// 创建土地挂牌
	CreateMarketListing(ctx context.Context, req request.CreateMarketListingRequest) error
	// 更新土地布局
	UpdateLandLayout(ctx context.Context, req request.UpdateLandLayoutRequest) error
	// 种植作物
	PlantCrop(ctx context.Context, req request.PlantCropRequest) error
	// 收获作物
	HarvestCrop(ctx context.Context, req request.HarvestCropRequest) error
	// 购买土地
	BuyLand(ctx context.Context, req request.BuyLandRequest) error
	// 取消土地租赁
	CancelRental(ctx context.Context, req request.CancelRentalRequest) error
}

type landServiceImpl struct {
	dao *dao.Dao
}

// 构造函数
func NewLandService(dao *dao.Dao) LandService {
	return &landServiceImpl{
		dao: dao,
	}
}

// GetUserLands 获取用户拥有的土地列表
func (s *landServiceImpl) GetUserLands(ctx context.Context, userAddress string) ([]*dao.LandInfo, error) {
	lands, err := s.dao.GetLandsByOwner(ctx, userAddress)
	if err != nil {
		logger.Errorf("获取用户土地列表失败: %v, userAddress: %s", err, userAddress)
		return nil, errors.Wrap(err, "获取土地列表失败")
	}
	return lands, nil
}

// GetLandDetail 获取土地详细信息
func (s *landServiceImpl) GetLandDetail(ctx context.Context, tokenID string) (*dao.LandInfo, error) {
	landInfo, err := s.dao.GetLandInfoByTokenID(ctx, tokenID)
	if err != nil {
		logger.Errorf("获取土地详情失败: %v, tokenID: %s", err, tokenID)
		return nil, errors.Wrap(err, "获取土地详情失败")
	}
	return landInfo, nil
}

// UpgradeLand 升级土地
func (s *landServiceImpl) UpgradeLand(ctx context.Context, req request.UpgradeLandRequest) error {

	// 1. 验证用户权限
	landInfo, err := s.dao.GetLandInfoByTokenID(ctx, req.LandTokenID)
	if err != nil {
		logger.Errorf("获取土地信息失败: %v, tokenID: %s", err, req.LandTokenID)
		return errors.Wrap(err, "获取土地信息失败")
	}
	if landInfo.OwnerAddress != req.UserAddress {
		logger.Errorf("用户无权限升级土地: tokenID=%s, userAddress=%s, ownerAddress=%s", req.LandTokenID, req.UserAddress, landInfo.OwnerAddress)
		return errors.New("无权限升级此土地")
	}

	// 2. 检查土地当前等级和升级条件
	nextLevel := landInfo.Level + 1
	upgradeCost := calculateUpgradeCost(int32(landInfo.Level))
	if upgradeCost == nil {
		return errors.New("已达到最高等级")
	}

	// 3. 扣减升级所需资源 (此处需调用资产服务)
	// TODO: 实现资源扣减逻辑

	// 4. 创建升级记录
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	upgradeRecord := dao.NewLandUpgrade(
		req.LandTokenID,
		req.UserAddress,
		landInfo.Level,
		nextLevel,
		upgradeCost.TokenAmount,
		upgradeCost.ItemIDs,
	)
	if err := s.dao.CreateLandUpgrade(ctx, tx, upgradeRecord); err != nil {
		tx.Rollback()
		logger.Errorf("创建升级记录失败: %v", err)
		return errors.Wrap(err, "创建升级记录失败")
	}

	// 5. 更新土地等级
	if err := s.dao.UpdateLevel(ctx, tx, req.LandTokenID, nextLevel); err != nil {
		tx.Rollback()
		logger.Errorf("更新土地等级失败: %v", err)
		return errors.Wrap(err, "更新土地等级失败")
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交事务失败: %v", err)
		return errors.Wrap(err.Error, "升级土地失败")
	}

	logger.Infof("土地升级成功: tokenID=%s, oldLevel=%d, newLevel=%d", req.LandTokenID, landInfo.Level, nextLevel)
	return nil
}

// calculateUpgradeCost 计算升级成本
func calculateUpgradeCost(currentLevel int32) *UpgradeCost {
	// 简单示例: 每级升级成本递增
	if currentLevel >= 10 {
		return nil // 最高等级
	}
	return &UpgradeCost{
		TokenAmount: uint64(100 * int64(currentLevel+1)),
		ItemIDs:     map[string]int{"1001": 1, "1002": 1}, // 示例道具ID
	}
}

// UpgradeCost 升级成本结构
type UpgradeCost struct {
	TokenAmount uint64         // 代币数量
	ItemIDs     map[string]int // 所需道具ID列表
}

// CreateRental 创建土地租赁订单
func (s *landServiceImpl) CreateRental(ctx context.Context, req request.CreateRentRequest) (*dao.LandRental, error) {
	// 1. 验证土地所有权
	landInfo, err := s.dao.GetLandInfoByTokenID(ctx, req.LandTokenID)
	if err != nil {
		logger.Errorf("获取土地信息失败: %v, tokenID: %s", err, req.LandTokenID)
		return nil, errors.Wrap(err, "获取土地信息失败")
	}
	if landInfo.OwnerAddress != req.UserAddress {
		logger.Errorf("用户非土地所有者: tokenID=%s, ownerAddress=%s, reqOwner=%s", req.LandTokenID, landInfo.OwnerAddress, req.UserAddress)
		return nil, errors.New("无权限出租此土地")
	}

	// 2. 检查土地是否已被租赁
	activeRental, err := s.dao.GetLandRentalByTokenID(ctx, req.LandTokenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("查询土地租赁状态失败: %v, tokenID: %s", err, req.LandTokenID)
		return nil, errors.Wrap(err, "查询租赁状态失败")
	}
	if activeRental != nil {
		return nil, errors.New("土地已处于租赁状态")
	}

	// 3. 验证租赁参数
	if req.RentalDuration <= 0 {
		return nil, errors.New("租赁时长必须大于0")
	}
	if req.RentPerSqm <= 0 {
		return nil, errors.New("租金必须大于0")
	}

	// 4. 检查租客余额 (此处需调用资产服务)
	// TODO: 实现余额检查逻辑

	// 5. 创建租赁订单
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建租赁记录
	landRental := dao.NewLandRental(
		req.LandTokenID,
		landInfo.OwnerAddress,
		req.RenterAddress,
		req.RentalDuration,
		req.RentPerSqm*0.05, // 5%系统手续费
		landInfo.Area,
	)

	if err := s.dao.CreateLandRental(ctx, tx, landRental); err != nil {
		tx.Rollback()
		logger.Errorf("创建租赁记录失败: %v", err)
		return nil, errors.Wrap(err, "创建租赁订单失败")
	}

	// 6. 扣减租客租金并转账给所有者 (此处需调用资产服务)
	// TODO: 实现租金转账逻辑

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交租赁事务失败: %v", err)
		return nil, errors.Wrap(err.Error, "创建租赁订单失败")
	}

	logger.Infof("土地租赁订单创建成功: tokenID=%s, renter=%s, duration=%ds", req.LandTokenID, req.RenterAddress, req.RentalDuration)
	return landRental, nil
}

// GetActiveRentals 获取活跃租赁订单
func (s *landServiceImpl) GetActiveRentals(ctx context.Context, userAddress string) ([]*dao.LandRental, error) {
	// 查询用户作为租客的活跃租赁订单
	rentals, err := s.dao.GetLandRentalByRenter(ctx, userAddress)
	if err != nil {
		logger.Errorf("获取用户活跃租赁订单失败: %v, userAddress: %s", err, userAddress)
		return nil, errors.Wrap(err, "获取租赁订单失败")
	}
	return rentals, nil
}

// CreateMarketListing 创建土地挂牌
func (s *landServiceImpl) CreateMarketListing(ctx context.Context, req request.CreateMarketListingRequest) error {
	// 1. 验证土地所有权
	landInfo, err := s.dao.GetLandInfoByTokenID(ctx, req.TokenID)

	if err != nil {
		logger.Errorf("获取土地信息失败: %v, tokenID: %s", err, req.TokenID)
		return errors.Wrap(err, "获取土地信息失败")
	}
	if landInfo.OwnerAddress != req.SellerAddress {
		logger.Errorf("用户非土地所有者: tokenID=%s, ownerAddress=%s, seller=%s", req.TokenID, landInfo.OwnerAddress, req.SellerAddress)
		return errors.New("无权限挂牌此土地")
	}

	// 2. 检查是否已有活跃挂牌
	activeListing, err := s.dao.GetLandMarketByTokenID(ctx, req.TokenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("查询土地挂牌状态失败: %v, tokenID: %s", err, req.TokenID)
		return errors.Wrap(err, "查询挂牌状态失败")
	}
	if activeListing != nil {
		return errors.New("土地已处于挂牌状态")
	}

	// 3. 验证挂牌价格
	if req.Price <= 0 {
		return errors.New("挂牌价格必须大于0")
	}

	// 4. 创建挂牌记录
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	listing := dao.NewLandMarket(
		req.TokenID,
		req.SellerAddress,
		landInfo.Area,
		req.Price,
	)

	if err := s.dao.CreateLandMarketListing(ctx, tx, listing); err != nil {
		tx.Rollback()
		logger.Errorf("创建土地挂牌失败: %v", err)
		return errors.Wrap(err, "创建土地挂牌失败")
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交挂牌事务失败: %v", err)
		return errors.Wrap(err.Error, "创建土地挂牌失败")
	}

	logger.Infof("土地挂牌成功: tokenID=%s, price=%d", req.TokenID, req.Price)
	return nil
}

// UpdateLandLayout 更新土地布局
func (s *landServiceImpl) UpdateLandLayout(ctx context.Context, req request.UpdateLandLayoutRequest) error {
	// 1. 验证土地所有权
	landInfo, err := s.dao.GetLandInfoByTokenID(ctx, req.TokenID)
	if err != nil {
		logger.Errorf("获取土地信息失败: %v, tokenID: %s", err, req.TokenID)
		return errors.Wrap(err, "获取土地信息失败")
	}
	if landInfo.OwnerAddress != req.UserAddress {
		logger.Errorf("用户非土地所有者: tokenID=%s, ownerAddress=%s, user=%s", req.TokenID, landInfo.OwnerAddress, req.UserAddress)
		return errors.New("无权限更新土地布局")
	}

	// 2. 验证布局数据 (简单检查非空)
	// if req.LayoutData == "" {
	// 	return errors.New("布局数据不能为空")
	// }

	// 3. 查询现有布局记录
	layout, err := s.dao.GetLayoutByTokenID(ctx, req.TokenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("查询土地布局失败: %v, tokenID: %s", err, req.TokenID)
		return errors.Wrap(err, "查询土地布局失败")
	}

	// 4. 更新或创建布局记录
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if layout == nil {
		// 创建新布局
		layout = dao.NewLandLayout(req.TokenID, req.Area, req.ZoneType, req.PosX, req.PosY, req.Width, req.Height)
		if err := s.dao.CreateLandLayout(ctx, tx, layout); err != nil {
			tx.Rollback()
			logger.Errorf("创建土地布局失败: %v", err)
			return errors.Wrap(err, "创建土地布局失败")
		}
	} else {
		// 更新现有布局
		layout.Area = req.Area
		layout.ZoneType = req.ZoneType
		layout.PositionX = req.PosX
		layout.PositionY = req.PosY
		layout.Width = req.Width
		layout.Height = req.Height
		layout.UpdateTime = time.Now()
		if err := s.dao.UpdateLandLayout(ctx, tx, layout); err != nil {
			tx.Rollback()
			logger.Errorf("更新土地布局失败: %v", err)
			return errors.Wrap(err, "更新土地布局失败")
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交布局事务失败: %v", err)
		return errors.Wrap(err.Error, "提交布局事务失败")
	}

	logger.Infof("土地布局更新成功: tokenID=%s", req.TokenID)
	return nil
}

// PlantCrop 种植作物
func (s *landServiceImpl) PlantCrop(ctx context.Context, req request.PlantCropRequest) error {
	// 1. 验证土地所有权
	landInfo, err := s.dao.GetLandInfoByTokenID(ctx, req.LandTokenID)
	if err != nil {
		logger.Errorf("获取土地信息失败: %v, tokenID: %s", err, req.LandTokenID)
		return errors.Wrap(err, "获取土地信息失败")
	}
	if landInfo.OwnerAddress != req.UserAddress {
		logger.Errorf("用户非土地所有者: tokenID=%s, ownerAddress=%s, user=%s", req.LandTokenID, landInfo.OwnerAddress, req.UserAddress)
		return errors.New("无权限种植作物")
	}

	// 2. 检查种植面积
	if req.Area <= 0 || req.Area > landInfo.Area {
		return errors.New("种植面积无效")
	}

	// 3. 检查土地肥力
	requiredFertility := int(req.Area) * 10 // 每单位面积消耗10点肥力
	if landInfo.Fertility < requiredFertility {
		return errors.New("土地肥力不足")
	}

	// 4. 创建种植活动
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 计算预计收获时间 (示例: 2小时)
	startTime := time.Now().Unix()
	endTime := startTime + 7200

	activity := dao.NewLandActivity(
		req.LandTokenID,
		req.UserAddress,
		dao.ActivityTypePlanting,
		req.CropAnimalID,
		"", // 作物名称可从配置表获取
		req.Area,
		int(requiredFertility),
		int(endTime),
	)

	if err := s.dao.CreateLandActivity(ctx, tx, activity); err != nil {
		tx.Rollback()
		logger.Errorf("创建种植活动失败: %v", err)
		return errors.Wrap(err, "种植作物失败")
	}

	// 5. 扣减土地肥力
	if err := s.dao.UpdateFertility(ctx, tx, req.LandTokenID, landInfo.Fertility-requiredFertility); err != nil {
		tx.Rollback()
		logger.Errorf("更新土地肥力失败: %v", err)
		return errors.Wrap(err, "种植作物失败")
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交种植事务失败: %v", err)
		return errors.Wrap(err.Error, "种植作物失败")
	}

	logger.Infof("作物种植成功: tokenID=%s, cropID=%d, area=%.2f", req.LandTokenID, req.CropAnimalID, req.Area)
	return nil
}

// HarvestCrop 收获作物
func (s *landServiceImpl) HarvestCrop(ctx context.Context, req request.HarvestCropRequest) error {
	// 1. 获取活动记录
	activity, err := s.dao.GetLandActivityByID(ctx, req.ActivityID)
	if err != nil {
		logger.Errorf("获取种植活动失败: %v, activityID: %d", err, req.ActivityID)
		return errors.Wrap(err, "获取活动信息失败")
	}

	// 2. 验证权限和状态
	if activity.OwnerAddress != req.UserAddress {
		logger.Errorf("用户非活动所有者: activityID=%d, ownerAddress=%s, user=%s", req.ActivityID, activity.OwnerAddress, req.UserAddress)
		return errors.New("无权限收获此作物")
	}

	if activity.Status != dao.ActivityStatusGrowing {
		return errors.New("作物未处于生长状态")
	}

	// 3. 检查是否成熟
	currentTime := time.Now().Unix()
	if currentTime < activity.ExpectedEndTime.Unix() {
		return errors.New("作物尚未成熟")
	}

	// 4. 更新活动状态和收获
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新活动状态
	if err := s.dao.UpdateLandActivityStatus(ctx, tx, activity.ID, dao.ActivityStatusHarvested); err != nil {
		tx.Rollback()
		logger.Errorf("更新活动状态失败: %v", err)
		return errors.Wrap(err, "收获作物失败")
	}

	// 5. 增加用户资产 (此处需调用资产服务)
	// TODO: 实现收获物发放逻辑

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交收获事务失败: %v", err)
		return errors.Wrap(err.Error, "收获作物失败")
	}

	logger.Infof("作物收获成功: activityID=%d, cropID=%d", req.ActivityID, activity.CropAnimalID)
	return nil
}

// BuyLand 购买土地
func (s *landServiceImpl) BuyLand(ctx context.Context, req request.BuyLandRequest) error {
	// 1. 查询市场挂牌信息
	listing, err := s.dao.GetLandMarketByID(ctx, req.MarketID)
	if err != nil {
		logger.Errorf("查询市场挂牌失败: %v, marketID: %d", err, req.MarketID)
		return errors.Wrap(err, "查询挂牌信息失败")
	}
	if listing == nil {
		return errors.New("土地挂牌不存在")
	}
	if listing.Status != dao.MarketStatusPending {
		return errors.New("土地挂牌已失效")
	}

	// 2. 查询土地信息
	_, err = s.dao.GetLandInfoByTokenID(ctx, listing.LandTokenID)
	if err != nil {
		logger.Errorf("获取土地信息失败: %v, tokenID: %s", err, listing.LandTokenID)
		return errors.Wrap(err, "获取土地信息失败")
	}

	//

	// 3. 验证买家余额 (此处需调用资产服务)
	// TODO: 实现余额检查逻辑

	// 4. 执行交易
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新土地所有者
	if err := s.dao.UpdateLandOwner(ctx, tx, listing.LandTokenID, req.BuyerAddress); err != nil {
		tx.Rollback()
		logger.Errorf("更新土地所有者失败: %v", err)
		return errors.Wrap(err, "购买土地失败")
	}

	// 更新挂牌状态为已售出
	if err := s.dao.UpdateMarketStatusByID(ctx, tx, req.MarketID, dao.MarketStatusSold, req.BuyerAddress); err != nil {
		tx.Rollback()
		logger.Errorf("更新挂牌状态失败: %v", err)
		return errors.Wrap(err, "购买土地失败")
	}

	// 创建交易记录
	// TODO: 实现交易记录创建逻辑

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交购买事务失败: %v", err)
		return errors.Wrap(err.Error, "购买土地失败")
	}

	logger.Infof("土地购买成功: marketID=%d, tokenID=%s, buyer=%s", req.MarketID, listing.LandTokenID, req.BuyerAddress)
	return nil
}

// CancelRental 取消土地租赁
func (s *landServiceImpl) CancelRental(ctx context.Context, req request.CancelRentalRequest) error {
	// 1. 查询租赁订单
	rental, err := s.dao.GetLandRentalByID(ctx, req.RentalID)
	if err != nil {
		logger.Errorf("查询租赁订单失败: %v, rentalID: %d", err, req.RentalID)
		return errors.Wrap(err, "查询租赁信息失败")
	}
	if rental == nil {
		return errors.New("租赁订单不存在")
	}

	// 2. 验证权限
	if rental.RenterAddress != req.UserAddress && rental.OwnerAddress != req.UserAddress {
		logger.Errorf("无权限取消租赁: rentalID=%d, user=%s, renter=%s, owner=%s",
			req.RentalID, req.UserAddress, rental.RenterAddress, rental.OwnerAddress)
		return errors.New("无权限取消此租赁")
	}

	// 3. 检查状态
	if rental.Status != dao.RentalStatusActive {
		return errors.New("租赁订单未处于活跃状态")
	}

	// 4. 取消租赁
	tx := s.dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新租赁状态
	if err := s.dao.UpdateLandRentalStatus(ctx, tx, rental.ID, dao.RentalStatusCancelled); err != nil {
		tx.Rollback()
		logger.Errorf("更新租赁状态失败: %v", err)
		return errors.Wrap(err, "取消租赁失败")
	}

	// TODO: 实现退款逻辑

	if err := tx.Commit(); err != nil {
		logger.Errorf("提交取消租赁事务失败: %v", err)
		return errors.Wrap(err.Error, "取消租赁失败")
	}

	logger.Infof("租赁取消成功: rentalID=%d, tokenID=%s", req.RentalID, rental.LandTokenID)
	return nil
}
