package main

import (
	"errors"
	"log"
	"os"
	"github.com/kylelemons/go-gypsy/yaml"
)

// start the logger early
func init() {
	// Log to stdout
	log.SetOutput(os.Stdout)
}

// Read config and start the service
func main() {
	// Read config from file
	config, err := readConfig("conf.yaml")
	if err != nil {
		// Panic should there's a problem with config
		panic(errors.New("unable to parse conf.yaml file"))
	}

	// Pass along the config and start the server
	Start(config)
}

// Reads conf.yaml into Config struct
func readConfig(configFile string) (Config, error) {
	log.Println("Reading configuration...")

	config := new(Config)
	c, err := yaml.ReadFile(configFile)
	if err != nil {
		// returning an empty config and an error
		return *config, err
	}

	// By choice: Translate conf.yaml config to Config struct
	// TODO check for error
	config.BackingRedisUrl, _ = c.Get("redis.url")
	config.Capacity, _ = c.GetInt("cache.capacity")
	config.Expiry, _ = c.GetInt("cache.expiry")
	config.Port, _ = c.Get("service.port")

	return *config, nil
}
