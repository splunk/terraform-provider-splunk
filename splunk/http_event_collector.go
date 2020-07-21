package splunk

// HTTP Input Response Schema
type Response struct {
	Entry []Entry `json:"entry"`
	Messages []ErrorMessage `json:"messages"`
}

type Entry struct {
	Name    string  `json:"name"`
	Content Content `json:"content"`
}

type Content struct {
	Token string `json:"token"`
	Index string `json:"index"`
}

type ErrorMessage struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
}
