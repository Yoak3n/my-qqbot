package command

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"my-qqbot/package/logger"
	"my-qqbot/plugin/bilibili"
	"net/http"
	"strings"
	"time"
)

const HotWord = "https://s.search.bilibili.com/main/hotword"

type (
	Event struct {
		Name    string
		Action  func(from *bilibili.From)
		From    *bilibili.From
		Timer   *time.Timer
		running bool
	}

	EventHub struct {
		Pool  []*Event
		Begin bool
	}
)

var (
	hub *EventHub
)

func dailyHotNews(from *bilibili.From) {
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
	hottestVideo := bilibili.SearchVideoFromKeyword(hottestWord)
	// 今日热搜前十
	now := time.Now().Format("1月2日")
	content := now + "热搜:\n"
	for index, word := range list {
		content += fmt.Sprintf("%d. %s\n", index+1, word.Get("show_name").String())
	}
	content, _ = strings.CutSuffix(content, "\n")
	notification := &bilibili.Notification{
		Private: from.Private,
		Target:  from.Id,
		Message: content,
	}
	bilibili.Notify <- notification

	message := fmt.Sprintf("今日热词: %s\n"+
		"推荐视频：%s\n"+
		"https://bilibili.com/video/%s\n"+
		"视频简介:%s", hottestWord, hottestVideo.Title, hottestVideo.BVID, hottestVideo.Description)
	pics := make([]string, 0)
	pics = append(pics, hottestVideo.Cover)
	notify := &bilibili.Notification{
		Private: from.Private,
		Target:  from.Id,
		Message: message,
		Picture: pics,
	}
	bilibili.Notify <- notify
}

func addNewsSub(from *bilibili.From) {
	now := time.Now()
	lunch := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	d := lunch.Sub(now)
	go dailyHotNews(from)
	if d < 0 {
		d += 24 * time.Hour
	}
	event := &Event{Name: "dailyHotNews", Action: dailyHotNews, Timer: time.NewTimer(d), From: from, running: false}
	if hub == nil {
		hub = &EventHub{
			Pool: make([]*Event, 0),
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

func runEventCircle() {
	for {
		for _, event := range hub.Pool {
			if event.Timer != nil && !event.running {
				go func() {
					if event.running {
						return
					}
					event.running = true
					<-event.Timer.C
					event.Action(event.From)
					event.Timer.Reset(24 * time.Hour)
					event.running = false
				}()
			}
		}
	}
}
