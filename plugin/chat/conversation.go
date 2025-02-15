package chat

import (
	"context"
	"fmt"
	"my-qqbot/config"
	"my-qqbot/internal/model"
	"my-qqbot/internal/queue"
	"my-qqbot/package/deep_seek"
	"my-qqbot/package/logger"
	"strings"
)

type (
	ConversationHub struct {
		Listener map[model.From]*Conversation
		queue    chan *Conversation
		started  bool
		client   *deep_seek.Client
	}
	Conversation struct {
		Ctx   context.Context
		Param *deep_seek.ChatCompletionNewParams
		From  model.From
	}
)

var ConversationHubInstance *ConversationHub

func NewConversationHub() *ConversationHub {
	c := &ConversationHub{
		Listener: make(map[model.From]*Conversation),
		started:  false,
		queue:    make(chan *Conversation, 3),
		client:   deep_seek.NewClient(),
	}
	c.client.SetBaseUrl(config.Conf.AIChat.BaseUrl)
	c.client.SetAPIKey(config.Conf.AIChat.Key)
	go c.Start()
	return c
}
func (c *ConversationHub) Start() {
	if c.started {
		return
	}
	c.started = true
	for {
		select {
		case con := <-c.queue:
			completion, err := c.client.ChatCompletion(con.Ctx, *con.Param)
			if err != nil {
				con.Reply(err.Error())
				logger.Logger.Error(err)
				return
			}
			if completion == nil {
				con.Reply("请求失败")
				logger.Logger.Error("请求失败")
				return
			}
			if len(completion.Choices) > 0 {
				answer := completion.Choices[0].Message.Content
				con.UpdateAssistantMessage(answer)
				// 兼容硅基流动的推理模型名
				if strings.HasSuffix(config.Conf.AIChat.Model, "reasoner") || strings.HasSuffix(config.Conf.AIChat.Model, "R1") {
					reason := completion.Choices[0].Message.ReasoningContent
					if reason != "" {
						con.Reply(fmt.Sprintf("推理过程：\n%s", reason))
					}
				}

				con.Reply(answer)
			}

		}
	}
}

func NewConversation(from model.From) *Conversation {
	return &Conversation{
		Ctx: context.Background(),
		Param: &deep_seek.ChatCompletionNewParams{
			Model:    config.Conf.AIChat.Model,
			Messages: make([]deep_seek.Message, 0),
		},
		From: from,
	}
}

func (c *Conversation) AddMessage(msg string) {
	m := deep_seek.Message{
		Role:    deep_seek.UserRole,
		Content: msg,
	}
	c.Param.Messages = append(c.Param.Messages, m)
}
func (c *Conversation) UpdateAssistantMessage(reply string) {
	m := deep_seek.Message{
		Role:    deep_seek.AssistantRole,
		Content: reply,
	}
	c.Param.Messages = append(c.Param.Messages, m)
}
func (c *Conversation) Reply(reply string) {
	notify := &model.Notification{
		Private: c.From.Private,
		Target:  c.From.Id,
		Message: reply,
	}
	queue.Notify <- notify
}
