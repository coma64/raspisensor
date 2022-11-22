package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("gotest").SetKeepAlive(2 * time.Second).SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Couldn't connect to broker")
	}
	defer client.Disconnect(250)

	for i := 0; i < 100; i++ {
		token := client.Publish("mqtt/sample", 0, false, fmt.Sprintf("Hello #%v", i))
		token.Wait()

		time.Sleep(time.Duration(rand.Intn(3000)) + time.Second)
	}
}
