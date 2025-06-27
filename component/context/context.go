package context

import (
	"context"
	"time"

	"MetaFarmBackend/component/config"
	"MetaFarmBackend/component/db"
	"MetaFarmBackend/component/logger"
	"MetaFarmBackend/component/redis"
	"MetaFarmBackend/dao"
	"MetaFarmBackend/service"

	"gorm.io/gorm"
)

type AppContext struct {
	DB                *gorm.DB
	Redis             *redis.Store
	Dao               *dao.Dao
	WalletAuthService service.WalletAuthService
}

func NewAppContext(config *config.Config) (*AppContext, error) {

	//初始化日志
	if err := logger.InitLogger(config); err != nil {
		// 处理错误
		panic(err)
	}

	//初始化gorm
	db, err := db.InitDB(config)
	if err != nil {
		// 处理错误
		panic(err)
	}

	//初始化redis
	redis, err := redis.InitRedis(config)
	if err != nil {
		// 处理错误
		panic(err)
	}
	dao := dao.NewDao(context.Background(), db, redis)
	walletAuthService := service.NewWalletAuthService(dao, time.Duration(config.API.SessionTTL)*time.Second)
	return &AppContext{
		DB:                db,
		Redis:             redis,
		Dao:               dao,
		WalletAuthService: walletAuthService,
	}, nil
}
