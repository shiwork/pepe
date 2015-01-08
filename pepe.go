package main

import (
	"os"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os/signal"
	"github.com/shiwork/pepe/config"
	"github.com/shiwork/slack/incoming"
	"path"
)

var conf_path string = os.Getenv("PEPE_CONFIG")

func main() {
	conf, err := config.Parse(conf_path)
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

	var watched = conf.Watches[0]

	go func() {
		log.Println("start watch")
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				switch {
				case event.Op&fsnotify.Create == fsnotify.Write:
					// create post message
					var _, file_name = path.Split(event.Name)

					var message string
					message = "<"+watched.Url+"/"+file_name+"|"+file_name+">"

					// post message to Slack
					err := incoming.Post(
						slack_incoming_conf,
						incoming.Payload{
							[]*incoming.Attachment{
								&incoming.Attachment{
									message,
									message,
									"",
									[]*incoming.Field{
										&incoming.Field{
											"",
											"",
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
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// 監視ディレクトリの追加
	for _, value := range conf.Watches {
		log.Println(value.Dir)
		err = watcher.Add(value.Dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	s := <-c
	log.Println("signal:", s)
}
