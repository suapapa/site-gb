package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	http.HandleFunc(rootPath+"img/", imgHandler)
	http.HandleFunc(rootPath+"send/", sendHandler)

	// start HTTPServer
	go func() {
		log.Printf("listening http on :%d", httpPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
			log.Fatal(err)
		}
	}()

	exitCh := make(chan any)
	<-exitCh
}
