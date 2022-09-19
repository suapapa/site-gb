package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/suapapa/site-gb/msg"
)

var (
	rootPath string
	httpPort int

	mqttC mqtt.Client
)

func main() {
	log.Println("homin.dev guestbook start")
	defer log.Println("homin.dev guestbook stop")

	flag.StringVar(&rootPath, "root", "/", "set subdomain root")
	flag.IntVar(&httpPort, "http", 8080, "set http port")
	flag.Parse()

	if !strings.HasSuffix(rootPath, "/") {
		rootPath += "/"
	}

	http.HandleFunc(rootPath+"", gbHandler)
	go func() {
		// start HTTPServer
		log.Printf("listening http on :%d", httpPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			log.Fatalf("ERR: %v", err)
		}
	}()

	var err error
	mqttC, err = connectBrokerByWSS(&Config{
		Host:     "homin.dev",
		Port:     9001,
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	defer mqttC.Disconnect(1000)

	go func() {
		porkV := msg.NewPorkMsg()
		tk := time.NewTicker(30 * time.Second)
		defer tk.Stop()
		for range tk.C {
			mqttPub(topic, porkV)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
