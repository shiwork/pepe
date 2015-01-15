package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Watch struct
type Watch struct {
	Dir string `json:"dir"`
	URL string `json:"url"`
}

// Config struct
type Config struct {
	Watches []Watch     `json:"watch"`
	Slack   SlackConfig `json:"slack"`
}

// SlackConfig struct
type SlackConfig struct {
	IncomingWebHook string `json:"incoming"`
}

// Parse config
func Parse(filename string) (Config, error) {
	var config Config
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Failed to read config file:", err)
		return config, err
	}

	err = json.Unmarshal(jsonString, &config)
	if err != nil {
		log.Println("Failed to json unmarshal:", err)
		return config, nil
	}
	return config, nil
}
