package app

import (
	"github.com/noolingo/swagger/internal/domain"
	"github.com/noolingo/swagger/internal/service"
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
