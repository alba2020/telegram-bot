package domain

import "context"

type Message struct {
	Text    string
	UserID  int64
	Buttons map[string]string // label => data
	Context context.Context
	Command string
}

func (m *Message) Reply(text string) *Message {
	return &Message{
		Text:    text,
		UserID:  m.UserID,
		Buttons: nil,
	}
}
