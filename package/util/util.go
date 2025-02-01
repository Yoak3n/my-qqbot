package util

import (
	"encoding/json"
	"os"
)

func CreateDirNotExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		e := os.MkdirAll(dir, os.ModePerm)
		if e != nil {
			println("Error creating directory: " + e.Error())
			return
		}
	}
}

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
