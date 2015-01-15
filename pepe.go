package main

import (
	"github.com/go-fsnotify/fsnotify"
	"github.com/shiwork/pepe/config"
	"github.com/shiwork/slack/incoming"
	"log"
	"os"
	"path"
	"strings"
)

var confPath = os.Getenv("PEPE_CONFIG")

func main() {
	conf, err := config.Parse(confPath)
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

	var slackIncomingConf incoming.IncomingConf
	slackIncomingConf.WebHookUrl = conf.Slack.IncomingWebHook

	go func() {
		log.Println("start watch")
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					var filePath string
					var dirPath string

					for _, value := range conf.Watches {
						if strings.HasPrefix(event.Name, value.Dir) {
							filePath = strings.Replace(event.Name, value.Dir, "", -1)
							dirPath = value.Url
							break
						}
					}

					// create post message
					var _, filename = path.Split(event.Name)

					var message string
					message = "<" + dirPath + filePath + "|" + filename + ">"

					// post message to Slack
					err := incoming.Post(
						slackIncomingConf,
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
