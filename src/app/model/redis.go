package model

import (
	"app/enmu"
	"github.com/aiscrm/redisgo"
	"github.com/golang/glog"
)

var redisCacher *redisgo.Cacher

func InitRedis() {
	redisCacher, _ = redisgo.New(redisgo.Options{
		Addr:     enmu.ServerHost + ":" + enmu.RedisPort,
		Db:       0,
		Password: "",
	})
	glog.Info("初始化Redis完成")
}
