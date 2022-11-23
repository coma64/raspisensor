package config

import "github.com/jinzhu/configor"

type config struct {
	Broker struct {
		URL   string
		Topic string
	}
	SensorGlob string
}

var Conf = &config{}

func init() {
	if err := configor.Load(Conf, "config.yml", "/app/config.yml"); err != nil {
		panic(err)
	}
}
