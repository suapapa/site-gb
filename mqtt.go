package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type Config struct {
	Host               string
	Port               int
	Username, Password string
	CaCert             string
}

func connectBrokerByWSS(config *Config) (mqtt.Client, error) {
	var tlsConfig tls.Config
	if config.CaCert == "" {
		config.CaCert = "/etc/ssl/certs/ca-certificates.crt"
	}

	certpool := x509.NewCertPool()
	ca, err := os.ReadFile(config.CaCert)
	if err != nil {
		return nil, errors.Wrap(err, "fail to connet broker")
	}
	certpool.AppendCertsFromPEM(ca)
	tlsConfig.RootCAs = certpool

	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("wss://%s:%d/mqtt", config.Host, config.Port)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetTLSConfig(&tlsConfig)
	opts.SetClientID("site-gb")
	// opts.SetKeepAlive(20)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, errors.Wrap(err, "fail to connet broker")
	}
	return client, nil
}
