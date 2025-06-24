package usecase

import (
	"context"

	"github.com/leonardo-gmuller/cep-weather-system/service-a/internal/app/config"
)

type UseCaseInterface interface {
	ValidateCEP(ctx context.Context, cep string) (bool, error)
}

type UseCase struct {
	AppName string
}

func New(config *config.Config) *UseCase {
	return &UseCase{
		AppName: config.App.Name,
	}
}
