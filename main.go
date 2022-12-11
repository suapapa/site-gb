package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/suapapa/site-gb/msg"
)

var (
	urlPrefix  string
	httpPort   int
	enableMQTT bool

	mqttC mqtt.Client
)

func main() {
	log.Println("homin.dev guestbook start")
	defer log.Println("homin.dev guestbook stop")

	flag.StringVar(&urlPrefix, "p", "/", "set url prefix")
	flag.IntVar(&httpPort, "http", 8080, "set http port")
	flag.Parse()

	if !strings.HasPrefix(urlPrefix, "/") {
		urlPrefix = "/" + urlPrefix
	}

	http.HandleFunc(urlPrefix, gbHandler)
	go func() {
		// start HTTPServer
		log.Printf("listening http on :%d", httpPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			log.Fatalf("ERR: %v", err)
		}
	}()

	if mqttURL, err := url.Parse(os.Getenv("MQTT_HOST")); err != nil {
		log.Printf("WARN: mqtt disabled by %v", err)
	} else {
		mqttHost := mqttURL.Hostname()
		mqttPort := mqttURL.Port()
		mqttC, err = connectBrokerByWS(&Config{
			Host:     mqttHost,
			Port:     mqttPort,
			Username: os.Getenv("MQTT_USERNAME"),
			Password: os.Getenv("MQTT_PASSWORD"),
			// CaCert:   "/etc/ssl/certs/ca-certificates.crt",
		})
		if err != nil {
			log.Printf("WARN: mqtt disabled by %v", err)
		} else {
			log.Printf("mqtt enabled")
			enableMQTT = true
		}
	}

	if enableMQTT {
		defer mqttC.Disconnect(1000)
		go func() {
			porkV := msg.NewPorkMsg()
			tk := time.NewTicker(30 * time.Minute)
			defer tk.Stop()
			for range tk.C {
				mqttPub(topic, porkV)
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
