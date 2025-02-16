package hub

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	model2 "my-qqbot/internal/model"
	"my-qqbot/internal/queue"
	"my-qqbot/package/logger"
	"my-qqbot/plugin/bilibili"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const HotWord = "https://s.search.bilibili.com/main/hotword"

// 每天12点更新，但计算时间间隔的时区问题仍有待测试
func dailyHotNews(from *model2.From) {
	res, err := http.Get(HotWord)
	if err != nil || res.StatusCode != 200 {
		e := fmt.Errorf("get hotword failed, err: %v", err)
		logger.Logger.Errorln(e)
		return
	}
	buf, _ := io.ReadAll(res.Body)
	result := gjson.ParseBytes(buf)
	list := result.Get("list").Array()
	hottestWord := list[0].Get("keyword").String()
	// 今日热搜前十
	now := time.Now().Format("1月2日")
	content := now + "热搜:\n"
	for index, word := range list {
		content += fmt.Sprintf("%d. %s\n", index+1, word.Get("show_name").String())
	}
	content, _ = strings.CutSuffix(content, "\n")
	notification := &model2.Notification{
		Private: from.Private,
		Target:  from.Id,
		Message: content,
	}
	queue.Notify <- notification
	hottestVideo := bilibili.SearchVideoFromKeyword(hottestWord)
	if hottestVideo == nil {
		return
	}
	message := ""
	if hottestVideo.Description != "" {
		message = fmt.Sprintf("今日热词: %s\n"+
			"推荐视频：%s\n"+
			"https://bilibili.com/video/%s\n"+
			"视频简介:%s", hottestWord, hottestVideo.Title, hottestVideo.BVID, hottestVideo.Description)
	} else {
		message = fmt.Sprintf("今日热词: %s\n"+
			"推荐视频：%s\n"+
			"https://bilibili.com/video/%s", hottestWord, hottestVideo.Title, hottestVideo.BVID)
	}
	pics := make([]string, 0)
	pics = append(pics, hottestVideo.Cover)
	notify := &model2.Notification{
		Private: from.Private,
		Target:  from.Id,
		Message: message,
		Picture: pics,
	}
	queue.Notify <- notify
}

func addNewsSub(from *model2.From) {
	now := time.Now().Local()
	// 每天12点更新，但计算时间间隔的时区问题仍有待测试
	// time.Local 返回的是系统时区，容器构建中默认使用标准时区
	lunch := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)
	d := lunch.Sub(now)
	go dailyHotNews(from)
	if d < 0 {
		d += 24 * time.Hour
	}
	event := &model2.Event{Name: "dailyHotNews" + strconv.FormatInt(from.Id, 10), Action: dailyHotNews, Timer: time.NewTimer(d), From: from, Running: false}
	if hub == nil {
		hub = &eventHub{
			Pool:  make([]*model2.Event, 0),
			Begin: false,
		}
	}
	exist := false
	for _, e := range hub.Pool {
		if e.Name == event.Name {
			exist = true
			break
		}
	}
	if !exist {
		hub.Pool = append(hub.Pool, event)
	}
	if !hub.Begin {
		hub.Begin = true
		go hub.runEventCircle()
	}
}

func cancelNewsSub(from *model2.From) {
	for index, event := range hub.Pool {
		if event.Name == "dailyHotNews" && event.From == from {
			event.Timer.Stop()
			hub.Pool = append(hub.Pool[:index], hub.Pool[index+1:]...)
		}
	}
}
