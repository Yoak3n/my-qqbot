package chat

import (
	"encoding/json"
	"fmt"
	"my-qqbot/config"
	"my-qqbot/package/request"
	"strconv"
)

func Ask(id int64, question string, option ...string) *ResponseBody {
	AskApi := config.Conf.QWen.Address + "/v1/chat"
	req := &RequestBody{
		Id:      strconv.FormatInt(id, 10),
		Content: question,
	}
	if len(option) > 0 {
		req.Preset = option[0]
	}
	data, err := json.Marshal(req)
	res := request.Post(AskApi, data)
	answer := &ResponseBody{}
	err = json.Unmarshal(res, answer)
	if err != nil {
		return nil
	}
	return answer
}

func Reset(id int64) bool {
	ResetApi := config.Conf.QWen.Address + "/v1/chat/reset"
	_, err := request.Get(fmt.Sprintf("%s/%d", ResetApi, id))
	if err != nil {
		return false
	}
	return true

}
