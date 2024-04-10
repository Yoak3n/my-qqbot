package command

import (
	"fmt"
	"my-qqbot/config"
	"my-qqbot/plugin/bilibili"
	"my-qqbot/plugin/chat"
	"os"
	"strconv"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func Register() {
	zero.OnCommand("登录哔哩哔哩").Handle(pluginMap["登录哔哩哔哩"])
	zero.OnRegex(`^/启用订阅直播间.*?(\d+)(?:\D+(\d+))*`).Handle(pluginMap["订阅直播间"])
	zero.OnRegex(`^/订阅直播间.*?(\d+)`).Handle(pluginMap["订阅直播间"])
	zero.OnRegex(`^/订阅动态.*?(\S+)`).Handle(pluginMap["订阅动态"])
	zero.OnRegex(`^/取消订阅.*?(\S+)`).Handle(pluginMap["取消订阅"])
	zero.OnCommandGroup([]string{"help", "帮助"}).Handle(pluginMap["帮助"])
	zero.OnCommandGroup([]string{"重置", "重置对话", "重置会话", "reset", "Reset"}).Handle(pluginMap["重置对话"])
	zero.OnCommand("test", zero.OnlyGroup).Handle(func(ctx *zero.Ctx) {
	})
	zero.OnMessage().Handle(pluginMap["ai对话"])
	go SendMessage(config.Conf.Self)
}

func SendMessage(self ...int64) {
	for notify := range bilibili.Notify {
		for _, item := range self {
			bot := zero.GetBot(item)
			var chain []message.MessageSegment
			if len(notify.Picture) == 0 {
				chain = append(chain, message.Text(notify.Message))
			} else {
				chain = append(chain, message.Text(notify.Message))
				for _, picture := range notify.Picture {
					chain = append(chain, message.Image(picture))
				}
			}
			m := (message.Message)(chain)
			if notify.Private {
				bot.SendPrivateMessage(notify.Target, m)
			} else {
				bot.SendGroupMessage(notify.Target, m)
			}
		}
	}
	// for {
	// 	select {
	// 	case notify := <-bilibili.Notify:
	// 		logger.Logger.Debugln("收到消息：" + notify.Message)
	// 		for _, item := range self {
	// 			bot := zero.GetBot(item)
	// 			var chain []message.MessageSegment
	// 			if len(notify.Picture) == 0 {
	// 				chain = append(chain, message.Text(notify.Message))
	// 			} else {
	// 				chain = append(chain, message.Text(notify.Message))
	// 				for _, picture := range notify.Picture {
	// 					chain = append(chain, message.Image(picture))
	// 				}
	// 			}
	// 			m := (message.Message)(chain)
	// 			if notify.Private {
	// 				bot.SendPrivateMessage(notify.Target, m)
	// 			} else {
	// 				bot.SendGroupMessage(notify.Target, m)
	// 			}
	// 		}
	// 	}
	// }

}

func loginBili(ctx *zero.Ctx) {
	go bilibili.GetCookie()
	ctx.Send("正在获取哔哩哔哩Cookie，请稍等...")
	for {
		msg := <-bilibili.Scan
		if msg == "scan" {
			data, _ := os.ReadFile("qrcode.png")
			ctx.SendChain(message.Text("请使用哔哩哔哩App扫描二维码"), message.ImageBytes(data))
		} else if msg == "done" {
			ctx.Send("登录成功！")
			return
		} else {
			ctx.Send("登录失败，请重新尝试登录！")
			return
		}

	}
}

func listenBili(ctx *zero.Ctx) {
	targets := ctx.State["regex_matched"].([]string)[1:]
	//bilibili.NewLiveRoomPlugin()
	var roomsID []int
	for _, target := range targets {
		id, err := strconv.Atoi(target)
		if err != nil {
			ctx.Send("直播间ID错误！")
		}
		roomsID = append(roomsID, id)
	}
	from := &bilibili.From{}
	if ctx.Event.MessageType == "group" {
		from = &bilibili.From{
			Id:      ctx.Event.GroupID,
			Private: false,
		}
	} else {
		from = &bilibili.From{
			Id:      ctx.Event.UserID,
			Private: true,
		}
	}
	err := bilibili.AddSub(from, roomsID...)
	if err != nil {
		ctx.Send("订阅失败：" + err.Error())
	}
	if len(targets) == 1 {
		ctx.Send("订阅" + targets[0] + "直播间成功！")
	} else if len(targets) > 1 {
		ctx.SendChain(message.Text(fmt.Sprintf("订阅b站直播间%s成功！", strings.Join(targets, ","))))
	}

}

func listenDynamic(ctx *zero.Ctx) {
	targets := ctx.State["regex_matched"].([]string)[1:]
	handleTargets := make([]string, 0)
	for _, i := range targets {
		s := strings.TrimSpace(i)
		exist := false
		for _, j := range handleTargets {
			if s == j {
				exist = true
				break
			}
		}
		if !exist && s != "" {
			handleTargets = append(handleTargets, s)
		}
	}
	if config.Conf.Cookie == "" {
		ctx.Send("请先【/登录哔哩哔哩】获取cookie")
		return
	}
	from := &bilibili.From{}
	if ctx.Event.MessageType == "group" {
		from = &bilibili.From{
			Id:      ctx.Event.GroupID,
			Private: false,
		}
	} else {
		from = &bilibili.From{
			Id:      ctx.Event.UserID,
			Private: true,
		}
	}
	bilibili.AddDynamic(*from, handleTargets[0])
	if len(handleTargets) == 1 {
		ctx.Send("订阅" + handleTargets[0] + "动态成功！")
	} else {
		t := strings.Join(handleTargets, ",")
		ctx.Send("订阅" + t + "动态成功！")
	}
}

func resetConversation(ctx *zero.Ctx) {
	ok := false
	if ctx.Event.MessageType == "private" {
		ok = chat.Reset(ctx.Event.UserID)
	} else {
		ok = chat.Reset(ctx.Event.GroupID)
	}
	if ok {
		ctx.SendChain(message.Text("重置对话成功"))
	} else {
		ctx.SendChain(message.Text("重置对话失败"))
	}
}

func aiChat(ctx *zero.Ctx) {
	if strings.HasPrefix(ctx.Event.RawMessage, "/") {
		return
	}
	r := &chat.ResponseBody{}
	if ctx.Event.MessageType == "group" {
		r = chat.Ask(ctx.Event.GroupID, ctx.Event.Message.String())
	} else {
		r = chat.Ask(ctx.Event.UserID, ctx.Event.Message.String())
	}
	ctx.SendChain(message.Text(r.Answer))
}

func cancelListenDynamic(ctx *zero.Ctx) {
	targets := ctx.State["regex_matched"].([]string)[1:]
	handleTargets := make([]string, 0)
	for _, i := range targets {
		s := strings.TrimSpace(i)
		exist := false
		for _, j := range handleTargets {
			if s == j {
				exist = true
				break
			}
		}
		if !exist && s != "" {
			handleTargets = append(handleTargets, s)
		}
	}
	if config.Conf.Cookie == "" {
		ctx.Send("请先【/登录哔哩哔哩】获取cookie")
		return
	}
	from := &bilibili.From{}
	if ctx.Event.MessageType == "group" {
		from = &bilibili.From{
			Id:      ctx.Event.GroupID,
			Private: false,
		}
	} else {
		from = &bilibili.From{
			Id:      ctx.Event.UserID,
			Private: true,
		}
	}
	err := bilibili.CancelDynamic(*from, handleTargets[0])
	if err != nil {
		ctx.Send(err.Error())
	}
	ctx.Send("取消订阅" + handleTargets[0] + "的动态成功！")

}
