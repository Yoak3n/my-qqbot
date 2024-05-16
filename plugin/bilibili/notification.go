package bilibili

import (
	"errors"
	"fmt"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/tidwall/gjson"
	"my-qqbot/config"
	"my-qqbot/model"
	"my-qqbot/package/logger"
	"my-qqbot/package/request"
	"time"
)

type (
	LiveRoomPlugin struct {
		targetRoomId []int
		Listeners    []*Client
	}
	Client struct {
		RoomID    int
		Client    *client.Client
		Room      *model.Room
		Listening bool
		From      *From
	}
	Notification struct {
		Private bool
		Message string
		Target  int64
		Picture []string
	}
	From struct {
		Private bool
		Id      int64
	}
)

var (
	hub    *LiveRoomPlugin
	Notify chan *Notification
)

func init() {
	Notify = make(chan *Notification, 100)
}

func AddSub(origin *From, targets ...int) error {
	if hub == nil {
		var err error
		hub, err = newLiveRoomPlugin(origin, targets...)
		if err != nil {
			if hub == nil {
				return err
			}
		}
	}
	for _, target := range targets {
		listener, err := genListener(target)
		if err != nil {
			return err
		}
		hub.targetRoomId = append(hub.targetRoomId, target)
		hub.Listeners = append(hub.Listeners, listener)
	}
	go hub.listenLiveStart()
	return nil
}

func newLiveRoomPlugin(origin *From, targets ...int) (*LiveRoomPlugin, error) {
	l := &LiveRoomPlugin{
		targetRoomId: make([]int, 0),
		Listeners:    make([]*Client, 0),
	}
	if len(targets) == 0 {
		return nil, errors.New("targets is empty")
	}
	for _, item := range targets {
		listener, err := genListener(item)
		listener.From = origin
		if err != nil {
			return nil, err
		}
		l.targetRoomId = append(l.targetRoomId, item)
		l.Listeners = append(l.Listeners, listener)

	}
	return l, nil
}

func genListener(id int) (*Client, error) {
	if id == 0 {
		return nil, errors.New("room id is empty")
	}
	c := client.NewClient(id)
	cookie := config.Conf.Bilibili.Cookie
	if cookie == "" {
		return nil, errors.New("请先【/登录哔哩哔哩】获取cookie")
	}
	c.SetCookie(cookie)
	listener := &Client{
		Client:    c,
		Listening: false,
		RoomID:    id,
	}
	info, err := getRoomInfo(id)
	if err != nil {
		return nil, errors.New("获取直播间信息失败,直播间可能不存在")
	}
	listener.Room = info
	return listener, nil
}

func (l *LiveRoomPlugin) listenLiveStart() {
	// 监听直播开始事件
	for _, c := range l.Listeners {
		if !c.Listening {
			c.Listening = true
			c.Client.OnLive(func(live *message.Live) {
				// 处理直播开始事件
				info, err := getRoomInfo(live.Roomid)
				if info == nil || err != nil {
					logger.Logger.Errorln("获取直播间信息失败:", err)
					return
				}
				if c.Room == nil {
					logger.Logger.Errorln("直播间信息为空")
				}
				c.Room = info
				msg := fmt.Sprintf("【%s】开始直播了！\n标题：%s\n观看链接：https://live.bilibili.com/%d", info.Name, info.Title, info.ShortId)
				notify := &Notification{
					Private: c.From.Private,
					Target:  c.From.Id,
					Message: msg,
					Picture: []string{info.Cover},
				}
				Notify <- notify
				logger.Logger.Printf("直播开始：%d", live.Roomid)
			})
			err := c.Client.Start()
			if err != nil {
				notify := &Notification{
					Private: c.From.Private,
					Target:  c.From.Id,
					Message: fmt.Sprintf("监听直播间【%d】失败：%s", c.RoomID, err.Error()),
				}
				Notify <- notify
			}
		}
	}
}

func getRoomInfo(id int) (*model.Room, error) {
	res, err := request.Get("https://api.live.bilibili.com/room/v1/Room/get_info", fmt.Sprintf("room_id=%d", id))
	if err != nil {
		return nil, err
	}
	room := &model.Room{
		ShortId: id,
		User:    &model.User{},
	}
	result := gjson.ParseBytes(res)
	if result.Get("code").Int() == 0 {
		room.User.UID = result.Get("data.uid").Int()
		room.LongId = result.Get("data.room_id").Int()
		room.FollowerCount = result.Get("data.attention").Int()
		room.Cover = result.Get("data.user_cover").String()
		room.Title = result.Get("data.title").String()
		user := getUserInfo(room.User.UID)
		if user != nil {
			room.User = user
		}
		return room, nil
	}
	return nil, errors.New("get room information failed")
}

func getUserInfo(uid int64) *model.User {
	// use local database to avoid anti-crawler
	count := 0
	for {
		res, err := request.Get("https://api.bilibili.com/x/web-interface/card", fmt.Sprintf("mid=%d", uid))
		if err != nil {
			continue
		}
		result := gjson.ParseBytes(res)
		if code := result.Get("code"); code.Exists() && code.Int() != 0 {
			count += 1
			if count > 10 {
				return nil
			}
			time.Sleep(time.Second)
			continue
		}
		data := result.Get("data")
		u := &model.User{
			UID:           uid,
			Avatar:        data.Get("card.face").String(),
			Name:          data.Get("card.name").String(),
			Sex:           data.Get("card.sex").String(),
			FollowerCount: data.Get("follower").Int(),
		}
		return u
	}
}
