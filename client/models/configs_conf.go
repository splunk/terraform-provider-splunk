package models

type ConfigsConfResponse struct {
	Entry    []ConfigsConfEntry `json:"entry"`
	Messages []ErrorMessage     `json:"messages"`
}

type ConfigsConfEntry struct {
	Name    string            `json:"name"`
	ACL     ACLObject         `json:"acl"`
	Content ConfigsConfObject `json:"content"`
}

type ConfigsConfObject struct {
	Variables map[string]string `json:"-"`
}
