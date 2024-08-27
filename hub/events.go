package hub

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"my-qqbot/model"
	"my-qqbot/package/logger"
	"my-qqbot/plugin/bilibili"
	"my-qqbot/queue"
	"net/http"
	"strings"
	"time"
)

const HotWord = "https://s.search.bilibili.com/main/hotword"

// 每天12点更新，但计算时间间隔的时区问题仍有待测试
func dailyHotNews(from *model.From) {
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
	notification := &model.Notification{
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
	notify := &model.Notification{
		Private: from.Private,
		Target:  from.Id,
		Message: message,
		Picture: pics,
	}
	queue.Notify <- notify
}

func addNewsSub(from *model.From) {
	now := time.Now().Local()
	// 每天12点更新，但计算时间间隔的时区问题仍有待测试
	lunch := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)
	d := lunch.Sub(now)
	go dailyHotNews(from)
	if d < 0 {
		d += 24 * time.Hour
	}
	event := &model.Event{Name: "dailyHotNews", Action: dailyHotNews, Timer: time.NewTimer(d), From: from, Running: false}
	if hub == nil {
		hub = &model.EventHub{
			Pool: make([]*model.Event, 0),
		}
		hub.Pool = append(hub.Pool, event)
	} else {
		hub.Pool = append(hub.Pool, event)
	}
	if !hub.Begin {
		hub.Begin = true
		go runEventCircle()
	}
}

func cancelNewsSub(from *model.From) {
	for index, event := range hub.Pool {
		if event.Name == "dailyHotNews" && event.From == from {
			event.Timer.Stop()
			hub.Pool = append(hub.Pool[:index], hub.Pool[index+1:]...)
		}
	}
}
