package command

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
)

var PluginMap map[string]func(ctx *zero.Ctx)

func init() {
	PluginMap = make(map[string]func(ctx *zero.Ctx))
	PluginMap["帮助"] = help
	PluginMap["登录哔哩哔哩"] = loginBili
	PluginMap["订阅直播间"] = listenBili
	PluginMap["订阅动态"] = listenDynamic
	PluginMap["ai对话"] = aiChat
	PluginMap["重置对话"] = resetConversation

}

func help(ctx *zero.Ctx) {
	msgArr := make([]string, 0)
	for k, _ := range PluginMap {
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
	ctx.SendChain(message.Text("帮助信息：\n" + strings.Join(msgLines, "")))
}
