package main

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func removeOldFilesThread(path string, olderThan time.Duration) {
	ticker := time.NewTicker(time.Hour)

	for _ = range ticker.C {
		log.Debugln("Removing old files...")
		err := removeOldFiles(path, olderThan)
		if err != nil {
			log.Println(err)
		}
		log.Debugln("Removing old files completed.")
	}
}

func removeOldFiles(path string, olderThan time.Duration) error {
	return filepath.Walk(config.StoragePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if info.ModTime().Add(olderThan).Before(time.Now()) {
			log.WithField("filepath", path).Debugln("Removing file...")
			err := os.Remove(path)
			if err != nil {
				log.Println(err)
			}
		}

		return nil
	})

}
