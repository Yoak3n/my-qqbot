package hub

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
)

var pluginMap map[string]func(ctx *zero.Ctx)

func init() {
	pluginMap = make(map[string]func(ctx *zero.Ctx))
	pluginMap["帮助"] = help
	pluginMap["登录哔哩哔哩"] = loginBili
	pluginMap["订阅直播间"] = listenBili
	pluginMap["订阅动态"] = listenDynamic
	pluginMap["ai对话"] = aiChat
	pluginMap["重置对话"] = resetConversation
	pluginMap["取消订阅"] = cancelListenDynamic
	pluginMap["每日新闻"] = subDailyNews
	pluginMap["取消每日新闻"] = cancelDailyNews
}

func help(ctx *zero.Ctx) {
	msgArr := make([]string, 0)
	for k := range pluginMap {
		msgArr = append(msgArr, k)
	}

	msgLines := make([]string, 0)

	for i := 0; i < len(msgArr); i++ {
		msg := fmt.Sprintf("%d.%s", i+1, msgArr[i])
		if i != len(msgArr)-1 {
			msg = msg + "\n"
		}
		msgLines = append(msgLines, msg)
	}
	ctx.SendChain(message.Text("命令列表：\n" + strings.Join(msgLines, "")))
}
