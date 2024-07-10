package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{ForceColors: true, DisableTimestamp: true})
	log.SetLevel(log.ErrorLevel)

	err := initConfig()
	if err != nil {
		log.Fatalln("initConfig:", err)
	}

	log.SetLevel(config.LogLevel)

	err = initTemplates()
	if err != nil {
		log.Fatalln("initTemplates:", err)
	}

	if config.RemoveFilePeriod > 0 {
		go removeOldFilesThread(time.Duration(config.RemoveFilePeriod) * time.Hour)
	}

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/icon", HandleIcon)
	http.HandleFunc("/upload", HandleUpload)
	http.HandleFunc("/download", HandleDownload)
	http.HandleFunc("/stream", HandleStream)
}

func main() {
	go func() {
		err := http.ListenAndServe(config.ListenAddress, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Debugln("Stop signal received.")
}
