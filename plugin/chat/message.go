package chat

type RequestBody struct {
	Id      string `json:"id"`
	Preset  string `json:"preset"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Answer string `json:"answer"`
	Id     string `json:"id"`
	Cost   string `json:"cost"`
}
