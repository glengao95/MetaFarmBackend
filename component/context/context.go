package context

import (
	"context"
	"time"

	"MetaFarmBackend/component/blockchain"
	"MetaFarmBackend/component/cache"
	"MetaFarmBackend/component/config"
	"MetaFarmBackend/component/db"
	"MetaFarmBackend/component/logger"
	"MetaFarmBackend/component/redis"
	"MetaFarmBackend/dao"
	"MetaFarmBackend/service"
)

type AppContext struct {
	Cache             *cache.CacheService
	Dao               *dao.Dao
	WalletAuthService service.WalletAuthService
	LandService       service.LandService
	EthClient         *blockchain.EthClient
	ZkSyncClient      *blockchain.ZkSync2Client
	ZkBridge          *blockchain.ZkSyncBridge
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
	//初始化缓存
	cache := cache.NewCacheService(redis)

	d := dao.NewDao(context.Background(), db, redis)
	//初始化表
	dao.InitTable()

	//初始化服务
	walletAuthService := service.NewWalletAuthService(d, time.Duration(config.API.SessionTTL)*time.Second)
	landService := service.NewLandService(d)

	// 初始化以太坊客户端
	ethClient, err := blockchain.NewEthClient(config.Ethereum.RPCURL, config.Ethereum.PrivateKey)
	if err != nil {
		panic(err)
	}

	// 初始化zkSync客户端
	zkSyncClient, err := blockchain.NewZkSync2Client(config.ZkSync.RPCURL, config.ZkSync.PrivateKey, ethClient.GetClient())
	if err != nil {
		panic(err)
	}

	// 初始化zkSync桥接
	zkBridge, err := blockchain.NewZkSyncBridge(ethClient, zkSyncClient, config.ZkSync.BridgeAddress)
	if err != nil {
		panic(err)
	}

	return &AppContext{
		Cache:             cache,
		Dao:               d,
		WalletAuthService: walletAuthService,
		LandService:       landService,
		EthClient:         ethClient,
		ZkSyncClient:      zkSyncClient,
		ZkBridge:          zkBridge,
	}, nil
}
