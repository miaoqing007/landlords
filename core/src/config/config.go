package config

import (
	"github.com/golang/glog"
	"gopkg.in/ini.v1"
)

var (
	GameIp, GamePort   string
	RedisIp, RedisPort string
)

func InitConfig() {
	cfg, err := ini.Load("gameconfig.ini")
	if err != nil {
		glog.Info("load config", err)
		return
	}

	// 获取默认分区的key
	GameIp = cfg.Section("GameAddr").Key("Ip").String()
	GamePort = cfg.Section("GameAddr").Key("Port").String()

	RedisIp = cfg.Section("RedisAddr").Key("Ip").String()
	RedisPort = cfg.Section("RedisAddr").Key("Port").String()
}
