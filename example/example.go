package main

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"my-qqbot/command"
	"my-qqbot/config"
)

func init() {

}

func main() {
	// 注册插件
	command.Register()
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{123456},
		Driver: []zero.Driver{
			// 正向 WS
			driver.NewWebSocketClient(config.Conf.WsDriver.Address, config.Conf.WsDriver.Token),
		},
	}, nil)
}
