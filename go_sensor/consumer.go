package main

import (
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var handler mqtt.MessageHandler = func(c mqtt.Client, msg mqtt.Message) {
	log.Info().Msgf("Got message: %v", string(msg.Payload()))
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Couldn't connect")
	}

	client.Subscribe("mqtt/sample", 0, handler)

	time.Sleep(10 * time.Second)
}
