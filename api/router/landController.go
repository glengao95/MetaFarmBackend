package router

import (
	"MetaFarmBackend/api/middleware"
	"MetaFarmBackend/api/request"
	"MetaFarmBackend/api/response"
	"MetaFarmBackend/component/logger"
	"MetaFarmBackend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// landController 土地相关控制器
type LandController struct {
	landService service.LandService
}

// 构造函数
func NewLandController(landService service.LandService) *LandController {
	return &LandController{
		landService: landService,
	}
}

// RegisterRoutes 注册土地相关路由
func (c *LandController) RegisterRoutes(router *gin.RouterGroup) {
	landRouter := router.Group("/land")
	{
		// 需要身份验证的路由
		landRouter.GET("/list", c.ListUserLands)
		landRouter.GET("/:tokenID/detail", c.GetLandDetail)
		landRouter.POST("/upgrade", c.UpgradeLand)
		landRouter.POST("/rent/list", c.ListRentLands)
		landRouter.POST("/rent/create", c.CreateRent)
		landRouter.POST("/rent/cancel", c.CancelRent)
		landRouter.POST("/market/buy", c.BuyLand)
		landRouter.POST("/layout/update", c.UpdateLayout)
		landRouter.POST("/activity/plant", c.PlantCrop)
		landRouter.POST("/activity/harvest", c.HarvestCrop)
	}
}

// ListUserLands 获取用户土地列表
// @Summary 获取用户所有土地
// @Description 获取当前登录用户的所有土地信息
// @Tags land
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]land.LandInfo}
// @Failure 400 {object} Response{error=string}
// @Router /land/list [get]
func (a *LandController) ListUserLands(ctx *gin.Context) {
	userAddr := ctx.GetHeader("user_address")
	lands, err := a.landService.GetUserLands(ctx, userAddr)
	if err != nil {
		logger.Error("获取土地列表失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, middleware.Response{Data: lands})
}

// GetLandDetail 获取土地详细信息
// @Summary 获取土地详细信息
// @Description 根据tokenID获取土地详细信息
// @Tags land
// @Accept json
// @Produce json
// @Param tokenID path string true "土地NFT TokenID"
// @Success 200 {object} Response{data=dao.LandInfo}
// @Failure 400 {object} Response{error=string}
// @Router /land/{tokenID}/detail [get]
func (a *LandController) GetLandDetail(ctx *gin.Context) {
	tokenID := ctx.Param("tokenID")
	landDetail, err := a.landService.GetLandDetail(ctx, tokenID)
	if err != nil {
		logger.Error("获取土地详情失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, middleware.Response{Data: landDetail})
}

// UpgradeLand 升级土地
// @Summary 升级土地
// @Description 升级指定土地的等级
// @Tags land
// @Accept json
// @Produce json
// @Param body body UpgradeLandRequest true "升级土地请求"
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} Response{error=string}
// @Router /land/upgrade [post]
func (a *LandController) UpgradeLand(ctx *gin.Context) {
	var req request.UpgradeLandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 从请求头获取用户地址
	userAddr := ctx.GetHeader("user_address")
	if userAddr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	req.UserAddress = userAddr

	// 调用服务层升级土地
	err := a.landService.UpgradeLand(ctx, req)
	if err != nil {
		logger.Error("升级土地失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, middleware.Response{Data: "土地升级成功"})
}

// ListRentLands 获取租赁订单列表
// @Summary 获取租赁订单列表
// @Description 获取当前用户的所有租赁订单
// @Tags land
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]dao.LandRental}
// @Failure 400 {object} Response{error=string}
// @Router /land/rent/list [get]
func (a *LandController) ListRentLands(ctx *gin.Context) {
	userAddr := ctx.GetHeader("user_address")
	lands, err := a.landService.GetActiveRentals(ctx, userAddr)
	if err != nil {
		logger.Error("获取租赁列表失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, middleware.Response{Data: lands})
}

// CreateRent 创建土地租赁订单
// @Summary 创建土地租赁订单
// @Description 创建土地租赁订单
// @Tags land
// @Accept json
// @Produce json
// @Param body body CreateRentRequest true "创建租赁请求"
// @Success 200 {object} Response{data=dao.LandRental}
// @Failure 400 {object} Response{error=string}
// @Router /land/rent/create [post]
func (c *LandController) CreateRent(ctx *gin.Context) {
	// 1. 绑定请求参数
	var req request.CreateRentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 调用服务层创建租赁订单
	rental, err := c.landService.CreateRental(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 返回标准化响应
	ctx.JSON(http.StatusOK, response.RentLandResponse{
		RentalID:        rental.ID,
		LandTokenID:     rental.LandTokenID,
		TotalRent:       rental.TotalRent,
		RentalStartTime: rental.RentalStartTime,
		RentalEndTime:   rental.RentalEndTime,
	})
}

// CancelRent 取消土地租赁
// @Summary 取消土地租赁
// @Description 取消土地租赁订单
// @Tags land
// @Accept json
// @Produce json
// @Param body body CancelRentRequest true "取消租赁请求"
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} Response{error=string}
// @Router /land/rent/cancel [post]
func (a *LandController) CancelRent(ctx *gin.Context) {
	var req request.CancelRentalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 从请求头获取用户地址
	userAddr := ctx.GetHeader("user_address")
	if userAddr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	req.UserAddress = userAddr

	// 调用服务层取消租赁
	err := a.landService.CancelRental(ctx, req)
	if err != nil {
		logger.Error("取消租赁失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, middleware.Response{Data: "租赁取消成功"})
}

// BuyLand 购买土地
// @Summary 购买土地
// @Description 从市场购买土地
// @Tags land
// @Accept json
// @Produce json
// @Param body body BuyLandRequest true "购买土地请求"
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} Response{error=string}
// @Router /land/market/buy [post]
func (a *LandController) BuyLand(ctx *gin.Context) {
	var req request.BuyLandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 从请求头获取用户地址
	buyerAddr := ctx.GetHeader("user_address")
	if buyerAddr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	req.BuyerAddress = buyerAddr

	// 调用服务层购买土地
	err := a.landService.BuyLand(ctx, req)
	if err != nil {
		logger.Error("购买土地失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, middleware.Response{Data: "土地购买成功"})
}

// UpdateLayout 更新土地布局
// @Summary 更新土地布局
// @Description 更新土地分区布局
// @Tags land
// @Accept json
// @Produce json
// @Param body body UpdateLayoutRequest true "更新布局请求"
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} Response{error=string}
// @Router /land/layout/update [post]
func (a *LandController) UpdateLayout(ctx *gin.Context) {
	var req request.UpdateLandLayoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 从请求头获取用户地址
	userAddr := ctx.GetHeader("user_address")
	if userAddr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	req.UserAddress = userAddr

	// 调用服务层更新布局
	err := a.landService.UpdateLandLayout(ctx, req)
	if err != nil {
		logger.Error("更新土地布局失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, middleware.Response{Data: "土地布局更新成功"})
}

// PlantCrop 种植作物
// @Summary 种植作物
// @Description 在土地上种植作物
// @Tags land
// @Accept json
// @Produce json
// @Param body body PlantCropRequest true "种植作物请求"
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} Response{error=string}
// @Router /land/activity/plant [post]
func (a *LandController) PlantCrop(ctx *gin.Context) {
	var req request.PlantCropRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 从请求头获取用户地址
	userAddr := ctx.GetHeader("user_address")
	if userAddr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	req.UserAddress = userAddr

	// 调用服务层种植作物
	err := a.landService.PlantCrop(ctx, req)
	if err != nil {
		logger.Error("种植作物失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, middleware.Response{Data: "作物种植成功"})
}

// HarvestCrop 收获作物
// @Summary 收获作物
// @Description 收获成熟的作物
// @Tags land
// @Accept json
// @Produce json
// @Param body body HarvestCropRequest true "收获作物请求"
// @Success 200 {object} Response{data=string}
// @Failure 400 {object} Response{error=string}
// @Router /land/activity/harvest [post]
func (a *LandController) HarvestCrop(ctx *gin.Context) {
	var req request.HarvestCropRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 从请求头获取用户地址
	userAddr := ctx.GetHeader("user_address")
	if userAddr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	req.UserAddress = userAddr

	// 调用服务层收获作物
	err := a.landService.HarvestCrop(ctx, req)
	if err != nil {
		logger.Error("收获作物失败: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, middleware.Response{Data: "作物收获成功"})
}
