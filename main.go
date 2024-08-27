package main

import (
	"my-qqbot/config"
	"my-qqbot/hub"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	// 注册插件
	hub.Register()
	zero.RunAndBlock(&zero.Config{
		NickName:      config.Conf.NickName,
		CommandPrefix: "/",
		SuperUsers:    []int64{config.Conf.Self},
		Driver: []zero.Driver{
			// 正向 WS
			driver.NewWebSocketClient(config.Conf.WsDriver.Address, config.Conf.WsDriver.Token),
		},
	}, nil)
}
