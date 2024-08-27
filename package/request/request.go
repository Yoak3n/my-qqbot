package request

import (
	"bytes"
	"fmt"
	"io"
	"my-qqbot/config"
	"my-qqbot/package/logger"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func Get(urlStr string, args ...string) ([]byte, error) {
	params := ""
	if l := len(args); l > 0 {
		params = "?" + url.PathEscape(strings.Join(args, "&"))
	} else {
		params = ""
	}
	client := &http.Client{}
	logger.Logger.Info("请求地址: " + params)
	req, err := http.NewRequest(http.MethodGet, params, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0")
	req.Header.Set("Cookie", config.Conf.Bilibili.Cookie)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func FetchArgs(s string) error {
	re, err := regexp.Compile(`^订阅动态.*?(\S+)(.*?(\S+))*`)
	if err != nil {
		return err
	}
	result := re.FindAllStringSubmatch(s, -1)
	for _, r := range result {
		fmt.Println(r)
		for _, i := range r {
			fmt.Println(i)
		}
	}
	return nil
}

func Post(uri string, data []byte) []byte {
	body := bytes.NewBuffer(data)
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	return buf
}
