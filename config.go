package main

import (
	"code.google.com/p/gcfg"
	"os"
)

type cfg struct {
	Notifier struct {
		Server []string
		To     string
	}
	Mail struct {
		Host     string
		Port     int
		User     string
		Password string
		From     string
	}
}

var Config = cfg{}

func loadConfig(configfile string) {
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		logger.Criticalf("Configuration file: %s not found", configfile)
		os.Exit(1)
	}

	gcfg.ReadFileInto(&Config, configfile)
}
