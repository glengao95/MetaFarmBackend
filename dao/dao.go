package dao

import (
	"context"

	"MetaFarmBackend/component/redis"

	"gorm.io/gorm"
)

// Dao is show dao.
type Dao struct {
	ctx     context.Context
	DB      *gorm.DB
	KvStore *redis.Store
}

func NewDao(ctx context.Context, db *gorm.DB, kvStore *redis.Store) *Dao {
	return &Dao{
		ctx:     ctx,
		DB:      db,
		KvStore: kvStore,
	}
}
