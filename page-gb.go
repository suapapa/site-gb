package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	loc = time.FixedZone("UTC+9", 9*60*60)
)

// func init() {
// 	var err error
// 	loc, err = time.LoadLocation("Asia/Seoul")
// 	if err != nil {
// 		panic(errors.Wrap(err, "can't get loc for Asia/Seoul"))
// 	}
// }

func gbHandler(w http.ResponseWriter, r *http.Request) {
	c := &PageContent{
		Title:     "💌 방명록 💌️",
		Img:       "https://homin.dev/asset/image/gb.jpg",
		Msg:       "익명이 가능하며, 저장되지 않습니다",
		LastWords: "<a href=\"/support\">대가없는 🥩 환영합니다</a>",
		Root:      rootPath,
	}

	switch r.Method {

	// when press SEND
	case "POST":
		c.Success = true
		err := tmplPage.Execute(w, c)
		if err != nil {
			log.Printf("ERR: %v", err)
			return
		}

		r.ParseForm()
		msg := map[string]string{
			"msg":        strings.TrimSpace(r.PostFormValue("msg")),
			"from":       strings.TrimSpace(r.PostFormValue("name")),
			"remoteAddr": r.RemoteAddr,
			"timestamp":  time.Now().In(loc).Format(time.RFC3339),
			// "timestamp": time.Now().Format(time.RFC3339),
		}
		if err = sendMsgToTelegram(makeMsgStringForTelegram(msg)); err != nil {
			log.Printf("ERR: %v", err)
			return
		}

		if err := mqttPub(topic, msg); err != nil {
			log.Printf("ERR: %v", err)
		}

	// form Page
	case "GET":
		if err := tmplPage.Execute(w, c); err != nil {
			log.Printf("ERR: %v", err)
		}
	}
}
