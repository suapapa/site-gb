package msg

import "testing"

func TestGetGuestBook(t *testing.T) {
	gbMsg := NewGuestBookMsg("Hello", "World")
	gb, err := gbMsg.GetGuestBook()
	if err != nil {
		t.Error(err)
	}
	if gb.Content != "World" || gb.From != "Hello" {
		t.Error("content mismatch")
	}
}
