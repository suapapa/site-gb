package msg

import (
	"fmt"
	"strings"
	"time"
)

const (
	MTGuestBook = "gb"
	MTPork      = "pork"
)

var (
	kst = time.FixedZone("KST", 9*60*60)
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

func (m *Message) GetGuestBook() (*GuestBook, error) {
	if m.Type != MTGuestBook {
		return nil, fmt.Errorf("not GuestBook Type")
	}

	gb, ok := m.Data.(*GuestBook)
	if !ok {
		return nil, fmt.Errorf("data is not a *GuestBook type")
	}

	return gb, nil
}

type GuestBook struct {
	From      string `json:"from"`
	Content   string `json:"content"`
	TimeStamp string `json:"ts"`
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
			TimeStamp: time.Now().In(kst).Format(time.RFC3339),
		},
	}
}

func NewPorkMsg() *Message {
	return &Message{
		Type: MTPork,
	}
}
