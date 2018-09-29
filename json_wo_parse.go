package main

import (
	"encoding/json"
	"fmt"
)

// 変数名の先頭を大文字にしないとダメみたい
type IssueData struct {
	Id          json.Number `json:"id"`
	Description string      `json:"description"`
}

type Message struct {
	Issue IssueData `json:"issue"`
}

func main() {
	b := []byte(`{"issue":{"id":4848, "description": "Alice"}}`)
	var m Message
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		return
	}
	println(m.Issue.Id)
	println(m.Issue.Description)
}
