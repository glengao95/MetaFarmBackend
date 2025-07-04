package router

import (
	"MetaFarmBackend/api/middleware"
	"MetaFarmBackend/component/context"

	"github.com/gin-gonic/gin"
)

func InitRouter(appContext *context.AppContext) *gin.Engine {
	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New() // 新建一个gin引擎实例

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORSMiddleware())

	authController := NewWalletAuthController(appContext.WalletAuthService)
	authController.RegisterRoutes(r)

	apiLand := r.Group("/api/land")
	{
		landController := NewLandController(appContext.LandService)
		landController.RegisterRoutes(apiLand)
	}
	return r
}
