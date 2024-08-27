package main

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"my-qqbot/config"
)

func init() {

}

func main() {
	// 注册插件
	zero.OnCommand("hello").
		Handle(func(ctx *zero.Ctx) {
			ctx.Send("world")
		})
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
