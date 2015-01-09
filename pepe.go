package main

import (
	"github.com/go-fsnotify/fsnotify"
	"github.com/shiwork/pepe/config"
	"github.com/shiwork/slack/incoming"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
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

	go func() {
		log.Println("start watch")
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					var file_path string
					var dir_path string

					for _, value := range conf.Watches {
						if strings.HasPrefix(event.Name, value.Dir) {
							file_path = strings.Replace(event.Name, value.Dir, "", -1)
							dir_path = value.Url
							break
						}
					}

					// create post message
					var _, file_name = path.Split(event.Name)

					var message string
					message = "<" + dir_path + file_path + "|" + file_name + ">"

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
