package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"strings"
)

type (
	Configuration struct {
		Self     int64
		NickName []string
		Bilibili Bilibili
		WsDriver WsDriver
		AIChat   AIChat
	}
	Bilibili struct {
		Cookie       string
		RefreshToken string
	}
	WsDriver struct {
		Address string
		Token   string
	}
	AIChat struct {
		BaseUrl string
		Model   string
		Key     string
	}
)

var (
	k    *koanf.Koanf
	Conf *Configuration
)

func init() {
	k = koanf.New(".")
	configPath := "config.yaml"
	err := k.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		panic(err)
	}
	Conf = &Configuration{
		Bilibili: Bilibili{},
		NickName: make([]string, 3),
	}
	loadToConfiguration()
}
func loadToConfiguration() {
	Conf.Self = k.Int64("self")
	Conf.NickName = k.Strings("nickname")
	Conf.Bilibili.Cookie = k.String("bilibili.cookie")
	Conf.Bilibili.RefreshToken = k.String("bilibili.refresh_token")

	Conf.WsDriver.Address = k.String("ws.address")
	Conf.WsDriver.Token = k.String("ws.token")

	Conf.AIChat.BaseUrl = k.String("ai_chat.base_url")
	Conf.AIChat.Model = k.String("ai_chat.model")
	Conf.AIChat.Key = k.String("ai_chat.key")

	Conf.AIChat.BaseUrl = strings.TrimRight(Conf.AIChat.BaseUrl, "/")

}
func UpdateBilibiliCookie(cookie string, refreshToken string) {
	Conf.Bilibili.Cookie = cookie
	err := k.Set("bilibili.cookie", cookie)
	err = k.Set("bilibili.refresh_token", refreshToken)
	if err != nil {
		fmt.Println(err)
	}
	buf, err := k.Marshal(yaml.Parser())
	if err != nil {
		return
	}
	fp, err := os.OpenFile("config.yaml", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	defer fp.Close()
	if err != nil {
		return
	}
	_, err = fp.Write(buf)
	if err != nil {
		return
	}
}
