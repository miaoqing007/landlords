package redisgo

import (
	"app/enmu"
	"github.com/aiscrm/redisgo"
	"github.com/golang/glog"
)

var _redisCacher *redisgo.Cacher

func InitRedis() {
	var err error
	_redisCacher, err = redisgo.New(redisgo.Options{
		Addr:     enmu.ServerHost + ":" + enmu.RedisPort,
		Db:       0,
		Password: "",
	})
	if err != nil {
		glog.Info("初始化Redis失败")
		return
	}
	glog.Info("初始化Redis完成")
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

func HSet(key, field string, expire int64) (interface{}, error) {
	return _redisCacher.HSet(key, field, expire)
}

func HMSet(key string, value interface{}, expire int) error {
	return _redisCacher.HMSet(key, value, expire)
}
