package deep_seek

type ResponseBody struct {
	Id                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}
type Choice struct {
	Index          int             `json:"index"`
	Message        MessageResponse `json:"message"`
	FinishedReason string          `json:"finish_reason"`
}

type MessageResponse struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}
type Usage struct {
	PromptTokens             int `json:"prompt_tokens"`
	CompletionTokens         int `json:"completion_tokens"`
	TotalTokens              int `json:"total_tokens"`
	*PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
	*CompletionTokensDetails `json:"completion_tokens_details"`
	PromptCachedHitTokens    int `json:"prompt_cached_hit_tokens"`
	PromptCachedMissTokens   int `json:"prompt_cached_miss_tokens"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}
