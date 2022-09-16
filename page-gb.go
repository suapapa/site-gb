package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	topic = "homin-dev/gb"
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
	log.Printf("hit %s", r.URL.Path)
	c := &PageContent{
		Title:     "💌 방명록 💌️",
		Img:       "https://homin.dev/asset/image/gb.jpg",
		Msg:       "익명 가능하며, 서버에 저장되지 않습니다",
		LastWords: "<a href=\"https://homin.dev/blog/post/20220910_live_print_guestbook_with_mqtt/\">💌는 어디로 가나?</a>",
		Root:      rootPath,
	}

	if r.Method == "POST" {
		c.Success = true
		c.Msg = "보냄❤️ <a href=\"/ingress\">대문으로 이동</a>"
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
		err = sendMsgToTelegram(makeMsgStringForTelegram(msg))
		if err != nil {
			log.Printf("ERR: %v", err)
			return
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
		defer mqttC.Disconnect(1000)
		mqttC.Publish(topic, 0, false, buf.Bytes())

		return
	}

	err := tmplPage.Execute(w, c)
	if err != nil {
		log.Printf("ERR: %v", err)
	}
}

func makeMsgStringForTelegram(in map[string]string) string {
	outFmt := `## 방명록 ##
%s
- %s -`
	out := fmt.Sprintf(outFmt,
		strings.ReplaceAll(in["msg"], "\r\n", "\n"),
		in["from"],
	)

	log.Println(out)
	return out
}
