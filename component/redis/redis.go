package redis

import (
	"MetaFarmBackend/component/config"
	"MetaFarmBackend/component/logger"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// getAndDelScript 获取并删除key所关联的值lua脚本
	getAndDelScript = `local current = redis.call('GET', KEYS[1]);
	if (current) then
		redis.call('DEL', KEYS[1]);
	end
	return current;`
)

type Store struct {
	kv.Store
	Redis *redis.Redis
}

func InitRedis(cfg *config.Config) (*Store, error) {
	if cfg.Kv.Redis == nil {
		return nil, errors.New("redis config is nil")
	}

	var kvConf kv.KvConf
	for _, con := range cfg.Kv.Redis {
		kvConf = append(kvConf, cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: con.Host,
				Type: con.Type,
				Pass: con.Pass,
			},
			Weight: 1,
		})
	}

	rd := redis.MustNewRedis(kvConf[0].RedisConf)
	store := &Store{
		Store: kv.NewStore(kvConf),
		Redis: rd,
	}

	logger.Info("Redis connected successfully")
	return store, nil
}
