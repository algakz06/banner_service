package main

import (
	"github.com/algakz/banner_service/config"
	"github.com/algakz/banner_service/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		logrus.Fatalf("%s", err.Error())
	}

	server := server.NewApp()
	if err := server.Run(viper.GetString("port")); err != nil {
		logrus.Fatalf("error occured while running server on port 3333: %s", err.Error())
	}
}
