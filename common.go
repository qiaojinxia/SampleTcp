package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var (
	AppConfig = Config{}
)

func Init() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("C:\\Users\\a\\go\\SampleTcp\\config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Error(err)
	}

}

type Config struct {
	Version    string `yaml:"Version"`
	Port       string `yaml:"Port"`
	ServerName string `yaml:"ServerName"`
}

var cacheBytes = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 1024)
		return &b
	},
}

var (
	process sync.WaitGroup
)

func HandlerFunc(f func() error) {
	process.Add(1)
	err := f()
	if err != nil {
		log.Error(err)
	}
	process.Done()
}

func HandlerAsyncFunc(f func() error) {
	process.Add(1)
	go func() {
		err := f()
		if err != nil {
			log.Error(err)
		}
		process.Done()
	}()
}

func isChanClose(ch chan os.Signal) bool {
	select {
	case _, received := <-ch:
		return !received
	default:
	}
	return false
}
