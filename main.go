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
)

var (
	rootPath string
	httpPort int
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
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
