package main

import (
	"log"
	"net/http"
)

func sendHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ERR: %v", err)
	}

	log.Println(r.Form)
}
