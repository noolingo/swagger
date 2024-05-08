package main

import (
	"os"

	"github.com/noolingo/swagger/internal/app"
	"github.com/sirupsen/logrus"
)

const configPath = "./configs/config.yml"

func main() {
	err := app.Run(configPath)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
