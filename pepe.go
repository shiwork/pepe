package main

import (
	"os"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os/signal"
	"github.com/shiwork/pepe/config"
	"github.com/shiwork/slack/incoming"
)

func main() {
	conf, err := config.Parse("./config.json")
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

	var slack_incoming_conf incoming.IncomingConf
	slack_incoming_conf.WebHookUrl = conf.Slack.IncomingWebHook

	go func() {
		log.Println("start watch")
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				// post Slack
				err := incoming.Post(
					slack_incoming_conf,
					incoming.Payload{
						[]*incoming.Attachment{
							&incoming.Attachment{
								event.Name,
								event.Name,
								"",
								[]*incoming.Field{
									&incoming.Field{
										"",
										"テストだよ",
										false,
									},
								},
							},
						},
					},
				)
				if err != nil {
					log.Println("error:", err)
				}
			case err := <-watcher.Errors:
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

	s := <-c
	log.Println("signal:", s)
}
