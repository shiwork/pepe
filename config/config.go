package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Dir   []string
	Slack SlackConfig
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
