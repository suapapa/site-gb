package main

import (
	"net/url"
	"testing"
)

func TestUrlParse(t *testing.T) {
	mqttURL, _ := url.Parse("wss://localhost:1883")
	host, port := mqttURL.Hostname(), mqttURL.Port()
	t.Log(host, port)

}
