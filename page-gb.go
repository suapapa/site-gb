package main

import (
	"log"
	"net/http"
)

func gbHandler(w http.ResponseWriter, r *http.Request) {
	c := &PageContent{
		Title:     "⚔️ Guest Book ⚔️",
		Img:       "iamfine",
		Msg:       "Leave a message.",
		LastWords: "<a href=\"/support\">대가없는 🥩 환영합니다</a>",
		Root:      rootPath,
	}

	err := tmplPage.Execute(w, c)
	if err != nil {
		log.Printf("ERR: %v", err)
	}
}
