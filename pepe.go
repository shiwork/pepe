package main

import (
	"os"
	"encoding/json"
	"net/http"
	"net/url"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os/signal"
	"github.com/shiwork/pepe/config"
)

var webhookUrl string = os.Getenv("PEPE_SLACK_INCOMING_URL")

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool `json:"short"`
}

type attachment struct {
	Fallback string `json:"fallback"`
	Pretext  string `json:"pretext"`
	Color    string `json:"color"`
	Fields   []*field `json:"fields"`
}

type payload struct{
	Attachments []*attachment `json:"attachments"`
}

func Post(title string, description string) error {
	p, err := json.Marshal(&payload{
		Attachments: []*attachment{
			&attachment{
				Fallback: title,
				Pretext: title,
				Fields: []*field{
					&field{
						Title: "",
						Value: description,
						Short: false,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = http.PostForm(webhookUrl, url.Values{
		"payload": []string{string(p)},
	})

	if err != nil {
		return err
	}
	return nil
}

func main() {
	conf, err := config.Parse("/Users/yuichi.shiwaku/go/src/github.com/shiwork/pepe/config.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		log.Println("start watch")
		for {
			select {
			case event:= <-watcher.Events:
				log.Println("event:", event)
			case err:= <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// 監視ディレクトリの追加
	for _, value := range conf.Dir {
		err = watcher.Add(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	s:=<-c
	log.Println("signal:", s)
}
