package context

import (
	"MetaFarmBankend/src/component/config"
	"MetaFarmBankend/src/component/db"
	"MetaFarmBankend/src/component/logger"
	"MetaFarmBankend/src/component/redis"

	"gorm.io/gorm"
)

type AppContext struct {
	DB    *gorm.DB
	Redis *redis.Store
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

	return &AppContext{
		DB:    db,
		Redis: redis,
	}, nil
}
