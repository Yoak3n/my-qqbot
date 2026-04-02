package chat

import (
	"my-qqbot/internal/model"
)

func Ask(from model.From, question string) {
	conversation, ok := GlobalConversationHub().Listener[from]
	if !ok {
		conversation = NewConversation(from)
		GlobalConversationHub().Listener[from] = conversation
	}
	conversation.AddMessage(question)
	GlobalConversationHub().queue <- conversation
}

func Reset(from *model.From) bool {
	delete(GlobalConversationHub().Listener, *from)
	return true
}
