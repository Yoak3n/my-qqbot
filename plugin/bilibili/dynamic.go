package bilibili

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"my-qqbot/config"
	"my-qqbot/package/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TextDynamic    = "DYNAMIC_TYPE_WORD"
	DrawDynamic    = "DYNAMIC_TYPE_DRAW"
	VideoDynamic   = "DYNAMIC_TYPE_AV"
	ForwardDynamic = "DYNAMIC_TYPE_FORWARD"
)

var dynamicHub *DynamicListener

type (
	Dynamic struct {
		Name      string
		UId       int64
		Text      string
		Picture   []string
		Timestamp int64
		Type      string
		Id        string
		Extra     string
	}
	DynamicListener struct {
		Listener  map[From][]string
		Listening bool
		done      []string
	}
)

func init() {
	dynamicHub = &DynamicListener{
		Listener:  make(map[From][]string),
		Listening: false,
		done:      make([]string, 0),
	}
}

func AddDynamic(origin From, target ...string) {
	if dynamicHub == nil {
		dynamicHub = &DynamicListener{
			Listener: make(map[From][]string),
		}
	}
	dynamicHub.Listener[origin] = append(dynamicHub.Listener[origin], target...)
	go GetDynamicListLoop()
}

func getDynamicList(baseline string) ([]Dynamic, error) {
	api := "https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/all?type=all"
	if baseline != "" {
		api += "&update_baseline=" + baseline
	}
	// 调用API获取动态列表数据
	// ...
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", api, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", config.Conf.Cookie)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	dynamicList := make([]Dynamic, 0)
	result := gjson.ParseBytes(body)
	// 解析返回的数据
	// ...
	if code := result.Get("code").Int(); code != 0 {
		return nil, fmt.Errorf("API error: %s", result.Get("message").String())
	}
	items := result.Get("data.items").Array()
	for _, item := range items {
		// 如果是上次已经获取过的动态，则跳过
		if id := item.Get("id_str").String(); id == baseline {
			logger.Logger.Debugln("Skipping dynamic:", id)
			break
		}
		typ := item.Get("type").String()
		switch typ {
		case DrawDynamic:
			pics := make([]string, 0)
			for _, pic := range item.Get("modules.module_dynamic.major.draw.items").Array() {
				pics = append(pics, pic.Get("src").String())
			}
			dynamic := Dynamic{
				Name:      item.Get("modules.module_author.name").String(),
				UId:       item.Get("modules.module_author.mid").Int(),
				Timestamp: item.Get("modules.module_author.pub_ts").Int(),
				Text:      item.Get("modules.module_dynamic.desc.text").String(),
				Id:        item.Get("id_str").String(),
				Picture:   pics,
				Type:      typ,
			}
			dynamicList = append(dynamicList, dynamic)
		case TextDynamic:
			dynamic := Dynamic{
				Name:      item.Get("modules.module_author.name").String(),
				UId:       item.Get("modules.module_author.mid").Int(),
				Timestamp: item.Get("modules.module_author.pub_ts").Int(),
				Id:        item.Get("id_str").String(),
				Text:      item.Get("modules.module_dynamic.desc.text").String(),
				Type:      typ,
			}
			dynamicList = append(dynamicList, dynamic)
		case VideoDynamic:
			video := item.Get("modules.module_dynamic.major.archive.cover").String()
			dynamic := Dynamic{
				Name:      item.Get("modules.module_author.name").String(),
				UId:       item.Get("modules.module_author.mid").Int(),
				Timestamp: item.Get("modules.module_author.pub_ts").Int(),
				Text:      item.Get("modules.module_dynamic.major.archive.bvid").String(),
				Id:        item.Get("id_str").String(),
				Extra:     item.Get("modules.module_dynamic.major.archive.title").String() + `%%` + item.Get("modules.module_dynamic.major.archive.desc").String(),
				Type:      typ,
				Picture:   []string{video},
			}
			dynamicList = append(dynamicList, dynamic)
		case ForwardDynamic:
			dynamic := Dynamic{
				Type:      typ,
				Id:        item.Get("id_str").String(),
				UId:       item.Get("modules.module_author.mid").Int(),
				Name:      item.Get("modules.module_author.name").String(),
				Timestamp: item.Get("modules.module_author.pub_ts").Int(),
				Text:      item.Get("modules.module_dynamic.desc.text").String(),
				Extra:     item.Get("orig.modules.module_author.name").String() + `%%` + item.Get("orig.modules.module_dynamic.desc.text").String(),
			}
			dynamicList = append(dynamicList, dynamic)
		}
	}
	logger.Logger.Debugln(dynamicList)
	return dynamicList, nil
}

func GetDynamicListLoop() {
	if dynamicHub.Listening {
		return
	}
	dynamicHub.Listening = true
	baseline := "" // invalidate baseline
	for {
		dynamicList, err := getDynamicList(baseline)
		if err != nil {
			logger.Logger.Errorln(err)
			err = RefreshCookie()
			if err != nil {
				logger.Logger.Errorln(err)
			}
			continue
		}
		for origin, target := range dynamicHub.Listener {
			for _, dynamic := range dynamicList {
				for _, t := range target {
					i, err := strconv.Atoi(t)
					if err != nil {
						if dynamic.Name == t {
							makeNotification(&origin, &dynamic)
						}
					}
					if dynamic.UId == int64(i) {
						makeNotification(&origin, &dynamic)
					}
				}
			}
		}
		if len(dynamicList) > 0 {
			baseline = dynamicList[0].Id
		}
		time.Sleep(time.Minute * 5)
	}

}

func makeNotification(origin *From, dynamic *Dynamic) {
	notify := &Notification{
		Private: origin.Private,
		Target:  origin.Id,
	}
	for _, i := range dynamicHub.done {
		if dynamic.Id == i {
			return
		}
	}
	switch dynamic.Type {
	case ForwardDynamic:
		arr := strings.SplitN(dynamic.Extra, `%%`, 2)
		text := fmt.Sprintf("@%s 转发了%s动态：\n%s\n原动态内容:\n%s\n", dynamic.Name, arr[0], dynamic.Text, arr[1])
		notify.Message = text
	case DrawDynamic:
		text := fmt.Sprintf("@%s 发布了一条动态：\n%s", dynamic.Name, dynamic.Text)
		notify.Message = text
		notify.Picture = dynamic.Picture
	case VideoDynamic:
		arr := strings.SplitN(dynamic.Extra, `%%`, 2)
		text := fmt.Sprintf("@%s 投稿了视频：\n《%s》\n"+
			"视频链接：https://www.bilibili.com/video/%s\n"+
			"视频简介：%s\n"+
			"视频封面：", dynamic.Name, arr[0], dynamic.Text, arr[1])
		notify.Message = text
		notify.Picture = dynamic.Picture
	case TextDynamic:
		text := fmt.Sprintf("@%s 发布了一条动态：\n%s", dynamic.Name, dynamic.Text)
		notify.Message = text
	}
	dynamicHub.done = append(dynamicHub.done, dynamic.Id)
	Notify <- notify

}
