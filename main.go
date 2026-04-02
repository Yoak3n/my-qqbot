package main

import (
	"my-qqbot/config"
	"my-qqbot/internal/hub"
	"github.com/Yoak3n/gulu/logger"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func init() {
	logger.Init()
	// 注册插件
	hub.Register()
}

func main() {
	drivers := make([]zero.Driver, 0)
	if config.Conf.WsDriver.Type == "server" {
		drivers = append(drivers, driver.NewWebSocketServer(30, config.Conf.WsDriver.Address, config.Conf.WsDriver.Token))
	} else {
		drivers = append(drivers, driver.NewWebSocketClient(config.Conf.WsDriver.Address, config.Conf.WsDriver.Token))
	}

	zero.RunAndBlock(&zero.Config{
		NickName:      config.Conf.NickName,
		CommandPrefix: "/",
		SuperUsers:    []int64{config.Conf.Self},
		Driver:        drivers,
	}, nil)
}
