package bilibili

import (
	"github.com/tidwall/gjson"
	"my-qqbot/model"
	"my-qqbot/package/request"
)

const VideoSearch = "https://api.bilibili.com/x/web-interface/wbi/search/type?search_type=video&keyword="

func SearchVideoFromKeyword(keyword string) *model.Video {
	res, err := request.Get(VideoSearch + keyword)
	if err != nil {
		return nil
	}
	result := gjson.ParseBytes(res)
	firstVideo := result.Get("data.result.0")

	video := &model.Video{
		AID:         firstVideo.Get("aid").Int(),
		BVID:        firstVideo.Get("bvid").String(),
		Title:       firstVideo.Get("title").String(),
		Description: firstVideo.Get("description").String(),
		Author:      firstVideo.Get("author").String(),
		Cover:       "https:" + firstVideo.Get("pic").String(),
	}

	return video
}
