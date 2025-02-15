package deep_seek

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"my-qqbot/package/logger"
	"net/http"
)

type Client struct {
	BaseUrl string
	APIKey  string
}

func NewClient() *Client {
	return &Client{}
}
func (c *Client) SetBaseUrl(url string) {
	c.BaseUrl = url
}
func (c *Client) SetAPIKey(key string) {
	c.APIKey = key
}

func (c *Client) ChatCompletion(ctx context.Context, param ChatCompletionNewParams) (*ResponseBody, error) {
	path := "/chat/completions"
	b, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(b)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	logger.Logger.Println(string(b))
	defer res.Body.Close()
	var body ResponseBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
