package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

const (
	topic = "homin-dev/gb"
)

type Config struct {
	Scheme             string
	Host               string
	Port               string
	Username, Password string
	CaCert             string
}

func connectBrokerByWS(config *Config) (mqtt.Client, error) {
	var tlsConfig tls.Config

	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("%s://%s:%s", config.Scheme, config.Host, config.Port)
	log.Infof("connecting to %s", broker)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	if config.CaCert != "" {
		certpool := x509.NewCertPool()
		ca, err := os.ReadFile(config.CaCert)
		if err != nil {
			return nil, errors.Wrap(err, "fail to connet broker")
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
		opts.SetTLSConfig(&tlsConfig)
	}
	opts.SetClientID("site-gb")
	// opts.SetKeepAlive(20)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, errors.Wrap(err, "fail to connect broker")
	}
	return client, nil
}

func mqttPub(topic string, jsonV any) error {
	if !enableMQTT {
		return errors.New("WARN: mqtt not enabled")
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(jsonV)
	if err != nil {
		return errors.Wrap(err, "fail to pub mqtt")
	}

	// Send it to mqtt
	tkn := mqttC.Publish(topic, 0, false, buf.Bytes())
	tkn.Wait()
	if err := tkn.Error(); err != nil {
		return errors.Wrap(err, "fail to pub mqtt")
	}
	return nil
}
