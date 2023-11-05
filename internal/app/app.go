package app

import (
	"github.com/MelnikovNA/noolingoswagger/internal/domain"
	"github.com/MelnikovNA/noolingoswagger/internal/service"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func Run(configPath string) error {
	logger := logrus.New()
	config := new(domain.Config)
	err := cleanenv.ReadConfig(configPath, config)

	if err == nil {
		parseFlags(config)
		err = service.Swagger(config, logger)
	}
	return err
}
