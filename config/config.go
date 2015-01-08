package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Watch struct {
	Dir string `json:"dir`
	Url string `json:"url`
}

type Config struct {
	Watches []Watch `json:"watch"`
	Slack SlackConfig `json:"slack"`
}

type SlackConfig struct {
	IncomingWebHook string `json:"incoming"`
}

func Parse(filename string) (Config, error) {
	var config Config
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Failed to read config file: %v", err)
		return config, err
	}

	err = json.Unmarshal(jsonString, &config)
	if err != nil {
		log.Println("Failed to json unmarshal: %v", err)
		return config, nil
	}
	return config, nil
}
