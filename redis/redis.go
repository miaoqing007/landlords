package redis

import (
	"landlords/enmu"
	"landlords/helper/redisgo"
	"github.com/golang/glog"
)

var _redisCacher *redisgo.Cacher

func InitRedis() {
	var err error
	_redisCacher, err = redisgo.New(redisgo.Options{
		Addr:     enmu.ServerHost + ":" + enmu.RedisPort,
		Db:       0,
		Password: "",
		Network:  "tcp",
	})
	if err != nil {
		panic("初始化redis失败")
		return
	}
	glog.Info("初始化redis完成")
}

func HGetAll(key string, value interface{}) error {
	return _redisCacher.HGetAll(key, value)
}

func HGet(key, field string) (interface{}, error) {
	return _redisCacher.HGet(key, field)
}

func Get(key string) (interface{}, error) {
	return _redisCacher.Get(key)
}

func Set(key string, value interface{}, expire int64) error {
	return _redisCacher.Set(key, value, expire)
}

func HSet(key, field string, val interface{}) (interface{}, error) {
	return _redisCacher.HSet(key, field, val)
}

func HMSet(key string, value interface{}, expire int) error {
	return _redisCacher.HMSet(key, value, expire)
}

func Exists(key string) bool {
	ok, _ := _redisCacher.Exists(key)
	return ok
}

func HExists(key string, val interface{}) bool {
	ok, _ := _redisCacher.HExists(key, val)
	return ok
}
