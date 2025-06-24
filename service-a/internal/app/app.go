package app

import (
	"github.com/leonardo-gmuller/cep-weather-system/service-a/internal/app/config"
	"github.com/leonardo-gmuller/cep-weather-system/service-a/internal/app/domain/usecase"
)

type App struct {
	UseCase usecase.UseCaseInterface
}

func New(config config.Config) (*App, error) {
	usecase := usecase.New(&config)
	return &App{
		UseCase: usecase,
	}, nil
}
