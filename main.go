package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TODO
// const SENSOR_PATH = "/sys/bus/w1/devices/28-*"
const SENSOR_PATH = "/tmp/bus/w1/devices/28-*"
const BROKER_URL = "tcp://141.144.238.227:8000"
const TOPIC = "sensors/temperature"

var client mqtt.Client

func getSensorPath() (string, error) {
	for {
		matches, err := filepath.Glob(SENSOR_PATH)
		if err != nil {
			return "", err
		}

		matchCount := len(matches)
		if matchCount > 1 {
			return "", errors.New("Too many sensor folders found")
		} else if matchCount == 1 {
			return matches[0], nil
		}

		log.Info().Msgf("Sensor folder not found. Sleeping 1s...", SENSOR_PATH)
		time.Sleep(1 * time.Second)
	}
}

func waitUntilSensorReady(path string) error {
	w1Slave := path + "/w1_slave"

	for {
		file, err := os.Open(w1Slave)
		if err != nil {
			return err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		isReady := bytes.Contains(content, []byte("YES"))
		if isReady {
			return nil
		} else {
			log.Info().Msg("Sensor not ready. Sleep 1s...")
			time.Sleep(1 * time.Second)
		}
	}
}

func readSensor(path string) (int, error) {
	temperaturePath := path + "/temperature"

	file, err := os.Open(temperaturePath)
	if err != nil {
		return 0, err
	}

	var content [16]byte
	_, err = file.Read(content[:])
	if err != nil {
		return 0, err
	}

	temp, err := strconv.Atoi(string(bytes.Trim(content[:], " \n\x00")))
	if err != nil {
		return 0, err
	}

	return temp / 1000, nil
}

func publishTemperature(temperature int) error {
    token := client.Publish(TOPIC, 0, false, strconv.Itoa(temperature))
    token.Wait()
    return token.Error()
}

func initMqttClient() error {
	opts := mqtt.NewClientOptions().AddBroker(BROKER_URL)
	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := initMqttClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't connect to broker")
	}
        defer client.Disconnect(1000)

	path, err := getSensorPath()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	log.Info().Msgf("Sensor folder found '%v'", path)

	err = waitUntilSensorReady(path)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	log.Info().Msg("Sensor ready!")

	for {
		time.Sleep(1 * time.Second)

		temp, err := readSensor(path)
		if err != nil {
			log.Warn().Err(err).Msg("Failed reading sensor")
			continue
		}

		go func() {
                    err := publishTemperature(temp)
                    if err != nil {
                        log.Warn().Err(err).Msg("Failed publishing temperature to broker")
                    } else {
                        log.Debug().Msgf("Published temperature: %v", temp)
                    }
                }()
	}
}
