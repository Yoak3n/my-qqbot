package util

import (
	"encoding/json"
)

type ReasoningMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
}

func GetReasoningContent(msg []byte) string {
	reasoning := ReasoningMessage{}
	err := json.Unmarshal(msg, &reasoning)
	if err != nil {
		return err.Error()
	}
	return reasoning.ReasoningContent
}
