package chat

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"my-qqbot/config"
	"my-qqbot/model"
	"my-qqbot/queue"
)

type (
	ConversationHub struct {
		Listener map[*model.From]*Conversation
		queue    chan *Conversation
		started  bool
		client   *openai.Client
	}
	Conversation struct {
		Ctx   context.Context
		Param *openai.ChatCompletionNewParams
		From  *model.From
	}
)

var ConversationHubInstance *ConversationHub

func NewConversationHub() *ConversationHub {
	c := &ConversationHub{
		Listener: make(map[*model.From]*Conversation),
		started:  false,
		queue:    make(chan *Conversation, 3),
		client: openai.NewClient(
			option.WithBaseURL(config.Conf.AIChat.BaseUrl),
			option.WithAPIKey(config.Conf.AIChat.Key), // defaults to os.LookupEnv("OPENAI_API_KEY")
		),
	}
	go c.Start()
	return c
}
func (c *ConversationHub) Start() {
	if c.started {
		return
	}
	for {
		select {
		case con := <-c.queue:
			completion, err := c.client.Chat.Completions.New(con.Ctx, *con.Param)
			if err != nil {
				con.Reply(err.Error())
				return
			}
			answer := completion.Choices[0].Message.Content
			con.UpdateAssistantMessage(answer)
			con.Reply(answer)
		}
	}
}

func NewConversation(from *model.From) *Conversation {
	return &Conversation{
		Ctx: context.Background(),
		Param: &openai.ChatCompletionNewParams{
			Model:    openai.F(config.Conf.AIChat.Model),
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{}),
		},
		From: from,
	}
}

func (c *Conversation) AddMessage(msg string) {
	c.Param.Messages.Value = append(c.Param.Messages.Value, openai.UserMessage(msg))
}
func (c *Conversation) UpdateAssistantMessage(reply string) {
	c.Param.Messages.Value = append(c.Param.Messages.Value, openai.AssistantMessage(reply))
}
func (c *Conversation) Reply(reply string) {
	notify := &model.Notification{
		Private: c.From.Private,
		Target:  c.From.Id,
		Message: reply,
	}
	queue.Notify <- notify
}
