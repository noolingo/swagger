package main

import (
	"os"

	"github.com/MelnikovNA/noolingoswagger/internal/app"
	"github.com/sirupsen/logrus"
)

const configPath = "./config/config.yml"

func main() {
	err := app.Run(configPath)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
