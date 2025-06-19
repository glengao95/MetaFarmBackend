package main

import (
	"MetaFarmBankend/src/api/router"
	"MetaFarmBankend/src/component/config"
	"MetaFarmBankend/src/component/context"
)

func main() {
	//读取配置文件
	config, err := config.LoadConfig("component/config/config.toml")
	if err != nil {
		// 处理错误
		panic(err)
	}

	appContext, err := context.NewAppContext(config)
	if err != nil {
		// 处理错误
		panic(err)
	}
	//初始化路由
	r := router.InitRouter(appContext)
	//启动服务
	r.Run(":" + config.API.Port)
}
