package bilibili

import (
	"github.com/tidwall/gjson"
	"my-qqbot/model"
	"my-qqbot/package/logger"
	"my-qqbot/package/request"
	"strings"
)

const VideoSearch = "https://api.bilibili.com/x/web-interface/wbi/search/type?search_type=video&keyword="

func SearchVideoFromKeyword(keyword string) *model.Video {
	logger.Logger.Println("Searching video from keyword: " + keyword)

	res, err := request.Get(VideoSearch + keyword)
	logger.Logger.Println(res)
	if err != nil {
		logger.Logger.Println(err)
		return nil
	}
	result := gjson.ParseBytes(res)
	firstVideo := result.Get("data.result.0")
	logger.Logger.Println(firstVideo)
	title := firstVideo.Get("title").String()
	title = strings.Replace(title, "<em class=\"keyword\">", "", -1) // remove keyword highlight
	title = strings.Replace(title, "</em>", "", -1)
	video := &model.Video{
		AID:         firstVideo.Get("aid").Int(),
		BVID:        firstVideo.Get("bvid").String(),
		Title:       title,
		Description: firstVideo.Get("description").String(),
		Author:      firstVideo.Get("author").String(),
		Cover:       "http:" + firstVideo.Get("pic").String(),
	}
	return video
}
