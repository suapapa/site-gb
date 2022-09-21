package main

import (
	"log"
	"net/http"
	"time"

	"github.com/suapapa/site-gb/msg"
)

var (
	lastGB *msg.GuestBook
)

func gbHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("hit %s", r.URL.Path)
	c := &PageContent{
		Title:     "💌 방명록 💌️",
		Img:       "https://homin.dev/asset/image/gb.jpg",
		Msg:       "익명 가능하며, 서버에 저장되지 않습니다",
		LastWords: "<a href=\"https://homin.dev/blog/post/20220910_live_print_guestbook_with_mqtt/\">💌는 어디로 가나?</a>",
	}

	switch r.Method {

	// when press SEND
	case "POST":
		c.Success = true
		c.Msg = "❤️ 보냄 ❤️"
		c.LastWords = "<a href=\"/ingress\">대문으로 이동</a>"
		err := tmplPage.Execute(w, c)
		if err != nil {
			log.Printf("ERR: %v", err)
			return
		}

		r.ParseForm()
		m := msg.NewGuestBookMsg(r.PostFormValue("name"), r.PostFormValue("msg"))
		if err = sendMsgToTelegram(m); err != nil {
			log.Printf("ERR: %v", err)
			return
		}

		gb := m.Data.(msg.GuestBook)
		if !gb.IsSame(lastGB) {
			if lastGB != nil {
				lastTS, _ := time.Parse(lastGB.TimeStamp, time.RFC3339)
				currTS, _ := time.Parse(gb.TimeStamp, time.RFC3339)
				if currTS.Sub(lastTS) > 3*time.Second {
					log.Printf("WARN: same msgs in 3 seconds")
					return
				}
			}
			lastGB = &gb
			if err := mqttPub(topic, m); err != nil {
				log.Printf("ERR: %v", err)
				return
			}
		}

	// form Page
	case "GET":
		if err := tmplPage.Execute(w, c); err != nil {
			log.Printf("ERR: %v", err)
		}
	}
}
