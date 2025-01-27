package chat

import (
	"my-qqbot/model"
)

func Ask(from model.From, question string) {
	if ConversationHubInstance == nil {
		ConversationHubInstance = NewConversationHub()
	}
	conversation, ok := ConversationHubInstance.Listener[from.Id]
	if !ok {
		conversation = NewConversation(from)
		ConversationHubInstance.Listener[from.Id] = conversation
	}
	conversation.AddMessage(question)
	ConversationHubInstance.queue <- conversation
}

func Reset(from *model.From) bool {
	delete(ConversationHubInstance.Listener, from.Id)
	return true
}
