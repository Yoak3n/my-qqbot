package deep_seek

type ChatCompletionNewParams struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}
type Role string

const (
	UserRole      = Role("user")
	AssistantRole = Role("assistant")
	SystemRole    = Role("system")
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func NewChatCompletionNewParams(model string) *ChatCompletionNewParams {
	return &ChatCompletionNewParams{
		Messages: make([]Message, 0),
		Model:    model,
		Stream:   false,
	}
}
