package main

import "github.com/coma64/raspisensor/mqtt_client"

func main() {
	defer mqtt_client.Disconnect()

	for i := 0; i < 200; i++ {
		err := mqtt_client.PublishTemperature(i)
		if err != nil {
			panic(err)
		}
	}
}
