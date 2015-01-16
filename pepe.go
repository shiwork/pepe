package main

import (
	"github.com/go-fsnotify/fsnotify"
	"github.com/shiwork/fsmonitor"
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
	}

	watcher, err := fsmonitor.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	var slackConfig incoming.Config
	slackConfig.WebHookURL = conf.Slack.IncomingWebHook

	// 監視ディレクトリの追加
	for _, value := range conf.Watches {
		log.Println(value.Dir)
		err = watcher.Watch(value.Dir)
		if err != nil {
			log.Fatal(err)
		}
	}

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
						dirPath = value.URL
						break
					}
				}

				// create post message
				var _, filename = path.Split(event.Name)

				var message string
				message = "<" + dirPath + filePath + "|" + filename + ">"

				// post message to Slack
				err := incoming.Post(
					slackConfig,
					incoming.Payload{
						[]incoming.Attachment{
							incoming.Attachment{
								message,
								message,
								"",
								[]incoming.Field{
									incoming.Field{
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
		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
}
