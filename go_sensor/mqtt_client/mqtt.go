package mqtt_client

import (
	"fmt"
	"github.com/coma64/raspisensor/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"strconv"
)

var (
	client mqtt.Client
)

func PublishTemperature(temperature int) error {
	log.Debug().
		Int("temperature", temperature).
		Msg("Publishing")

	token := client.Publish(config.Conf.Broker.Topic, 0, false, strconv.Itoa(temperature))
	token.Wait()
	return token.Error()
}

func Disconnect() {
	client.Disconnect(1000)
}

func init() {
	opts := mqtt.NewClientOptions().AddBroker(config.Conf.Broker.URL)
	client = mqtt.NewClient(opts)

	log.Info().Str("broker", config.Conf.Broker.URL).Msg("Connecting to broker")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Errorf("unable to connect to broker: %w", token.Error()))
	}
}
