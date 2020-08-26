package models

type ConfStanzaResponse struct {
	Entry    []ConfStanzaEntry `json:"entry"`
	Messages []ErrorMessage    `json:"messages"`
}

type ConfStanzaEntry struct {
	Name     string           `json:"name"`
	ACL      ACLObject        `json:"acl"`
	Content  ConfStanzaObject `json:"content"`
}

type ConfStanzaObject struct {
	Variables map[string]string `json:"-"`
}