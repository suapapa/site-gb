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
		Title:     "π λ°©λͺλ‘ ποΈ",
		Img:       "https://homin.dev/asset/image/gb.jpg",
		Msg:       "μ΅λͺ κ°λ₯νλ©°, μλ²μ μ μ₯λμ§ μμ΅λλ€",
		LastWords: "<a href=\"https://homin.dev/blog/post/20220910_live_print_guestbook_with_mqtt/\">πλ μ΄λλ‘ κ°λ?</a>",
	}

	switch r.Method {

	// when press SEND
	case "POST":
		c.Success = true
		c.Msg = "β€οΈ λ³΄λ β€οΈ"
		c.LastWords = "<a href=\"/ingress\">λλ¬ΈμΌλ‘ μ΄λ</a>"
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

		gb, err := m.GetGuestBook()
		if err != nil {
			log.Printf("ERR: %v", err)
			return
		}

		if !gb.IsSame(lastGB) {
			if lastGB != nil {
				lastTS, _ := time.Parse(lastGB.TimeStamp, time.RFC3339)
				currTS, _ := time.Parse(gb.TimeStamp, time.RFC3339)
				if currTS.Sub(lastTS) > 3*time.Second {
					log.Printf("WARN: same msgs in 3 seconds")
					return
				}
			}
			lastGB = gb
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
