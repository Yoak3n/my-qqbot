package chat

import (
	"my-qqbot/internal/model"
)

func Ask(from model.From, question string) {
	if ConversationHubInstance == nil {
		ConversationHubInstance = NewConversationHub()
	}
	conversation, ok := ConversationHubInstance.Listener[from]
	if !ok {
		conversation = NewConversation(from)
		ConversationHubInstance.Listener[from] = conversation
	}
	conversation.AddMessage(question)
	ConversationHubInstance.queue <- conversation
}

func Reset(from *model.From) bool {
	delete(ConversationHubInstance.Listener, *from)
	return true
}
