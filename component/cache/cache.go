package cache

import (
	"MetaFarmBackend/component/redis"
	"context"
	"encoding/json"

	"reflect"

	"github.com/pkg/errors"

	"MetaFarmBackend/component/convert"
)

const (
	// getAndDelScript 获取并删除key所关联的值lua脚本
	getAndDelScript = `local current = redis.call('GET', KEYS[1]);
	if (current) then
		redis.call('DEL', KEYS[1]);
	end
	return current;`
)

type CacheService struct {
	store *redis.Store
}

func NewCacheService(store *redis.Store) *CacheService {
	return &CacheService{store: store}
}

// 提供通用缓存方法
func (c *CacheService) GetWithCtx(ctx context.Context, key string) (string, error) {
	value, err := c.store.GetCtx(ctx, key)
	return value, err
}

// GetInt 返回给定key所关联的int值
func (c *CacheService) GetInt(key string) (int, error) {
	value, err := c.store.Get(key)
	if err != nil {
		return 0, err
	}

	return convert.ToInt(value), nil
}

// SetInt 将int value关联到给定key，seconds为key的过期时间（秒）
func (c *CacheService) SetInt(key string, value int, seconds ...int) error {
	return c.SetString(key, convert.ToString(value), seconds...)
}

// GetInt64 返回给定key所关联的int64值
func (c *CacheService) GetInt64(key string) (int64, error) {
	value, err := c.store.Get(key)
	if err != nil {
		return 0, err
	}

	return convert.ToInt64(value), nil
}

// SetInt64 将int64 value关联到给定key，seconds为key的过期时间（秒）
func (c *CacheService) SetInt64(key string, value int64, seconds ...int) error {
	return c.SetString(key, convert.ToString(value), seconds...)
}

// GetBytes 返回给定key所关联的[]byte值
func (c *CacheService) GetBytes(key string) ([]byte, error) {
	value, err := c.store.Get(key)
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}

// GetDel 返回并删除给定key所关联的string值
func (c *CacheService) GetDel(key string) (string, error) {
	resp, err := c.store.Eval(getAndDelScript, key)
	if err != nil {
		return "", errors.Wrap(err, "eval script err")
	}

	return convert.ToString(resp), nil
}

// SetString 将string value关联到给定key，seconds为key的过期时间（秒）
func (c *CacheService) SetString(key, value string, seconds ...int) error {
	if len(seconds) != 0 {
		return errors.Wrapf(c.store.Setex(key, value, seconds[0]), "setex by seconds = %v err", seconds[0])
	}

	return errors.Wrap(c.store.Set(key, value), "set err")
}

// Read 将给定key所关联的值反序列化到obj对象
// 返回false时代表给定key不存在
func (c *CacheService) Read(key string, obj interface{}) (bool, error) {
	if !isValid(obj) {
		return false, errors.New("obj is invalid")
	}

	value, err := c.GetBytes(key)
	if err != nil {
		return false, errors.Wrap(err, "get bytes err")
	}
	if len(value) == 0 {
		return false, nil
	}

	err = json.Unmarshal(value, obj)
	if err != nil {
		return false, errors.Wrap(err, "json unmarshal value to obj err")
	}

	return true, nil
}

// Write 将对象obj序列化后关联到给定key，seconds为key的过期时间（秒）
func (c *CacheService) Write(key string, obj interface{}, seconds ...int) error {
	value, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "json marshal obj err")
	}

	return c.SetString(key, string(value), seconds...)
}

// GetFunc 给定key不存在时调用的数据获取函数
type GetFunc func() (interface{}, error)

// ReadOrGet 将给定key所关联的值反序列化到obj对象
// 若给定key不存在则调用数据获取函数，调用成功时赋值至obj对象
// 并将其序列化后关联到给定key，seconds为key的过期时间（秒）
func (c *CacheService) ReadOrGet(key string, obj interface{}, gf GetFunc, seconds ...int) error {
	isExist, err := c.Read(key, obj)
	if err != nil {
		return errors.Wrap(err, "read obj by err")
	}

	if !isExist {
		data, err := gf()
		if err != nil {
			return err
		}

		if !isValid(data) {
			return errors.New("get data is invalid")
		}

		ov, dv := reflect.ValueOf(obj).Elem(), reflect.ValueOf(data).Elem()
		if ov.Type() != dv.Type() {
			return errors.New("obj type and get data type are not equal")
		}
		ov.Set(dv)

		_ = c.Write(key, data, seconds...)
	}

	return nil
}

// isValid 判断对象是否合法
func isValid(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return false
	}

	return true
}
