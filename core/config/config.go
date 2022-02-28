package config

import (
	"gopkg.in/ini.v1"
	"log"
)

var (
	OnlineGRPCPort, PvpGRPCPort, GatewayTCPPort, RedisAddr string
)

func init() {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Printf("load config", err)
		return
	}
	online := cfg.Section("Online")
	if online != nil {
		OnlineGRPCPort = online.Key("GRPCPort").String()
	}
	pvp := cfg.Section("Pvp")
	if pvp != nil {
		PvpGRPCPort = pvp.Key("GRPCPort").String()
	}
	gateway := cfg.Section("Gateway")
	if gateway != nil {
		GatewayTCPPort = gateway.Key("TcpPort").String()
	}
	redis := cfg.Section("Redis")
	if redis != nil {
		RedisAddr = redis.Key("Addr").String()
	}
}
