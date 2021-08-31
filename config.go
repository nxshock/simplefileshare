package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

var config *Config

type Config struct {
	ListenAddress    string
	StoragePath      string
	RemoveFilePeriod uint // hours
	LogLevel         log.Level
}

func initConfig() error {
	log.Debugln("Сonfig initialization started.")
	defer log.Debugln("Сonfig initialization finished.")

	var configFilePath string

	if len(os.Args) < 2 {
		configFilePath = defaultConfigFilePath
	} else {
		configFilePath = os.Args[1]
	}

	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return err
	}

	stat, err := os.Stat(config.StoragePath)
	if err != nil {
		return fmt.Errorf("os.Stat(config.StoragePath): %v", err)
	}
	if !stat.IsDir() {
		return errors.New("StoragePath is not a dir")
	}

	return nil
}
