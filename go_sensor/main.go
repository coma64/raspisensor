package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/coma64/raspisensor/config"
	"github.com/coma64/raspisensor/mqtt_client"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func getSensorPath() (string, error) {
	for {
		matches, err := filepath.Glob(config.Conf.SensorGlob)
		if err != nil {
			return "", fmt.Errorf("unable to expand glob: %w", err)
		}

		matchCount := len(matches)
		if matchCount > 1 {
			return "", errors.New("Too many sensor folders found")
		} else if matchCount == 1 {
			return matches[0], nil
		}

		log.Debug().Str("glob", config.Conf.SensorGlob).Msg("Sensor folder not found. Sleeping...")
		time.Sleep(1 * time.Second)
	}
}

func waitUntilSensorReady(path string) error {
	w1Slave := path + "/w1_slave"

	for {
		file, err := os.Open(w1Slave)
		if err != nil {
			return fmt.Errorf("unable to open slave file: %w", err)
		}

		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("unable to read slave file: %w", err)
		}

		isReady := bytes.Contains(content, []byte("YES"))
		if isReady {
			return nil
		} else {
			log.Debug().Str("path", w1Slave).Msg("Sensor not ready. Sleeping...")
			time.Sleep(1 * time.Second)
		}

		_ = file.Close()
	}
}

func readSensor(path string) (int, error) {
	temperaturePath := path + "/temperature"

	file, err := os.Open(temperaturePath)
	if err != nil {
		return 0, fmt.Errorf("unable to open temperature file: %w", err)
	}

	var content [16]byte
	_, err = file.Read(content[:])
	if err != nil {
		return 0, fmt.Errorf("unable to read temperature file: %w", err)
	}

	temp, err := strconv.Atoi(string(bytes.Trim(content[:], " \n\x00")))
	if err != nil {
		return 0, fmt.Errorf("unable to parse temperature: %w", err)
	}

	return temp, nil
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if config.Conf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	defer mqtt_client.Disconnect()

	path, err := getSensorPath()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting sensor path")
	}

	log.Info().Str("path", path).Msgf("Sensor folder found")

	err = waitUntilSensorReady(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Error waiting for sensor")
	}

	log.Info().Msg("Sensor ready!")

	log.Info().
		Str("topic", config.Conf.Broker.Topic).
		Str("broker", config.Conf.Broker.URL).
		Msg("Starting to publish")

	for {
		time.Sleep(1 * time.Second)

		temp, err := readSensor(path)
		if err != nil {
			log.Warn().Err(err).Msg("Error reading sensor")
			continue
		}

		go func() {
			err := mqtt_client.PublishTemperature(temp)
			if err != nil {
				log.Warn().Err(err).Msg("Error publishing temperature to broker")
			}
		}()
	}
}
