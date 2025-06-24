package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/config"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/domain/usecase"
)

type Handler struct {
	cfg     config.Config
	useCase usecase.UseCaseInterface
}

func New(cfg config.Config, useCase usecase.UseCaseInterface) Handler {
	return Handler{
		cfg:     cfg,
		useCase: useCase,
	}
}

func RegisterHealthCheckRoute(router chi.Router) {
	router.Get("/healthcheck", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
}

func RegisterRoutes(router chi.Router, cfg config.Config, usecase usecase.UseCaseInterface) {
	handler := New(cfg, usecase)
	handler.WeatherSetup(router)
}
