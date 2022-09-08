package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	topic = "homin-dev/gb"
)

func gbHandler(w http.ResponseWriter, r *http.Request) {
	c := &PageContent{
		Title:     "⚔️ Guest Book ⚔️",
		Img:       "iamfine",
		Msg:       "Leave a message.",
		LastWords: "<a href=\"/support\">대가없는 🥩 환영합니다</a>",
		Root:      rootPath,
	}

	if r.Method == "POST" {
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
			"timestamp":  time.Now().Format(time.RFC3339),
		}
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(msg)

		// Send it to mqtt
		mqttC, err := connectBrokerByWSS(&Config{
			Host:     "homin.dev",
			Port:     9001,
			Username: os.Getenv("MQTT_USERNAME"),
			Password: os.Getenv("MQTT_PASSWORD"),
		})
		if err != nil {
			log.Printf("ERR: %v", err)
			return
		}
		mqttC.Publish(topic, 0, false, buf.Bytes())

		return
	}

	err := tmplPage.Execute(w, c)
	if err != nil {
		log.Printf("ERR: %v", err)
	}
}
