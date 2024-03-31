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
		Self   int64
		Cookie string

		RefreshToken string
		WsDriver     WsDriver
		QWen         QWen
	}
	WsDriver struct {
		Address string
		Token   string
	}
	QWen struct {
		Address string
		Token   string
	}
)

var (
	k    *koanf.Koanf
	Conf *Configuration
)

func InitConfig() {
	k = koanf.New(".")
	configPath := "config.yaml"
	err := k.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		panic(err)
	}
	Conf = &Configuration{
		Cookie:       "",
		RefreshToken: "",
	}
	loadToConfiguration()
	fmt.Println("config loaded:", Conf)
}
func loadToConfiguration() {
	Conf.Self = k.Int64("self")

	Conf.Cookie = k.String("bilibili.cookie")
	Conf.RefreshToken = k.String("bilibili.refresh_token")

	Conf.WsDriver.Address = k.String("ws.address")
	Conf.WsDriver.Token = k.String("ws.token")

	Conf.QWen.Address = strings.TrimRight(k.String("qwen.address"), "/")
}
func UpdateBilibiliCookie(cookie string, refreshToken string) {
	Conf.Cookie = cookie
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
