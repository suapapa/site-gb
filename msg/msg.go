package msg

import (
	"strings"
	"time"
)

const (
	MTGuestBook = "gb"
	MTPork      = "pork"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type GuestBook struct {
	From      string    `json:"from"`
	Content   string    `json:"content"`
	TimeStamp time.Time `json:"ts"`
}

func (b *GuestBook) IsSame(b2 *GuestBook) bool {
	if b2 == nil {
		return false
	}
	return b.From == b2.From && b.Content == b2.Content
}

func NewGuestBookMsg(from, content string) *Message {
	return &Message{
		Type: MTGuestBook,
		Data: &GuestBook{
			From:      from,
			Content:   strings.ReplaceAll(content, "\r\n", "\n"),
			TimeStamp: time.Now(),
		},
	}
}

func NewPorkMsg() *Message {
	return &Message{
		Type: MTPork,
	}
}
